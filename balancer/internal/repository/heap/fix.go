package heap

func (h *Heap) Fix(idx int) {
	if idx < 0 || idx >= len(h.data) {
		return
	}
	if !h.down(idx) {
		h.up(idx)
	}
}
