package store

func (ms *MemoryStore) startCleanupRoutine() {
	for {
		select {
		case <-ms.cleanupTicker.C:
			ms.performCleanup()
		case <-ms.stopCleanup:
			ms.cleanupTicker.Stop()
			ms.logger.Info("Rate limit bucket cleanup routine stopped.")

			return
		}
	}
}
