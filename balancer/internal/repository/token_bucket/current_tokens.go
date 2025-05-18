package token_bucket

import "time"

func (tb *TokenBucket) CurrentTokens() int64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.lastAccessed = time.Now()

	tb.refill()

	return tb.tokens
}
