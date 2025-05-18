package store

import (
	"log/slog"
	"time"
)

func (ms *MemoryStore) performCleanup() {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	now := time.Now()
	cleanedCount := 0

	for clientID, bucket := range ms.buckets {
		if now.Sub(bucket.LastAccessed()) > ms.inactiveThreshold {
			delete(ms.buckets, clientID)
			cleanedCount++

			ms.logger.Debug("Cleaned up inactive rate limit bucket",
				slog.String("client_id", clientID))
		}
	}

	if cleanedCount > 0 {
		ms.logger.Info("Performed rate limit bucket cleanup",
			slog.Int("cleaned_count", cleanedCount),
			slog.Int("remaining_count", len(ms.buckets)))
	}
}
