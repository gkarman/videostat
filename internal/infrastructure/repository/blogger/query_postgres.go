package blogger

import (
	"context"

	"github.com/gkarman/demo/internal/application/blogger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QueryPostgres struct {
	db *pgxpool.Pool
}

func NewQueryPostgres(db *pgxpool.Pool) *QueryPostgres {
	return &QueryPostgres{db: db}
}

func (r *QueryPostgres) List(ctx context.Context) ([]blogger.BloggerRow, error) {
	rows, err := r.db.Query(ctx, `
		SELECT b.id, b.url, p.name
		FROM bloggers b
		JOIN platforms p ON p.id = b.platform_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []blogger.BloggerRow

	for rows.Next() {
		var v blogger.BloggerRow
		if err := rows.Scan(&v.ID, &v.URL, &v.Platform); err != nil {
			return nil, err
		}
		result = append(result, v)
	}

	return result, rows.Err()
}