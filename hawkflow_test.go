package hawkflow

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

type ClientMock struct {
	returnStatusCode int
	request          *http.Request
	retryCount       uint8
	clientError      error
}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	c.request = req
	c.retryCount++
	if nil != c.clientError {
		return nil, c.clientError
	}
	return &http.Response{
		StatusCode: c.returnStatusCode,
		Body:       io.NopCloser(bytes.NewReader([]byte("test"))),
	}, nil
}

func TestOptionMaxRetries(t *testing.T) {
	maxRetries := uint8(7)
	hfc := New("api_key", OptionMaxRetries(maxRetries))

	if hfc.maxRetries != maxRetries {
		t.Errorf("Setting max retries failed.")
	}
}

func TestOptionTimeout(t *testing.T) {
	timeout := 123 * time.Millisecond
	hfc := New("api_key", OptionTimeout(timeout))

	if fmt.Sprint(reflect.Indirect(reflect.ValueOf(hfc.httpClient)).FieldByName("Timeout")) != "123ms" {
		t.Errorf("Setting timeout failed.")
	}
}

func TestOptionDebug(t *testing.T) {
	hfc := New("api_key", OptionDebug(true))

	if hfc.debug != true {
		t.Errorf("Setting debug failed.")
	}
}

func TestStart(t *testing.T) {
	// test that proper request was send
}

func TestEnd(t *testing.T) {
	// test that proper request was send
}

func TestException(t *testing.T) {
	// test that proper request was send
}

func TestMetrics(t *testing.T) {
	// test that proper request was send
}

func TestSendWithRetry(t *testing.T) {
	// valid, retreid X times
	// invalid, when `send` returns 0, retry interrupted with error from send
	// when retry is 0, return error

	testCases := map[string]struct {
		statusCode    int
		clientError   error
		retry         uint8
		expectedCount uint8
		error         string
	}{
		"Valid finished without retry": {
			statusCode:    201,
			retry:         4,
			expectedCount: 1,
			error:         "",
		},
		"Valid 4 retries": {
			statusCode:    500,
			retry:         4,
			expectedCount: 4,
			error:         "Connection failed permanently. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
		},
		"Invalid interrupted": {
			statusCode:    0,
			clientError:   errors.New("server is down"),
			retry:         4,
			expectedCount: 1,
			error:         "server is down",
		},
		"Invalid 0 retries set": {
			statusCode:    201,
			retry:         0,
			expectedCount: 0,
			error:         "Connection failed permanently. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			c := &ClientMock{returnStatusCode: testCase.statusCode, clientError: testCase.clientError}
			hfc := New("api_key", OptionHTTPClient(c))
			req := &request{}
			err := hfc.sendWithRetry(req, "/v1/test", testCase.retry)
			errorMsg := ""
			if err != nil {
				errorMsg = err.Error()
			}

			if c.retryCount != testCase.expectedCount {
				t.Errorf("%v expected, got %v", testCase.expectedCount, c.retryCount)
			}
			if errorMsg != testCase.error {
				t.Errorf("%v expected, got %v", testCase.error, errorMsg)
			}
		})
	}
}

func TestSendValidateApiKeyIsCalled(t *testing.T) {
	testCases := map[string]struct {
		apiKey     string
		statusCode int
		error      string
	}{
		"Valid API key": {
			apiKey:     "api_key",
			statusCode: 201,
			error:      "",
		},
		"Empty API key": {
			apiKey:     "",
			statusCode: 0,
			error:      "No API Key set. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
		},
		"Invalid API key": {
			apiKey:     "api % # key",
			statusCode: 0,
			error:      "Invalid API Key format. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			c := &ClientMock{returnStatusCode: 201}
			hfc := New(testCase.apiKey, OptionHTTPClient(c))
			req := &request{}
			statusCode, err := hfc.send(req, "/v1/test")
			errorMsg := ""
			if err != nil {
				errorMsg = err.Error()
			}

			if statusCode != testCase.statusCode {
				t.Errorf("%v expected, got %v", testCase.statusCode, statusCode)
			}
			if errorMsg != testCase.error {
				t.Errorf("%v expected, got %v", testCase.error, errorMsg)
			}
		})
	}
}

func TestSendRequest(t *testing.T) {
	testCases := map[string]struct {
		path           string
		req            *request
		expectedUrl    string
		expectedMethod string
		expectedBody   string
	}{
		"Valid request": {
			path: "test",
			req: &request{
				Process: "test_process",
				Meta:    "test_meta",
			},
			expectedUrl:    "https://api.hawkflow.ai/v1/test",
			expectedMethod: "POST",
			expectedBody:   `{"process":"test_process","meta":"test_meta"}`,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			c := &ClientMock{returnStatusCode: 201}
			hfc := New("api_key", OptionHTTPClient(c))
			_, _ = hfc.send(testCase.req, testCase.path)
			reqBody, _ := io.ReadAll(c.request.Body)

			if c.request.URL.String() != testCase.expectedUrl {
				t.Errorf("%v expected, got %v", testCase.expectedUrl, c.request.URL.String())
			}
			if c.request.Method != testCase.expectedMethod {
				t.Errorf("%v expected, got %v", testCase.expectedMethod, c.request.Method)
			}
			if strings.TrimSpace(string(reqBody)) != testCase.expectedBody {
				t.Errorf("%v expected, got %v", testCase.expectedBody, strings.TrimSpace(string(reqBody)))
			}
		})
	}
}

func TestLog(t *testing.T) {
	testCases := map[string]struct {
		message string
		debug   bool
		log     string
	}{
		"Valid log": {
			message: "test log message",
			debug:   true,
			log:     "HF test log message\n",
		},
		"No log when debug is disabled": {
			message: "test log message",
			debug:   false,
			log:     "",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			buf := bytes.NewBufferString("")
			logger := log.New(buf, "", 0)
			hfc := New("api_key", OptionLogger(logger), OptionDebug(testCase.debug))
			hfc.log(testCase.message)

			if buf.String() != testCase.log {
				t.Errorf("%v expected, got %v", testCase.log, buf.String())
			}
		})
	}
}
