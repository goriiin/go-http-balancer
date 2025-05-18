package store

import (
	"github.com/goriiin/go-http-balancer/balancer/internal/repository/token_bucket"
	"log/slog"
)

func (ms *MemoryStore) GetBucket(clientID string) (*token_bucket.TokenBucket, error) {
	ms.mu.RLock()
	bucket, exists := ms.buckets[clientID]
	ms.mu.RUnlock()

	if exists {
		return bucket, nil
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	if bucket, exists = ms.buckets[clientID]; exists {
		return bucket, nil
	}

	var capacity, rate int64
	clientConf, confExists := ms.clientConfigs[clientID]
	if confExists {
		capacity = clientConf.Capacity
		rate = clientConf.Rate

		ms.logger.Debug("Using specific config for client",
			slog.String("client_id", clientID))

	} else {
		capacity = ms.defaultConfig.Capacity
		rate = ms.defaultConfig.Rate

		ms.logger.Debug("Using default config for client",
			slog.String("client_id", clientID))
	}

	if capacity <= 0 || rate <= 0 {
		ms.logger.Warn("Client config has invalid capacity or rate, falling back to store defaults",
			slog.String("client_id", clientID),
			slog.Int64("config_capacity", capacity), slog.Int64("config_rate", rate),
			slog.Int64("default_capacity", ms.defaultConfig.Capacity), slog.Int64("default_rate", ms.defaultConfig.Rate))

		capacity = ms.defaultConfig.Capacity
		rate = ms.defaultConfig.Rate
	}

	newBucket := token_bucket.NewTokenBucket(capacity, rate)
	ms.buckets[clientID] = newBucket

	ms.logger.Info("Created new token bucket",
		slog.String("client_id", clientID),
		slog.Int64("capacity", capacity),
		slog.Int64("rate", rate))

	return newBucket, nil
}
