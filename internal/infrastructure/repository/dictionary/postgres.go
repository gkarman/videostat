package dictionary

import (
	"context"
	"errors"

	"github.com/gkarman/demo/internal/domain/dictionary"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgres(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}

func (r *PostgresRepo) GetPlatformByName(ctx context.Context, name string) (*dictionary.Platform, error) {
	const q = `
		SELECT id, name
		FROM platforms
		WHERE name = $1
	`

	row := r.db.QueryRow(ctx, q, name)

	var p dictionary.Platform
	if err := row.Scan(&p.ID, &p.Name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dictionary.ErrPlatformNotFound
		}
		return nil, err
	}

	return &p, nil
}