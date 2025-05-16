package back

import (
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type breaker interface {
	Allow() bool
}

type Backend struct {
	URL     *url.URL
	Proxy   *httputil.ReverseProxy
	alive   atomic.Bool
	breaker breaker
}

func NewBackend() {

}
