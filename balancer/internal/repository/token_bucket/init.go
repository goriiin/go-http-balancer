package token_bucket

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity     int64
	tokens       int64
	refillRate   int64
	lastRefill   time.Time
	mu           sync.Mutex
	lastAccessed time.Time
}

func NewTokenBucket(capacity, refillRatePerSecond int64) *TokenBucket {
	if capacity <= 0 {
		capacity = 1
	}

	if refillRatePerSecond <= 0 {
		refillRatePerSecond = 1
	}

	return &TokenBucket{
		capacity:     capacity,
		tokens:       capacity,
		refillRate:   refillRatePerSecond,
		lastRefill:   time.Now(),
		lastAccessed: time.Now(),
	}
}
