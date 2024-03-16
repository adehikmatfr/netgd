package nethttp

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/adehikmatfr/netgd/httpclient"
)

const (
	defaultRetryCount  = 0
	defaultHTTPTimeout = 30 * time.Second
)

type Client struct {
	httpClient *httpclient.Client
}

// NewClient returns a new instance of http Client
func NewClient(opts ...httpclient.Option) *Client {
	httpClient := &httpclient.Client{}

	for _, opt := range opts {
		opt(httpClient)
	}

	if httpClient.GetNetHttpClient() == nil {
		httpClient.SetNetHttpClient(&http.Client{
			Timeout: defaultHTTPTimeout,
		})
	}

	return &Client{
		httpClient: httpClient,
	}
}

// Get makes a HTTP GET request to provided URL
func (c *Client) Get(url string, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return response, err
	}

	request.Header = headers

	return c.Do(request)
}

// Post makes a HTTP POST request to provided URL and requestBody
func (c *Client) Post(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return response, err
	}

	request.Header = headers

	return c.Do(request)
}

// Put makes a HTTP PUT request to provided URL and requestBody
func (c *Client) Put(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return response, err
	}

	request.Header = headers

	return c.Do(request)
}

// Patch makes a HTTP PATCH request to provided URL and requestBody
func (c *Client) Patch(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return response, err
	}

	request.Header = headers

	return c.Do(request)
}

// Delete makes a HTTP DELETE request with provided URL
func (c *Client) Delete(url string, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return response, err
	}

	request.Header = headers

	return c.Do(request)
}

// Do makes an HTTP request with the native `http.Do` interface
func (c *Client) Do(request *http.Request) (*http.Response, error) {
	request.Close = true

	var bodyReader *bytes.Reader

	if request.Body != nil {
		reqData, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(reqData)
		request.Body = ioutil.NopCloser(bodyReader) // prevents closing the body between retries
	}

	errStrList := []string{}
	var response *http.Response

	for i := 0; i <= c.httpClient.GetRetryCount(); i++ {
		if response != nil {
			response.Body.Close()
		}

		c.onRequestStart(request)
		var err error
		response, err = c.httpClient.GetNetHttpClient().Do(request)
		if bodyReader != nil {
			// Reset the body reader after the request since at this point it's already read
			// Note that it's safe to ignore the error here since the 0,0 position is always valid
			_, _ = bodyReader.Seek(0, 0)
		}

		if err != nil {
			errStrList = append(errStrList, err.Error())
			c.onError(request, err)
			backoffTime := c.httpClient.GetRetrier().NextInterval(i)
			time.Sleep(backoffTime)
			continue
		}

		c.onRequestEnd(request, response)

		if response.StatusCode >= http.StatusInternalServerError {
			backoffTime := c.httpClient.GetRetrier().NextInterval(i)
			time.Sleep(backoffTime)
			continue
		}

		errStrList = []string{} // Clear errors if any iteration succeeds
		break
	}

	return response, errors.New(strings.Join(errStrList, ","))
}

func (c *Client) onRequestStart(request *http.Request) {
	for _, interceptor := range c.httpClient.GetInterceptors() {
		interceptor.OnRequestStart(request)
	}
}

func (c *Client) onError(request *http.Request, err error) {
	for _, interceptor := range c.httpClient.GetInterceptors() {
		interceptor.OnError(request, err)
	}
}

func (c *Client) onRequestEnd(request *http.Request, response *http.Response) {
	for _, interceptor := range c.httpClient.GetInterceptors() {
		interceptor.OnRequestEnd(request, response)
	}
}
