package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	DSN             string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

func NewPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	poolCfg, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("parse DSN: %w", err)
	}

	poolCfg.MaxConns = cfg.MaxConns
	poolCfg.MinConns = cfg.MinConns
	poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolCfg.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

func validateConfig(cfg Config) error {
	if cfg.DSN == "" {
		return fmt.Errorf("DSN is required")
	}
	if cfg.MaxConns <= 0 {
		return fmt.Errorf("MaxConns must be positive, got %d", cfg.MaxConns)
	}
	if cfg.MaxConns > 500 {
		return fmt.Errorf("MaxConns too large: %d", cfg.MaxConns)
	}
	if cfg.MinConns < 0 {
		return fmt.Errorf("MinConns must be non-negative, got %d", cfg.MinConns)
	}
	if cfg.MinConns > cfg.MaxConns {
		return fmt.Errorf("MinConns (%d) cannot be greater than MaxConns (%d)", cfg.MinConns, cfg.MaxConns)
	}
	return nil
}
