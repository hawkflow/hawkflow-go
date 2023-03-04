package hawkflow

import (
	"bytes"
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
	returnBody       string
	request          *http.Request
	count            uint8
	clientError      error
}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	c.request = req
	c.count++
	if nil != c.clientError {
		return nil, c.clientError
	}
	return &http.Response{
		StatusCode: c.returnStatusCode,
		Body:       io.NopCloser(bytes.NewReader([]byte(c.returnBody))),
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
	testCases := map[string]struct {
		process               string
		meta                  string
		uid                   string
		count                 uint8
		statusCode            int
		expectedCount         uint8
		expectedRequestBody   string
		expectedRequestMethod string
		expectedRequestUrl    string
		error                 string
	}{
		"Proper request was created": {
			process:               "test_process",
			statusCode:            201,
			count:                 2,
			expectedCount:         1,
			expectedRequestBody:   `{"process":"test_process"}`,
			expectedRequestMethod: "POST",
			expectedRequestUrl:    "https://api.hawkflow.ai/v1/start",
		},
		"Proper full request was created": {
			process:               "test_process",
			meta:                  "test_meta",
			uid:                   "test_uid",
			statusCode:            201,
			count:                 2,
			expectedCount:         1,
			expectedRequestBody:   `{"process":"test_process","meta":"test_meta","uid":"test_uid"}`,
			expectedRequestMethod: "POST",
			expectedRequestUrl:    "https://api.hawkflow.ai/v1/start",
		},
		"Request was validated - process": {
			process:       "invalid process ❌",
			expectedCount: 0,
			error:         "Process parameter contains unsupported characters. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
		},
		"Request was validated - missing process": {
			expectedCount: 0,
			error:         "No process set. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
		},
		"Request was validated - meta": {
			process:       "test_process",
			meta:          "invalid meta ❌",
			expectedCount: 0,
			error:         "Meta parameter contains unsupported characters. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
		},
		"Request was validated - uid": {
			process:       "test_process",
			uid:           "invalid uid ❌",
			expectedCount: 0,
			error:         "UID parameter contains unsupported characters. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
		},
		"Request was retried X times": {
			process:               "test_process",
			statusCode:            500,
			count:                 2,
			expectedCount:         2,
			expectedRequestBody:   `{"process":"test_process"}`,
			expectedRequestMethod: "POST",
			expectedRequestUrl:    "https://api.hawkflow.ai/v1/start",
			error:                 "Connection failed permanently. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			c := &ClientMock{returnStatusCode: testCase.statusCode}
			hfc := New("api_key", OptionHTTPClient(c), OptionMaxRetries(testCase.count))
			err := hfc.Start(testCase.process, testCase.meta, testCase.uid)
			errorMsg := ""
			if err != nil {
				errorMsg = err.Error()
			}

			if nil != c.request {
				reqBody, _ := io.ReadAll(c.request.Body)
				if strings.TrimSpace(string(reqBody)) != testCase.expectedRequestBody {
					t.Errorf("%v expected, got %v", testCase.expectedRequestBody, strings.TrimSpace(string(reqBody)))
				}
				if c.request.URL.String() != testCase.expectedRequestUrl {
					t.Errorf("%v expected, got %v", testCase.expectedRequestUrl, c.request.URL.String())
				}
				if c.request.Method != testCase.expectedRequestMethod {
					t.Errorf("%v expected, got %v", testCase.expectedRequestMethod, c.request.Method)
				}
			}
			if c.count != testCase.expectedCount {
				t.Errorf("%v expected, got %v", testCase.expectedCount, c.count)
			}
			if errorMsg != testCase.error {
				t.Errorf("%v expected, got %v", testCase.error, errorMsg)
			}
		})
	}
}

