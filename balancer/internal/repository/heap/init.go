package heap

import (
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
)

type Heap struct {
	data []*domain.Backend
}
