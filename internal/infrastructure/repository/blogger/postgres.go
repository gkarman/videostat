package blogger

import (
	"context"
	"errors"
	"fmt"

	"github.com/gkarman/demo/internal/domain/blogger"
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

func (r *PostgresRepo) GetById(ctx context.Context, id string) (*blogger.Blogger, error) {
	const q = `
		SELECT id, platform_id, url
		FROM bloggers
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, q, id)

	var b blogger.Blogger
	if err := row.Scan(&b.ID, &b.PlatformID, &b.URL); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, blogger.ErrBloggerNotFound
		}
		return nil, err
	}

	return &b, nil
}

func (r *PostgresRepo) SaveVideo(ctx context.Context, v *blogger.Video) error {
	const q = `
		INSERT INTO videos 
		(id, blogger_id, external_id, url, title, views, likes, comments, published_at, created_at, status)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
		    $11   
		)
		ON CONFLICT (blogger_id, external_id)
		DO UPDATE SET
			url          = EXCLUDED.url,
			title        = EXCLUDED.title,
			views        = EXCLUDED.views,
			likes        = EXCLUDED.likes,
			comments     = EXCLUDED.comments,
			published_at = EXCLUDED.published_at
`
	_, err := r.db.Exec(ctx,
		q,
		v.ID,
		v.BloggerID,
		v.ExternalID,
		v.URL,
		v.Title,
		v.Views,
		v.Likes,
		v.Comments,
		v.PublishedAt,
		v.CreatedAt,
		v.Status,
	)
	if err != nil {
		return fmt.Errorf("save video: %w", err)
	}
	return nil
}

func (r *PostgresRepo) ListVideosByBlogger(ctx context.Context, bloggerID string) ([]*blogger.Video, error) {
	const q = `
		SELECT 
			id,
			blogger_id,
			external_id,
			url,
			title,
			views,
			likes,
			comments,
			published_at,
			created_at
		FROM videos
		WHERE blogger_id = $1
		ORDER BY views DESC
	`

	rows, err := r.db.Query(ctx, q, bloggerID)
	if err != nil {
		return nil, fmt.Errorf("list videos by blogger: %w", err)
	}
	defer rows.Close()

	var result []*blogger.Video

	for rows.Next() {
		var v blogger.Video
		if err := rows.Scan(
			&v.ID,
			&v.BloggerID,
			&v.ExternalID,
			&v.URL,
			&v.Title,
			&v.Views,
			&v.Likes,
			&v.Comments,
			&v.PublishedAt,
			&v.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan video: %w", err)
		}

		result = append(result, &v)
	}

	return result, nil
}

func (r *PostgresRepo) List(ctx context.Context) ([]*blogger.Blogger, error) {
	return nil, nil
}

func (r *PostgresRepo) UpdateStatus(ctx context.Context, videoID string, from blogger.VideoStatus, to blogger.VideoStatus) error {
	const q = `
		UPDATE videos
		SET status = $1
		WHERE id = $2
		AND status = $3
	`

	result, err := r.db.Exec(ctx, q, to, videoID, from)
	if err != nil {
		return err
	}

	rows := result.RowsAffected()
	if rows == 0 {
		return blogger.ErrConcurrentUpdate
	}

	return nil
}