func TestEnd(t *testing.T) {
	testCases := map[string]struct {
		process               string
		meta                  string
		uid                   string
		count                 uint8
		statusCode            int
		expectedCount         uint8
		expectedRequestBody   string
		expectedRequestMethod string
		expectedRequestUrl    string
		error                 string
	}{
		"Proper request was created": {
			process:               "test_process",
			statusCode:            201,
			count:                 2,
			expectedCount:         1,
			expectedRequestBody:   `{"process":"test_process"}`,
			expectedRequestMethod: "POST",
			expectedRequestUrl:    "https://api.hawkflow.ai/v1/end",
		},
		"Proper full request was created": {
			process:               "test_process",
			meta:                  "test_meta",
			uid:                   "test_uid",
			statusCode:            201,
			count:                 2,
			expectedCount:         1,
			expectedRequestBody:   `{"process":"test_process","meta":"test_meta","uid":"test_uid"}`,
			expectedRequestMethod: "POST",
			expectedRequestUrl:    "https://api.hawkflow.ai/v1/end",
		},
		"Request was validated - process": {
			process:       "invalid process ❌",
			expectedCount: 0,
			error:         "Process parameter contains unsupported characters. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
		},
		"Request was validated - missing process": {
			expectedCount: 0,
			error:         "No process set. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
		},
		"Request was validated - meta": {
			process:       "test_process",
			meta:          "invalid meta ❌",
			expectedCount: 0,
			error:         "Meta parameter contains unsupported characters. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
		},
		"Request was validated - uid": {
			process:       "test_process",
			uid:           "invalid uid ❌",
			expectedCount: 0,
			error:         "UID parameter contains unsupported characters. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
		},
		"Request was retried X times": {
			process:               "test_process",
			statusCode:            500,
			count:                 2,
			expectedCount:         2,
			expectedRequestBody:   `{"process":"test_process"}`,
			expectedRequestMethod: "POST",
			expectedRequestUrl:    "https://api.hawkflow.ai/v1/end",
			error:                 "Connection failed permanently. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			c := &ClientMock{returnStatusCode: testCase.statusCode}
			hfc := New("api_key", OptionHTTPClient(c), OptionMaxRetries(testCase.count))
			err := hfc.End(testCase.process, testCase.meta, testCase.uid)
			errorMsg := ""
			if err != nil {
				errorMsg = err.Error()
			}

			if nil != c.request {
				reqBody, _ := io.ReadAll(c.request.Body)
				if strings.TrimSpace(string(reqBody)) != testCase.expectedRequestBody {
					t.Errorf("%v expected, got %v", testCase.expectedRequestBody, strings.TrimSpace(string(reqBody)))
				}
				if c.request.URL.String() != testCase.expectedRequestUrl {
					t.Errorf("%v expected, got %v", testCase.expectedRequestUrl, c.request.URL.String())
				}
				if c.request.Method != testCase.expectedRequestMethod {
					t.Errorf("%v expected, got %v", testCase.expectedRequestMethod, c.request.Method)
				}
			}
			if c.count != testCase.expectedCount {
				t.Errorf("%v expected, got %v", testCase.expectedCount, c.count)
			}
			if errorMsg != testCase.error {
				t.Errorf("%v expected, got %v", testCase.error, errorMsg)
			}
		})
	}
}

func TestException(t *testing.T) {
	testCases := map[string]struct {
		process               string
		meta                  string
		message               string
		count                 uint8
		statusCode            int
		expectedCount         uint8
		expectedRequestBody   string
		expectedRequestMethod string
		expectedRequestUrl    string
		error                 string
	}{
		"Proper request was created": {
			process:               "test_process",
			message:               "test exception message",
			statusCode:            201,
			count:                 2,
			expectedCount:         1,
			expectedRequestBody:   `{"process":"test_process","exception":"test exception message"}`,
			expectedRequestMethod: "POST",
			expectedRequestUrl:    "https://api.hawkflow.ai/v1/exception",
		},
		"Proper full request was created": {
			process:               "test_process",
			meta:                  "test_meta",
			message:               "test exception message",
			statusCode:            201,
			count:                 2,
			expectedCount:         1,
			expectedRequestBody:   `{"process":"test_process","meta":"test_meta","exception":"test exception message"}`,
			expectedRequestMethod: "POST",
			expectedRequestUrl:    "https://api.hawkflow.ai/v1/exception",
		},
		"Request was validated - missing process": {
			expectedCount: 0,
			error:         "No process set. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
		},
		"Request was validated - missing message": {
			process:               "test_process",
			statusCode:            201,
			count:                 2,
			expectedCount:         1,
			expectedRequestBody:   `{"process":"test_process"}`,
			expectedRequestMethod: "POST",
			expectedRequestUrl:    "https://api.hawkflow.ai/v1/exception",
		},
		"Request was validated - process": {
			process:       "invalid process ❌",
			expectedCount: 0,
			error:         "Process parameter contains unsupported characters. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
		},
		"Request was retried X times": {
			process:               "test_process",
			message:               "test exception message",
			statusCode:            500,
			count:                 2,
			expectedCount:         2,
			expectedRequestBody:   `{"process":"test_process","exception":"test exception message"}`,
			expectedRequestMethod: "POST",
			expectedRequestUrl:    "https://api.hawkflow.ai/v1/exception",
			error:                 "Connection failed permanently. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			c := &ClientMock{returnStatusCode: testCase.statusCode}
			hfc := New("api_key", OptionHTTPClient(c), OptionMaxRetries(testCase.count))
			err := hfc.Exception(testCase.process, testCase.meta, testCase.message)
			errorMsg := ""
			if err != nil {
				errorMsg = err.Error()
			}

			if nil != c.request {
				reqBody, _ := io.ReadAll(c.request.Body)
				if strings.TrimSpace(string(reqBody)) != testCase.expectedRequestBody {
					t.Errorf("%v expected, got %v", testCase.expectedRequestBody, strings.TrimSpace(string(reqBody)))
				}
				if c.request.URL.String() != testCase.expectedRequestUrl {
					t.Errorf("%v expected, got %v", testCase.expectedRequestUrl, c.request.URL.String())
				}
				if c.request.Method != testCase.expectedRequestMethod {
					t.Errorf("%v expected, got %v", testCase.expectedRequestMethod, c.request.Method)
				}
			}
			if c.count != testCase.expectedCount {
				t.Errorf("%v expected, got %v", testCase.expectedCount, c.count)
			}
			if errorMsg != testCase.error {
				t.Errorf("%v expected, got %v", testCase.error, errorMsg)
			}
		})
	}
}

