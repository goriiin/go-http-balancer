package heap

import "github.com/goriiin/go-http-balancer/balancer/internal/domain"

func (h *Heap) Remove(idx int) *domain.Backend {
	n := len(h.data)
	if idx < 0 || idx >= n {
		return nil
	}

	h.Swap(idx, n-1)
	be := h.Pop()
	if idx < n-1 {
		h.Fix(idx)
	}

	return be
}
