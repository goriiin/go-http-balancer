package configs

import (
	"fmt"
	"github.com/spf13/viper"
)

var defaultConfig AppConfig

func init() {
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")

	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("не удалось прочитать  %w", err))
	}

	if err := viper.Sub("default").Unmarshal(&defaultConfig); err != nil {
		panic(fmt.Errorf("не удалось распарсить default: %w", err))
	}
}

func GetAppConfig(key string) (AppConfig, error) {
	cfg := defaultConfig

	sub := viper.Sub(key)
	if sub == nil {
		return cfg, fmt.Errorf("%q section not found — default returned", key)
	}

	if err := sub.Unmarshal(&cfg); err != nil {
		return defaultConfig, fmt.Errorf("section parsing error %q: %w — default returned", key, err)
	}

	return cfg, nil
}