func TestMetrics(t *testing.T) {
	testCases := map[string]struct {
		process               string
		meta                  string
		items                 map[string]float64
		count                 uint8
		statusCode            int
		expectedCount         uint8
		expectedRequestBody   string
		expectedRequestMethod string
		expectedRequestUrl    string
		error                 string
	}{
		"Proper request was created": {
			process:               "test_process",
			items:                 map[string]float64{"key": 123},
			statusCode:            201,
			count:                 2,
			expectedCount:         1,
			expectedRequestBody:   `{"process":"test_process","items":{"key":123}}`,
			expectedRequestMethod: "POST",
			expectedRequestUrl:    "https://api.hawkflow.ai/v1/metrics",
		},
		"Proper full request was created": {
			process:               "test_process",
			meta:                  "test_meta",
			items:                 map[string]float64{"key": 123},
			statusCode:            201,
			count:                 2,
			expectedCount:         1,
			expectedRequestBody:   `{"process":"test_process","meta":"test_meta","items":{"key":123}}`,
			expectedRequestMethod: "POST",
			expectedRequestUrl:    "https://api.hawkflow.ai/v1/metrics",
		},
		"Request was validated - process": {
			process:       "invalid process ❌",
			expectedCount: 0,
			error:         "Process parameter contains unsupported characters. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
		},
		"Request was validated - missing process": {
			expectedCount: 0,
			error:         "No process set. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
		},
		"Request was validated - missing items": {
			process:       "test_process",
			expectedCount: 0,
			error:         "No items set. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
		},
		"Request was retried X times": {
			process:               "test_process",
			items:                 map[string]float64{"key": 123},
			statusCode:            500,
			count:                 2,
			expectedCount:         2,
			expectedRequestBody:   `{"process":"test_process","items":{"key":123}}`,
			expectedRequestMethod: "POST",
			expectedRequestUrl:    "https://api.hawkflow.ai/v1/metrics",
			error:                 "Connection failed permanently. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			c := &ClientMock{returnStatusCode: testCase.statusCode}
			hfc := New("api_key", OptionHTTPClient(c), OptionMaxRetries(testCase.count))
			err := hfc.Metrics(testCase.process, testCase.meta, testCase.items)
			errorMsg := ""
			if err != nil {
				errorMsg = err.Error()
			}

			if nil != c.request {
				reqBody, _ := io.ReadAll(c.request.Body)
				if strings.TrimSpace(string(reqBody)) != testCase.expectedRequestBody {
					t.Errorf("%v expected, got %v", testCase.expectedRequestBody, strings.TrimSpace(string(reqBody)))
				}
				if c.request.URL.String() != testCase.expectedRequestUrl {
					t.Errorf("%v expected, got %v", testCase.expectedRequestUrl, c.request.URL.String())
				}
				if c.request.Method != testCase.expectedRequestMethod {
					t.Errorf("%v expected, got %v", testCase.expectedRequestMethod, c.request.Method)
				}
			}
			if c.count != testCase.expectedCount {
				t.Errorf("%v expected, got %v", testCase.expectedCount, c.count)
			}
			if errorMsg != testCase.error {
				t.Errorf("%v expected, got %v", testCase.error, errorMsg)
			}
		})
	}
}

