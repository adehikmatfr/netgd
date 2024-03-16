package fasthttp

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/adehikmatfr/netgd/httpclient"
	"github.com/valyala/fasthttp"
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

	if httpClient.GetFastHttpClient() == nil {
		httpClient.SetFastHttpClient(&fasthttp.Client{
			ReadTimeout:  defaultHTTPTimeout,
			WriteTimeout: defaultHTTPTimeout,
		})
	}

	return &Client{
		httpClient: httpClient,
	}
}

// Post makes a HTTP POST request to provided URL and requestBody
func (c *Client) Post(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

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
func (c *Client) Do(request *fasthttp.Request, response *fasthttp.Response) error {
	var bodyReader *bytes.Reader

	errStrList := []string{}

	for i := 0; i <= c.httpClient.GetRetryCount(); i++ {
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
