package back

import (
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type breaker interface {
	Allow() bool
}

type Backend struct {
	URL      *url.URL
	Proxy    *httputil.ReverseProxy
	domainBe *domain.Backend
	log      *slog.Logger
}

func NewBackend(domainBe *domain.Backend, logger *slog.Logger) *Backend {
	proxy := httputil.NewSingleHostReverseProxy(domainBe.URL)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = domainBe.URL.Host
	}

	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		logger.Error("backend proxy error",
			slog.String("backend_url", domainBe.URL.String()),
			slog.String("request_path", req.URL.Path),
			slog.String("error", err.Error()),
		)

		rw.WriteHeader(http.StatusBadGateway)
	}

	return &Backend{
		URL:      domainBe.URL,
		Proxy:    proxy,
		domainBe: domainBe,
		log:      logger,
	}
}
