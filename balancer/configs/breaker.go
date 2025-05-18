package configs

import "time"

type BreakerConfig struct {
	FailureThreshold int           `mapstructure:"failure_threshold"`
	ResetTimeout     time.Duration `mapstructure:"reset_timeout"`
	HalfOpenRequests int           `mapstructure:"half_open_requests"`
}
