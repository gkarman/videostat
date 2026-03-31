package platform

import (
	"context"
	"time"

	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/infrastructure/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(parent context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(parent, 10*time.Second)
	defer cancel()

	return db.NewPool(ctx, db.Config{
		DSN:             cfg.DB.DSN(),
		MaxConns:        cfg.DB.MaxConnections,
		MinConns:        cfg.DB.MinConnections,
		MaxConnLifetime: time.Duration(cfg.DB.MaxConnectionLifeTimeMinutes) * time.Minute,
		MaxConnIdleTime: time.Duration(cfg.DB.MaxConnectionIdleTimeMinutes) * time.Minute,
	})
	
}