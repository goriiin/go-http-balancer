package store

func (ms *MemoryStore) Stop() {
	if ms.cleanupTicker != nil {
		select {
		case <-ms.stopCleanup:
			return
		default:
			close(ms.stopCleanup)
		}
	}
}
