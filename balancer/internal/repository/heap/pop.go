package heap

import "github.com/goriiin/go-http-balancer/balancer/internal/domain"

func (h *Heap) Pop() *domain.Backend {
	n := len(h.data)
	if n == 0 {
		return nil
	}

	be := h.data[n-1]
	h.data[n-1] = nil
	h.data = h.data[:n-1]
	be.Index = -1

	return be
}
