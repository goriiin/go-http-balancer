package store

import (
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
	"github.com/goriiin/go-http-balancer/balancer/internal/repository/token_bucket"
	"log/slog"
	"sync"
	"time"
)

type MemoryStore struct {
	buckets           map[string]*token_bucket.TokenBucket
	clientConfigs     map[string]domain.ClientInfo
	mu                sync.RWMutex
	defaultConfig     domain.ClientInfo
	logger            *slog.Logger
	cleanupInterval   time.Duration
	inactiveThreshold time.Duration
	cleanupTicker     *time.Ticker
	stopCleanup       chan struct{}
}

func NewMemoryStore(
	defaultCapacity, defaultRatePerSecond int64,
	cleanupInterval, inactiveThreshold time.Duration,
	logger *slog.Logger,
) *MemoryStore {
	store := &MemoryStore{
		buckets:       make(map[string]*token_bucket.TokenBucket),
		clientConfigs: make(map[string]domain.ClientInfo),
		defaultConfig: domain.ClientInfo{
			Capacity: defaultCapacity,
			Rate:     defaultRatePerSecond,
		},
		logger:            logger.With(slog.String("component", "ratelimit_memorystore")),
		cleanupInterval:   cleanupInterval,
		inactiveThreshold: inactiveThreshold,
		stopCleanup:       make(chan struct{}),
	}

	if store.cleanupInterval > 0 && store.inactiveThreshold > 0 {
		store.cleanupTicker = time.NewTicker(store.cleanupInterval)
		go store.startCleanupRoutine()

		store.logger.Info("Rate limit bucket cleanup routine started",
			slog.Duration("interval", store.cleanupInterval),
			slog.Duration("threshold", store.inactiveThreshold))

	} else {
		store.logger.Info("Rate limit bucket cleanup routine will not start (interval or threshold is zero)")
	}

	return store
}
