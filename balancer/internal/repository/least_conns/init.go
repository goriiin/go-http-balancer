package least_conns

import (
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
	"sync"
)

type heap interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
	Push(x *domain.Backend)
	Pop() *domain.Backend
	Fix(idx int)
	Remove(idx int) *domain.Backend
	GetAll() []*domain.Backend
}

type Pool struct {
	mu         sync.Mutex
	healthy    heap
	unhealthy  map[string]*domain.Backend
	backendMap map[string]*domain.Backend
}

func NewPool(heap heap) *Pool {
	return &Pool{
		healthy:    heap,
		unhealthy:  make(map[string]*domain.Backend),
		backendMap: make(map[string]*domain.Backend),
	}
}
