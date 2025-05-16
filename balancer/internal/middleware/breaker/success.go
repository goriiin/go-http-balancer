package breaker

func (b *Breaker) Success() {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case Closed:
		b.failures = 0
	case HalfOpen:
		b.halfOpenLeft--
		if b.halfOpenLeft == 0 {
			b.state = Closed
			b.failures = 0
		}
	}
}
