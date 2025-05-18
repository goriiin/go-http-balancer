package store

import (
	"errors"
	"fmt"
	"log/slog"
)

func (ms *MemoryStore) DeleteClientConfig(clientID string) error {
	if clientID == "" {
		return errors.New("client ID cannot be empty")
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.clientConfigs[clientID]; !exists {
		return fmt.Errorf("no configuration found for client ID %s to delete", clientID)
	}

	delete(ms.clientConfigs, clientID)
	if _, bucketExists := ms.buckets[clientID]; bucketExists {
		delete(ms.buckets, clientID)
		ms.logger.Info("Removed existing bucket for client due to config deletion",
			slog.String("client_id", clientID))
	}

	ms.logger.Info("Deleted client config",
		slog.String("client_id", clientID))

	return nil
}
