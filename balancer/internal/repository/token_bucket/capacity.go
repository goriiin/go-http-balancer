package token_bucket

func (tb *TokenBucket) Capacity() int64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	return tb.capacity
}
