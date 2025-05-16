package logger

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func load(project string) (LogConfig, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		return LogConfig{}, err
	}

	logSub := v.Sub("log")
	if logSub == nil {
		return LogConfig{}, errors.New("no log config")
	}

	projSub := logSub.Sub("project")
	if projSub == nil {
		return LogConfig{}, errors.New("no log project")
	}

	var cfg LogConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		_ = v.Unmarshal(&cfg)
	})

	return cfg, nil
}
