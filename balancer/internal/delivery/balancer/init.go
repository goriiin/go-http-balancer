package balancer

import (
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
)

type pool interface {
	Next() *domain.Backend
	MarkUnhealthy(u *domain.Backend)
}

type Balancer struct {
	m    map[*url.URL]http.Handler
	mux  *sync.RWMutex
	pool pool
	log  *slog.Logger
}

func NewBalancer(p pool, log *slog.Logger) *Balancer {
	return &Balancer{
		m:    make(map[*url.URL]http.Handler, 100),
		mux:  new(sync.RWMutex),
		pool: p,
		log:  log,
	}
}
