package httpclient

import (
	"io"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

type NetHttpDoer interface {
	Do(request *http.Request) (*http.Response, error)
}

type FastHttpDoer interface {
	Do(req *fasthttp.Request, resp *fasthttp.Response) error
}

type HttpClient interface {
	Get(url string, headers http.Header) (*http.Response, error)
	Post(url string, body io.Reader, headers http.Header) (*http.Response, error)
	Put(url string, body io.Reader, headers http.Header) (*http.Response, error)
	Patch(url string, body io.Reader, headers http.Header) (*http.Response, error)
	Delete(url string, headers http.Header) (*http.Response, error)
}

type Client struct {
	netHttpClient  NetHttpDoer
	fastHttpClient FastHttpDoer
	timeout        time.Duration
	retryCount     int
	retrier        Retriable
	interceptors   []Interceptor
}

// Getter & setter
func (c *Client) GetNetHttpClient() NetHttpDoer {
	return c.netHttpClient
}

func (c *Client) SetNetHttpClient(d NetHttpDoer) {
	c.netHttpClient = d
}

func (c *Client) GetFastHttpClient() FastHttpDoer {
	return c.fastHttpClient
}

func (c *Client) SetFastHttpClient(d FastHttpDoer) {
	c.fastHttpClient = d
}

func (c *Client) GetInterceptors() []Interceptor {
	return c.interceptors
}

func (c *Client) GetTimeout() time.Duration {
	return c.timeout
}

func (c *Client) GetRetryCount() int {
	return c.retryCount
}

func (c *Client) GetRetrier() Retriable {
	return c.retrier
}

func NewClient(hc HttpClient) *HttpClient {
	return &hc
}
