package health_checker

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

func (hc *HealthChecker) Start(ctx context.Context) {
	ticker := time.NewTicker(hc.Period)
	client := &http.Client{Timeout: hc.Timeout}

	for {
		select {
		case <-ticker.C:
			backs := hc.Pool.GetAll()
			for _, b := range backs {
				probe := b.URL.ResolveReference(&url.URL{Path: b.URL.Path})
				resp, err := client.Head(probe.String())

				if err != nil || resp.StatusCode >= 400 {
					hc.Pool.MarkUnhealthy(b)
				} else {
					hc.Pool.MarkHealthy(b)
				}
			}

		case <-ctx.Done():
			ticker.Stop()

			return
		}
	}
}
