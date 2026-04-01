package video

import (
	"context"
	"errors"

	"github.com/gkarman/demo/internal/domain/video"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) SaveSnapshot(ctx context.Context, snap *video.AccountSnapshot) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	accountID, err := r.upsertAccount(ctx, tx, snap.Account)
	if err != nil {
		return err
	}

	if err := r.insertAccountStats(ctx, tx, accountID, snap.AccountStats); err != nil {
		return err
	}

	for _, c := range snap.Contents {
		contentID, err := r.upsertContent(ctx, tx, accountID, c.Content)
		if err != nil {
			return err
		}

		if err := r.insertContentStats(ctx, tx, contentID, c.Stats); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgresRepo) upsertAccount(ctx context.Context, tx pgx.Tx, a *video.Account) (int64, error) {
	const q = `
		INSERT INTO video_accounts (platform_id, external_id, title, url)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (platform_id, external_id)
		DO UPDATE SET title = EXCLUDED.title, url = EXCLUDED.url, updated_at = NOW()
		RETURNING id
	`

	var id int64
	err := tx.QueryRow(ctx, q,
		a.PlatformID,
		a.ExternalID,
		a.Title,
		a.URL,
	).Scan(&id)

	return id, err
}

func (r *PostgresRepo) upsertContent(ctx context.Context, tx pgx.Tx, accountID int64, c *video.Content) (int64, error) {
	const q = `
		INSERT INTO video_content (account_id, external_id, title, published_at, duration_sec)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (account_id, external_id)
		DO UPDATE SET title = EXCLUDED.title
		RETURNING id
	`

	var id int64
	err := tx.QueryRow(ctx, q,
		accountID,
		c.ExternalID,
		c.Title,
		c.PublishedAt,
		c.DurationSec,
	).Scan(&id)

	return id, err
}

func (r *PostgresRepo) insertAccountStats(ctx context.Context, tx pgx.Tx, accountID int64, s *video.AccountStats) error {
	const q = `
		INSERT INTO video_account_stats_history (account_id, followers, total_views, collected_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := tx.Exec(ctx, q,
		accountID,
		s.Followers,
		s.TotalViews,
		s.CollectedAt,
	)

	return err
}

func (r *PostgresRepo) insertContentStats(ctx context.Context, tx pgx.Tx, contentID int64, s *video.ContentStats) error {
	const q = `
		INSERT INTO video_content_stats_history (content_id, views, likes, comments, shares, collected_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := tx.Exec(ctx, q,
		contentID,
		s.Views,
		s.Likes,
		s.Comments,
		s.Shares,
		s.CollectedAt,
	)

	return err
}

func (r *PostgresRepo) ExistsByPlatformAndExternalID(ctx context.Context, platformID int16, externalID string) (bool, error) {
	const q = `
		SELECT 1 FROM video_accounts
		WHERE platform_id = $1 AND external_id = $2
		LIMIT 1
	`

	var tmp int
	err := r.db.QueryRow(ctx, q, platformID, externalID).Scan(&tmp)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
