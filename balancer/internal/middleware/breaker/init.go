package breaker

import (
	"github.com/goriiin/go-http-balancer/balancer/configs"
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
	cfg          configs.BreakerConfig
	mu           sync.Mutex
	state        State
	failures     int
	lastFailure  time.Time
	halfOpenLeft int
}

func NewBreaker(cfg configs.BreakerConfig) *Breaker {
	return &Breaker{cfg: cfg, state: Closed}
}
