package random

import (
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
	"math/rand"
	"sync"
	"time"
)

type Pool struct {
	mu        sync.RWMutex
	healthy   []*domain.Backend
	unhealthy map[string]*domain.Backend
	rng       *rand.Rand
}

func NewPool() *Pool {
	return &Pool{
		healthy:   make([]*domain.Backend, 0),
		unhealthy: make(map[string]*domain.Backend),
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())), // Seeded RNG
	}
}
