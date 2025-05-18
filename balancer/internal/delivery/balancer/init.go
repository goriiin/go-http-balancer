package balancer

import (
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
	"log/slog"
	"net/http"
	"sync"
)

type pool interface {
	Next() *domain.Backend
	Add(be *domain.Backend)
	MarkHealthy(be *domain.Backend)
	MarkUnhealthy(be *domain.Backend)
	Done(be *domain.Backend)
	GetAll() []*domain.Backend
}

type Balancer struct {
	backends map[string]http.Handler
	mu       sync.RWMutex
	pool     pool
	log      *slog.Logger
}

func NewBalancer(p pool, log *slog.Logger) *Balancer {
	return &Balancer{
		backends: make(map[string]http.Handler),
		pool:     p,
		log:      log,
	}
}

func (b *Balancer) AddBackendHandler(urlStr string, handler http.Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.backends[urlStr] = handler

	b.log.Info("added backend handler to balancer's dispatch map", slog.String("url", urlStr))
}
