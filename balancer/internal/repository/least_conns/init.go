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
}

type Pool struct {
	mu      sync.Mutex
	healthy heap
}
