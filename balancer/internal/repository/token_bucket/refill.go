package token_bucket

import "time"

func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)

	tokensToAdd := (elapsed.Nanoseconds() * tb.refillRate) / time.Second.Nanoseconds()

	if tokensToAdd > 0 {
		tb.tokens += tokensToAdd
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}

		tb.lastRefill = now
	}
}
