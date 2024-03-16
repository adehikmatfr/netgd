package httpclient

import (
	"net/http"
)

type Interceptor interface {
	OnRequestStart(*http.Request)
	OnRequestEnd(*http.Request, *http.Response)
	OnError(*http.Request, error)
}
