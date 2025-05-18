package heap

import "github.com/goriiin/go-http-balancer/balancer/internal/domain"

func (h *Heap) Push(x *domain.Backend) {
	x.Index = len(h.data)

	h.data = append(h.data, x)
}
