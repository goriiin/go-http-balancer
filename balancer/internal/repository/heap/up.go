package heap

func (h *Heap) up(j int) {
	for {
		i := (j - 1) / 2
		if j == 0 || !h.Less(j, i) {
			break
		}

		h.Swap(i, j)
		j = i
	}
}
