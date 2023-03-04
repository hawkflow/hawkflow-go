package hawkflow

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	// _ENDPOINT must end with backslash
	_ENDPOINT = "https://api.hawkflow.ai/v1/"
	_TIMEOUT  = 1000 * time.Millisecond
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type client struct {
	apiKey     string
	endpoint   string
	maxRetries uint8
	debug      bool
	logger     logger
	httpClient httpClient
}

type logger interface {
	Output(int, string) error
	Print(...interface{})
}

type request struct {
	Process          string             `json:"process"`
	Meta             string             `json:"meta,omitempty"`
	UID              string             `json:"uid,omitempty"`
	ExceptionMessage string             `json:"exception,omitempty"`
	Items            map[string]float64 `json:"items,omitempty"`
}

type option func(*client)

func OptionMaxRetries(maxRetries uint8) func(*client) {
	return func(hfc *client) { hfc.maxRetries = maxRetries }
}

// OptionTimeout overwrites httpClient
func OptionTimeout(timeout time.Duration) func(*client) {
	return func(hfc *client) { hfc.httpClient = &http.Client{Timeout: timeout} }
}

func OptionDebug(b bool) func(*client) {
	return func(hfc *client) { hfc.debug = b }
}

func OptionLogger(l logger) func(*client) {
	return func(hfc *client) {
		hfc.logger = l
	}
}

func OptionHTTPClient(c httpClient) func(*client) {
	return func(hfc *client) { hfc.httpClient = c }
}

func New(apiKey string, options ...option) *client {
	hfc := &client{
		apiKey:     apiKey,
		endpoint:   _ENDPOINT,
		maxRetries: 3,
		debug:      false,
		logger:     log.New(os.Stderr, "hawkflow", log.LstdFlags|log.Lshortfile),
		httpClient: &http.Client{
			Timeout: _TIMEOUT,
		},

		//// https://github.com/hashicorp/go-cleanhttp/blob/02f12f05b908335f65f124c1a7f1ec45f5c42a35/cleanhttp.go#L26
		//httpClient: &http.Client{
		//	Transport: &http.Transport{
		//		MaxIdleConns:          100,
		//		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		//		IdleConnTimeout:       3 * time.Second,
		//		TLSHandshakeTimeout:   1 * time.Second,
		//		ExpectContinueTimeout: 1 * time.Second,
		//		ForceAttemptHTTP2:     true,
		//	},
		//	Timeout: 1 * time.Second,
		//},
	}

	for _, opt := range options {
		opt(hfc)
	}

	return hfc
}

func (hfc *client) log(m string) {
	if hfc.debug {
		hfc.logger.Print(fmt.Sprintf("HF %s\n", m))
	}
}

func (hfc *client) Start(process, meta, uid string) error {
	r := &request{
		Process: process,
		Meta:    meta,
		UID:     uid,
	}

	err := validateTimedData(r)
	if err != nil {
		return err
	}

	hfc.log(fmt.Sprintf("Start: %s", process))

	return hfc.sendWithRetry(r, "start", hfc.maxRetries)
}

func (hfc *client) End(process, meta, uid string) error {
	r := &request{
		Process: process,
		Meta:    meta,
		UID:     uid,
	}

	err := validateTimedData(r)
	if err != nil {
		return err
	}

	hfc.log(fmt.Sprintf("End: %s", process))

	return hfc.sendWithRetry(r, "end", hfc.maxRetries)
}

func (hfc *client) Exception(process, meta, message string) error {
	r := &request{
		Process:          process,
		Meta:             meta,
		ExceptionMessage: message,
	}

	err := validateException(r)
	if err != nil {
		return err
	}

	hfc.log(fmt.Sprintf("Exception: %s", process))

	return hfc.sendWithRetry(r, "exception", hfc.maxRetries)
}

func (hfc *client) Metrics(process, meta string, items map[string]float64) error {
	r := &request{
		Process: process,
		Meta:    meta,
		Items:   items,
	}

	err := validateMetrics(r)
	if err != nil {
		return err
	}

	hfc.log(fmt.Sprintf("Metrics: %s", process))

	return hfc.sendWithRetry(r, "metrics", hfc.maxRetries)
}

func (hfc *client) sendWithRetry(r *request, path string, count uint8) error {
	if 0 >= count {
		return createError("Connection failed permanently.")
	}

	retry, err := hfc.send(r, path)
	if nil != err {
		hfc.log(fmt.Sprintf("Connection failed on attempt %d with error: %s", count, err))
		if !retry {
			return err
		}
		return hfc.sendWithRetry(r, path, count-1)
	}

	return nil
}

func (hfc *client) send(r *request, path string) (bool, error) {
	err := validateApiKey(hfc.apiKey)
	if err != nil {
		return false, err
	}

	body := new(bytes.Buffer)
	err = json.NewEncoder(body).Encode(r)
	if err != nil {
		return false, err
	}

	hfc.log(fmt.Sprintf("Requesting path: %s", path))
	hfc.log(fmt.Sprintf("Sending data: %s", body))

	req, err := http.NewRequest("POST", hfc.endpoint+path, body)
	if err != nil {
		return false, err
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-hawkflow-api-key", hfc.apiKey)

	resp, err := hfc.httpClient.Do(req)
	if err != nil {
		return true, err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	hfc.log(fmt.Sprintf("Response Status: %s", resp.Status))
	hfc.log(fmt.Sprintf("Response Body: %s", respBody))

	if http.StatusUnauthorized == resp.StatusCode {
		return false, errors.New(string(respBody))
	}
	if http.StatusCreated != resp.StatusCode {
		return true, errors.New(string(respBody))
	}

	return false, nil
}
