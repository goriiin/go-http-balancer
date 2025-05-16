package breaker

import "time"

func (b *Breaker) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case Closed:
		return true
	case Open:
		if time.Since(b.lastFailure) > b.cfg.ResetTimeout {
			b.state = HalfOpen
			b.halfOpenLeft = b.cfg.HalfOpenRequests

			return true
		}

		return false
	case HalfOpen:
		return b.halfOpenLeft > 0
	default:
		return false
	}
}
