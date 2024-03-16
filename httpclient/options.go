package httpclient

import (
	"time"
)

// Option represents the client options
type Option func(*Client)

// WithHTTPTimeout sets hystrix timeout
func WithHTTPTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithRetryCount sets the retry count for the hystrixHTTPClient
func WithRetryCount(retryCount int) Option {
	return func(c *Client) {
		c.retryCount = retryCount
	}
}

// WithRetrier sets the strategy for retrying
func WithRetrier(retrier Retriable) Option {
	return func(c *Client) {
		c.retrier = retrier
	}
}

// WithNetHTTPClient sets a custom http client
func WithNetHTTPClient(client NetHttpDoer) Option {
	return func(c *Client) {
		c.netHttpClient = client
	}
}

// WithNetHTTPClient sets a custom http client
func WithFastHTTPClient(client FastHttpDoer) Option {
	return func(c *Client) {
		c.fastHttpClient = client
	}
}

// WithInterceptor sets interceptor client
func WithInterceptor(i []Interceptor) Option {
	return func(c *Client) {
		c.interceptors = i
	}
}
