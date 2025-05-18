package token_bucket

import "time"

func (tb *TokenBucket) LastAccessed() time.Time {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	return tb.lastAccessed
}
