package blogger

import (
	"context"
	"fmt"

	"github.com/gkarman/demo/internal/domain/blogger"
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

func (r *PostgresRepo) Save(ctx context.Context, b *blogger.Blogger) error {
	const q = `
		INSERT INTO bloggers (id, platform_id, url) VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(
		ctx,
		q,
		b.ID,
		b.PlatformID,
		b.URL,
	)

	return err
}

func (r *PostgresRepo) ExistByUrl(ctx context.Context, url string) (bool, error) {
	const q = `
		SELECT EXISTS (
			SELECT 1
			FROM bloggers
			WHERE url = $1
		);
	
	`
	var exists bool
	if err := r.db.QueryRow(ctx, q, url).Scan(&exists); err != nil {
		return false, fmt.Errorf("check blogger exists by url: %w", err)
	}

	return exists, nil

}

