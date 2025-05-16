package breaker

import "time"

func (b *Breaker) Failure() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.failures++
	b.lastFailure = time.Now()

	if b.state == HalfOpen || b.failures >= b.cfg.FailureThreshold {
		b.state = Open
	}
}
