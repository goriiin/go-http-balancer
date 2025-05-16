package breaker

import "time"

type Config struct {
	FailureThreshold int
	ResetTimeout     time.Duration
	HalfOpenRequests int
}
