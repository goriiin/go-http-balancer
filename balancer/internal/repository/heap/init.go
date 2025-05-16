package heap

import (
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
	"sync/atomic"
)

type Heap struct {
	data []*domain.Backend
}

func (h *Heap) Len() int {
	return len(h.data)
}

func (h *Heap) Less(i, j int) bool {
	return atomic.LoadInt32(&h.data[i].ConnsCount) < atomic.LoadInt32(&h.data[j].ConnsCount)
}
func (h *Heap) Swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
	h.data[i].Index, h.data[j].Index = i, j
}

func (h *Heap) Push(x *domain.Backend) {
	x.Index = len(h.data)

	h.data = append(h.data, x)
}

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

func (h *Heap) Fix(idx int) {
	if idx < 0 || idx >= len(h.data) {
		return
	}
	if !h.down(idx) {
		h.up(idx)
	}
}

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

func (h *Heap) up(j int) {
	for {
		i := (j - 1) / 2 // родитель
		if j == 0 || !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		j = i
	}
}

func (h *Heap) down(i int) bool {
	n := len(h.data)
	orig := i
	for {
		left := 2*i + 1
		if left >= n {
			break
		}
		smallest := left
		if right := left + 1; right < n && h.Less(right, left) {
			smallest = right
		}
		if !h.Less(smallest, i) {
			break
		}
		h.Swap(i, smallest)
		i = smallest
	}
	return i != orig
}