func TestSendWithRetry(t *testing.T) {
	testCases := map[string]struct {
		apiKey        string
		statusCode    int
		clientError   error
		count         uint8
		returnBody    string
		expectedCount uint8
		error         string
	}{
		"Valid finished without retry": {
			apiKey:        "api_key",
			statusCode:    201,
			count:         4,
			expectedCount: 1,
		},
		"Valid 4 retries": {
			apiKey:        "api_key",
			statusCode:    500,
			count:         4,
			expectedCount: 4,
			error:         "Connection failed permanently. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
		},
		"Invalid interrupted": {
			apiKey:        "invalid api key ❌",
			statusCode:    0,
			count:         4,
			expectedCount: 0,
			error:         "Invalid API Key format. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
		},
		"Invalid 0 retries set": {
			apiKey:        "api_key",
			statusCode:    201,
			count:         0,
			expectedCount: 0,
			error:         "Connection failed permanently. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
		},
		"Unauthorised": {
			apiKey:        "api_key",
			statusCode:    401,
			count:         4,
			returnBody:    `{"status":"401","message":"unauthorized"}`,
			expectedCount: 1,
			error:         `{"status":"401","message":"unauthorized"}`,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			c := &ClientMock{returnStatusCode: testCase.statusCode, returnBody: testCase.returnBody, clientError: testCase.clientError}
			hfc := New(testCase.apiKey, OptionHTTPClient(c))
			req := &request{}
			err := hfc.sendWithRetry(req, "/v1/test", testCase.count)
			errorMsg := ""
			if err != nil {
				errorMsg = err.Error()
			}

			if c.count != testCase.expectedCount {
				t.Errorf("%v expected, got %v", testCase.expectedCount, c.count)
			}
			if errorMsg != testCase.error {
				t.Errorf("%v expected, got %v", testCase.error, errorMsg)
			}
		})
	}
}

func TestSendValidateApiKeyIsCalled(t *testing.T) {
	testCases := map[string]struct {
		apiKey        string
		statusCode    int
		error         string
		expectedRetry bool
	}{
		"Valid API key": {
			apiKey:        "api_key",
			statusCode:    201,
			error:         "",
			expectedRetry: false,
		},
		"Empty API key": {
			apiKey:        "",
			statusCode:    -1,
			error:         "No API Key set. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
			expectedRetry: false,
		},
		"Invalid API key": {
			apiKey:        "invalid api key ❌",
			statusCode:    -1,
			error:         "Invalid API Key format. Please see documentation at https://docs.hawkflow.ai/integration/index.html",
			expectedRetry: false,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			c := &ClientMock{returnStatusCode: 201}
			hfc := New(testCase.apiKey, OptionHTTPClient(c))
			req := &request{}
			retry, err := hfc.send(req, "/v1/test")
			errorMsg := ""
			if err != nil {
				errorMsg = err.Error()
			}

			if retry != testCase.expectedRetry {
				t.Errorf("%v expected, got %v", testCase.expectedRetry, retry)
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
		statusCode     int
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
			statusCode:     201,
			expectedUrl:    "https://api.hawkflow.ai/v1/test",
			expectedMethod: "POST",
			expectedBody:   `{"process":"test_process","meta":"test_meta"}`,
		},
		"Invalid response": {
			path:           "test",
			statusCode:     500,
			expectedUrl:    "https://api.hawkflow.ai/v1/test",
			expectedMethod: "POST",
			expectedBody:   `{"status":"500","message":"Server error"}`,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			c := &ClientMock{returnStatusCode: testCase.statusCode, returnBody: testCase.expectedBody}
			hfc := New("api_key", OptionHTTPClient(c))
			_, err := hfc.send(testCase.req, testCase.path)
			reqBody, _ := io.ReadAll(c.request.Body)

			if c.request.URL.String() != testCase.expectedUrl {
				t.Errorf("%v expected, got %v", testCase.expectedUrl, c.request.URL.String())
			}
			if c.request.Method != testCase.expectedMethod {
				t.Errorf("%v expected, got %v", testCase.expectedMethod, c.request.Method)
			}
			if nil != err {
				if err.Error() != testCase.expectedBody {
					t.Errorf("%v expected, got %v", testCase.expectedBody, err)
				}
			} else {
				if strings.TrimSpace(string(reqBody)) != testCase.expectedBody {
					t.Errorf("%v expected, got %v", testCase.expectedBody, strings.TrimSpace(string(reqBody)))
				}
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
