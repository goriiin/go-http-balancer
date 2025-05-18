package health_checker

import (
	"context"
	"log/slog"
	"time"
)

func (hc *HealthChecker) Start(ctx context.Context) {
	if hc.config.Period <= 0 {
		hc.log.Info("health check period is not positive, checker will not run")

		return
	}
	ticker := time.NewTicker(hc.config.Period)
	defer ticker.Stop()

	hc.log.Info("health checker started",
		slog.String("path", hc.config.Path),
		slog.Duration("period", hc.config.Period),
		slog.Duration("timeout", hc.config.Timeout),
	)

	hc.performChecks()

	for {
		select {
		case <-ticker.C:
			hc.performChecks()
		case <-ctx.Done():
			hc.log.Info("health checker stopping due to context cancellation")
			return
		}
	}
}
