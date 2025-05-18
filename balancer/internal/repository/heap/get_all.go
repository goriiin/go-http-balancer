package heap

import "github.com/goriiin/go-http-balancer/balancer/internal/domain"

func (h *Heap) GetAll() []*domain.Backend {
	dataCopy := make([]*domain.Backend, len(h.data))
	copy(dataCopy, h.data)

	return dataCopy
}
