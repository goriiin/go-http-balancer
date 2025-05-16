package round_robin

import (
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
	"sync"
)

type Pool struct {
	mu        sync.RWMutex
	healthy   []*domain.Backend
	unhealthy []*domain.Backend
	pos       int
}
