package round_robin

import (
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
	"sync"
)

type Pool struct {
	mu        sync.RWMutex
	healthy   []*domain.Backend
	unhealthy map[string]*domain.Backend
	pos       uint32
}

func NewPool() *Pool {
	return &Pool{
		healthy:   make([]*domain.Backend, 0),
		unhealthy: make(map[string]*domain.Backend),
	}
}
