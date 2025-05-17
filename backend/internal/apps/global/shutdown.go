package global

import "context"

func (a *App) Shutdown() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		a.config.GracefulShutdownTimeout,
	)
	defer cancel()

	a.log.Info("shutting down server")

	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Error("server shutdown error", "error", err)
	}

	if closer, ok := a.dataDelivery.(interface{ Close() error }); ok {
		if err := closer.Close(); err != nil {
			a.log.Error("database connection closure error", "error", err)
		}
	}

	a.log.Info("server stopped")
}
