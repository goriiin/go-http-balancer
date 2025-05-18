package configs

import "time"

type RateLimiterConfig struct {
	DefaultCapacity            int64         `mapstructure:"default_capacity"`
	DefaultRefillRatePerSecond int64         `mapstructure:"default_refill_rate_per_second"`
	ClientIDHeader             string        `mapstructure:"client_id_header"`
	TrustXForwardedFor         bool          `mapstructure:"trust_x_forwarded_for"`
	TrustedProxies             []string      `mapstructure:"trusted_proxies"`
	CleanupInterval            time.Duration `mapstructure:"cleanup_interval"`
	InactiveThreshold          time.Duration `mapstructure:"inactive_threshold"`
}
