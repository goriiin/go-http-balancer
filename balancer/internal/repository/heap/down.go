package heap

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
