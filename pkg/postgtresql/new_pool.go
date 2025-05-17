package postgtresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

type pgConfig struct {
	Username string
	Password string
	DBName   string
}

func NewPool(dbc DBConnConfig) (*pgxpool.Pool, error) {
	cfg, err := getEnvConfig(dbc)
	if err != nil {
		return nil, fmt.Errorf("environment config error: %w", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s",
		cfg.Username,
		cfg.Password,
		dbc.Host,
		cfg.DBName,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("pool config parse error: %w", err)
	}

	poolCfg.MaxConns = int32(dbc.MaxConnections)
	poolCfg.MinConns = int32(dbc.MinConnections)
	poolCfg.MaxConnLifetime = dbc.MaxLifetime
	poolCfg.MaxConnIdleTime = dbc.MaxIdle

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, fmt.Errorf("connection pool creation error: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database ping error: %w", err)
	}

	return pool, nil
}

func getEnvConfig(dbc DBConnConfig) (pgConfig, error) {
	username := os.Getenv(dbc.UsernameKey)
	password := os.Getenv(dbc.PasswordKey)
	dbname := os.Getenv(dbc.NameKey)

	// Проверяем обязательные поля
	if username == "" || password == "" || dbname == "" {
		return pgConfig{}, fmt.Errorf("missing required environment variables")
	}

	return pgConfig{
		Username: username,
		Password: password,
		DBName:   dbname,
	}, nil
}
