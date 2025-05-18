package configs

import (
	"github.com/goriiin/go-http-balancer/balancer/internal/usecase/health_checker"
	"github.com/goriiin/go-http-balancer/pkg/logger"
	"time"
)

type BalancerAppConfig struct {
	ListenAddr              string                `mapstructure:"listen_addr"`
	Algorithm               string                `mapstructure:"algorithm"`
	Backends                []string              `mapstructure:"backends"`
	HealthCheck             health_checker.Config `mapstructure:"health_check"`
	CircuitBreaker          BreakerConfig         `mapstructure:"circuit_breaker"`
	GracefulShutdownTimeout time.Duration         `mapstructure:"graceful_shutdown_timeout"`
	Logger                  logger.Logger         `mapstructure:"logger"`
	RateLimiter             RateLimiterConfig     `mapstructure:"rate_limiter"`
}

func DefaultBalancerCfg() BalancerAppConfig {
	return BalancerAppConfig{
		ListenAddr: ":8080",
		Algorithm:  "round-robin",
		Backends:   []string{},
		HealthCheck: health_checker.Config{
			Path:    "/api/v1/health",
			Period:  15 * time.Second,
			Timeout: 5 * time.Second,
		},
		CircuitBreaker: BreakerConfig{
			FailureThreshold: 5,
			ResetTimeout:     30 * time.Second,
			HalfOpenRequests: 2,
		},
		GracefulShutdownTimeout: 15 * time.Second,
		Logger: logger.Logger{
			Level: "info",
		},
		RateLimiter: RateLimiterConfig{
			DefaultCapacity:            100,
			DefaultRefillRatePerSecond: 10,
			ClientIDHeader:             "",
			TrustXForwardedFor:         false,
			TrustedProxies:             []string{},
			CleanupInterval:            5 * time.Minute,
			InactiveThreshold:          30 * time.Minute,
		},
	}
}
