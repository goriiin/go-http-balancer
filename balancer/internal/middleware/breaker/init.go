package breaker

import (
	"sync"
	"time"
)

type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

type Breaker struct {
	cfg          Config
	mu           sync.Mutex
	state        State
	failures     int
	lastFailure  time.Time
	halfOpenLeft int
}

func NewBreaker(cfg Config) *Breaker {
	return &Breaker{cfg: cfg, state: Closed}
}
