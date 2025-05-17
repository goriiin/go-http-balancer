package postgtresql

import "time"

type DBConnConfig struct {
	URI string `mapstructure:"uri"`

	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	UsernameKey string `mapstructure:"username_key"`
	PasswordKey string `mapstructure:"password_key"`
	NameKey     string `mapstructure:"name_key"`

	MaxConnections int           `mapstructure:"max_connections"`
	MinConnections int           `mapstructure:"min_connections"`
	MaxLifetime    time.Duration `mapstructure:"max_lifetime"`
	MaxIdle        time.Duration `mapstructure:"max_idle"`
}
