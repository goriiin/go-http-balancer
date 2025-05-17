package configs

import (
	"github.com/goriiin/go-http-balancer/pkg/logger"
	"github.com/goriiin/go-http-balancer/pkg/postgtresql"
	"time"
)

type AppConfig struct {
	GracefulShutdownTimeout time.Duration            `mapstructure:"graceful_shutdown_timeout"`
	ReadTimeout             time.Duration            `mapstructure:"read_timeout"`
	WriteTimeout            time.Duration            `mapstructure:"write_timeout"`
	IdleTimeout             time.Duration            `mapstructure:"idle_timeout"`
	Logger                  logger.Logger            `mapstructure:"logger"`
	DB                      postgtresql.DBConnConfig `mapstructure:"db"`
}
