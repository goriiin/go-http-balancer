package heap

import "sync/atomic"

func (h *Heap) Less(i, j int) bool {
	return atomic.LoadInt32(&h.data[i].ConnsCount) < atomic.LoadInt32(&h.data[j].ConnsCount)
}
