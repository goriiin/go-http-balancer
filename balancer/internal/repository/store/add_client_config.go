package store

import (
	"errors"
	"fmt"
	"github.com/goriiin/go-http-balancer/balancer/internal/domain"
	"log/slog"
)

func (ms *MemoryStore) AddClientConfig(config domain.ClientInfo) error {
	if config.ClientID == "" {
		return errors.New("client ID cannot be empty")
	}

	if config.Capacity <= 0 || config.Rate <= 0 {
		return fmt.Errorf("capacity and rate must be positive for client %s", config.ClientID)
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.clientConfigs[config.ClientID]; exists {
		return fmt.Errorf("configuration for client ID %s already exists, use UpdateClientConfig", config.ClientID)
	}

	ms.clientConfigs[config.ClientID] = config
	if _, bucketExists := ms.buckets[config.ClientID]; bucketExists {
		delete(ms.buckets, config.ClientID)

		ms.logger.Info("Removed existing bucket for client due to new config addition",
			slog.String("client_id", config.ClientID))
	}

	ms.logger.Info("Added client config",
		slog.String("client_id", config.ClientID),
		slog.Int64("capacity", config.Capacity),
		slog.Int64("rate", config.Rate))

	return nil
}
