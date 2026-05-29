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

func (r *PostgresRepo) UpdateVideoStatus(ctx context.Context, videoID string, from blogger.VideoStatus, to blogger.VideoStatus) error {
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

func (r *PostgresRepo) GetVideoByUrl(ctx context.Context, url string) (*blogger.Video, error) {
	const q = `
		SELECT
			id, blogger_id, external_id, url, title,
			views, likes, comments,
			status, error_stage, error_message,
			published_at, created_at
		FROM videos
		WHERE url = $1
		LIMIT 1
	`

	v, err := scanVideo(r.db.QueryRow(ctx, q, url))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, blogger.ErrVideoNotFound
		}
		return nil, fmt.Errorf("get video by url: %w", err)
	}

	return v, nil
}

func (r *PostgresRepo) GetVideoByID(ctx context.Context, id string) (*blogger.Video, error) {
	const q = `
		SELECT
			id, blogger_id, external_id, url, title,
			views, likes, comments,
			status, error_stage, error_message,
			published_at, created_at
		FROM videos
		WHERE id = $1
		LIMIT 1
	`

	v, err := scanVideo(r.db.QueryRow(ctx, q, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, blogger.ErrVideoNotFound
		}
		return nil, fmt.Errorf("get video by id: %w", err)
	}

	return v, nil
}

func (r *PostgresRepo) UpdateVideoState(ctx context.Context, v *blogger.Video) error {
	const q = `
		UPDATE videos
		SET
			status        = $1,
			error_stage   = $2,
			error_message = $3
		WHERE id = $4
	`

	_, err := r.db.Exec(
		ctx,
		q,
		v.Status,
		v.ErrorStage,
		v.ErrorMessage,
		v.ID,
	)
	if err != nil {
		return fmt.Errorf("update video state: %w", err)
	}

	return nil
}

func (r *PostgresRepo) SaveVideoAnalysis(ctx context.Context, va *blogger.VideoAnalysis) error {
	const q = `
		INSERT INTO video_analysis (id, video_id, provider, raw_payload, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, q, va.ID, va.VideoID, va.Provider, va.RawPayload, va.CreatedAt)
	if err != nil {
		return fmt.Errorf("save video analysis: %w", err)
	}
	return nil
}

func (r *PostgresRepo) GetVideoAnalysisByVideoID(ctx context.Context, videoID string) (*blogger.VideoAnalysis, error) {
	const q = `
		SELECT id, video_id, provider, raw_payload, created_at
		FROM video_analysis
		WHERE video_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	var va blogger.VideoAnalysis
	err := r.db.QueryRow(ctx, q, videoID).Scan(&va.ID, &va.VideoID, &va.Provider, &va.RawPayload, &va.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, blogger.ErrVideoNotFound
		}
		return nil, fmt.Errorf("get video analysis by video id: %w", err)
	}
	return &va, nil
}

func (r *PostgresRepo) SaveVideoPrompt(ctx context.Context, vp *blogger.VideoPrompt) error {
	const q = `
		INSERT INTO video_prompts (id, video_id, llm_provider, llm_model, brief, script, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, q, vp.ID, vp.VideoID, vp.LLMProvider, vp.LLMModel, vp.Brief, vp.Script, vp.CreatedAt)
	if err != nil {
		return fmt.Errorf("save video prompt: %w", err)
	}
	return nil
}

func (r *PostgresRepo) GetVideoPromptByVideoID(ctx context.Context, videoID string) (*blogger.VideoPrompt, error) {
	const q = `
		SELECT id, video_id, llm_provider, llm_model, brief, script, created_at
		FROM video_prompts
		WHERE video_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	var vp blogger.VideoPrompt
	err := r.db.QueryRow(ctx, q, videoID).Scan(
		&vp.ID, &vp.VideoID, &vp.LLMProvider, &vp.LLMModel, &vp.Brief, &vp.Script, &vp.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, blogger.ErrVideoNotFound
		}
		return nil, fmt.Errorf("get video prompt by video id: %w", err)
	}
	return &vp, nil
}

func (r *PostgresRepo) SaveVideoGeneration(ctx context.Context, vg *blogger.VideoGeneration) error {
	const q = `
		INSERT INTO video_generations (id, video_id, platform, external_id, status, s3_url, error_message, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, q,
		vg.ID, vg.VideoID, vg.Platform, vg.ExternalID, vg.Status, vg.S3URL, vg.ErrorMessage, vg.CreatedAt, vg.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("save video generation: %w", err)
	}
	return nil
}

func (r *PostgresRepo) GetVideoGenerationByExternalID(ctx context.Context, externalID string) (*blogger.VideoGeneration, error) {
	const q = `
		SELECT id, video_id, platform, external_id, status, s3_url, error_message, created_at, updated_at
		FROM video_generations
		WHERE external_id = $1
		LIMIT 1
	`
	var vg blogger.VideoGeneration
	err := r.db.QueryRow(ctx, q, externalID).Scan(
		&vg.ID, &vg.VideoID, &vg.Platform, &vg.ExternalID, &vg.Status, &vg.S3URL, &vg.ErrorMessage, &vg.CreatedAt, &vg.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, blogger.ErrVideoNotFound
		}
		return nil, fmt.Errorf("get video generation by external id: %w", err)
	}
	return &vg, nil
}

func (r *PostgresRepo) ListPendingVideoGenerations(ctx context.Context) ([]*blogger.VideoGeneration, error) {
	const q = `
		SELECT id, video_id, platform, external_id, status, s3_url, error_message, created_at, updated_at
		FROM video_generations
		WHERE status IN ('pending', 'processing')
		ORDER BY created_at ASC
	`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list pending video generations: %w", err)
	}
	defer rows.Close()

	var result []*blogger.VideoGeneration
	for rows.Next() {
		var vg blogger.VideoGeneration
		if err := rows.Scan(
			&vg.ID, &vg.VideoID, &vg.Platform, &vg.ExternalID, &vg.Status, &vg.S3URL, &vg.ErrorMessage, &vg.CreatedAt, &vg.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan video generation: %w", err)
		}
		result = append(result, &vg)
	}
	return result, rows.Err()
}

func (r *PostgresRepo) UpdateVideoGeneration(ctx context.Context, vg *blogger.VideoGeneration) error {
	const q = `
		UPDATE video_generations
		SET status = $1, s3_url = $2, error_message = $3, updated_at = $4
		WHERE id = $5
	`
	_, err := r.db.Exec(ctx, q, vg.Status, vg.S3URL, vg.ErrorMessage, vg.UpdatedAt, vg.ID)
	if err != nil {
		return fmt.Errorf("update video generation: %w", err)
	}
	return nil
}

func scanVideo(row pgx.Row) (*blogger.Video, error) {
	var v blogger.Video
	var errorStage *string
	var errorMessage *string

	err := row.Scan(
		&v.ID,
		&v.BloggerID,
		&v.ExternalID,
		&v.URL,
		&v.Title,
		&v.Views,
		&v.Likes,
		&v.Comments,
		&v.Status,
		&errorStage,
		&errorMessage,
		&v.PublishedAt,
		&v.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if errorStage != nil {
		es := blogger.VideoErrorStage(*errorStage)
		v.ErrorStage = &es
	}

	v.ErrorMessage = errorMessage

	return &v, nil
}
func (r *PostgresRepo) ExistsVideoWatcher(ctx context.Context, videoID string, chatID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM video_watchers WHERE video_id = $1 AND chat_id = $2)`,
		videoID, chatID,
	).Scan(&exists)
	return exists, err
}

func (r *PostgresRepo) SaveVideoWatcher(ctx context.Context, w *blogger.VideoWatcher) error {
	const q = `
		INSERT INTO video_watchers (id, video_id, chat_id, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(ctx, q, w.ID, w.VideoID, w.ChatID, w.CreatedAt)
	if err != nil {
		return fmt.Errorf("save video watcher: %w", err)
	}
	return nil
}

func (r *PostgresRepo) ListVideoWatchers(ctx context.Context, videoID string) ([]*blogger.VideoWatcher, error) {
	const q = `
		SELECT id, video_id, chat_id, created_at
		FROM video_watchers
		WHERE video_id = $1
	`
	rows, err := r.db.Query(ctx, q, videoID)
	if err != nil {
		return nil, fmt.Errorf("list video watchers: %w", err)
	}
	defer rows.Close()

	var result []*blogger.VideoWatcher
	for rows.Next() {
		var w blogger.VideoWatcher
		if err := rows.Scan(&w.ID, &w.VideoID, &w.ChatID, &w.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan video watcher: %w", err)
		}
		result = append(result, &w)
	}
	return result, rows.Err()
}

func (r *PostgresRepo) SaveBrollSegments(ctx context.Context, videoID string, segments []*blogger.BrollSegment) error {
	const q = `
		INSERT INTO video_broll_segments (id, video_id, position, start_ms, end_ms, text, broll_prompt, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	for _, s := range segments {
		s.VideoID = videoID
		if _, err := r.db.Exec(ctx, q, s.ID, s.VideoID, s.Position, s.StartMS, s.EndMS, s.Text, s.BrollPrompt, s.CreatedAt); err != nil {
			return fmt.Errorf("save broll segment: %w", err)
		}
	}
	return nil
}

func (r *PostgresRepo) ListBrollSegmentsByVideoID(ctx context.Context, videoID string) ([]*blogger.BrollSegment, error) {
	const q = `
		SELECT id, video_id, position, start_ms, end_ms, text, broll_prompt, broll_url,
		       generation_status, generation_external_id, generation_error, created_at
		FROM video_broll_segments
		WHERE video_id = $1
		ORDER BY position
	`
	rows, err := r.db.Query(ctx, q, videoID)
	if err != nil {
		return nil, fmt.Errorf("list broll segments: %w", err)
	}
	defer rows.Close()

	var result []*blogger.BrollSegment
	for rows.Next() {
		s, err := scanBrollSegment(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, rows.Err()
}

func (r *PostgresRepo) ListPendingBrollSegments(ctx context.Context, videoID string) ([]*blogger.BrollSegment, error) {
	const q = `
		SELECT id, video_id, position, start_ms, end_ms, text, broll_prompt, broll_url,
		       generation_status, generation_external_id, generation_error, created_at
		FROM video_broll_segments
		WHERE video_id = $1 AND generation_status = 'pending'
		ORDER BY position
	`
	rows, err := r.db.Query(ctx, q, videoID)
	if err != nil {
		return nil, fmt.Errorf("list pending broll segments: %w", err)
	}
	defer rows.Close()

	var result []*blogger.BrollSegment
	for rows.Next() {
		s, err := scanBrollSegment(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, rows.Err()
}

func (r *PostgresRepo) ListVideosReadyToCompose(ctx context.Context) ([]string, error) {
	const q = `
		SELECT vg.video_id
		FROM video_generations vg
		WHERE vg.status = 'completed'
		  AND vg.s3_url != ''
		  AND NOT EXISTS (
		    SELECT 1 FROM video_broll_segments vbs
		    WHERE vbs.video_id = vg.video_id
		      AND vbs.generation_status IN ('pending', 'processing')
		  )
		  AND NOT EXISTS (
		    SELECT 1 FROM video_compositions vc
		    WHERE vc.video_id = vg.video_id
		  )
	`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list videos ready to compose: %w", err)
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		result = append(result, id)
	}
	return result, rows.Err()
}

func (r *PostgresRepo) ListVideosWithPendingBrollSegments(ctx context.Context) ([]string, error) {
	const q = `SELECT DISTINCT video_id FROM video_broll_segments WHERE generation_status = 'pending'`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list videos with pending broll segments: %w", err)
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		result = append(result, id)
	}
	return result, rows.Err()
}

func (r *PostgresRepo) ListProcessingBrollSegments(ctx context.Context) ([]*blogger.BrollSegment, error) {
	const q = `
		SELECT id, video_id, position, start_ms, end_ms, text, broll_prompt, broll_url,
		       generation_status, generation_external_id, generation_error, created_at
		FROM video_broll_segments
		WHERE generation_status = 'processing'
		ORDER BY created_at
		LIMIT 20
	`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list processing broll segments: %w", err)
	}
	defer rows.Close()

	var result []*blogger.BrollSegment
	for rows.Next() {
		s, err := scanBrollSegment(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, rows.Err()
}

func (r *PostgresRepo) UpdateBrollSegment(ctx context.Context, s *blogger.BrollSegment) error {
	const q = `
		UPDATE video_broll_segments
		SET generation_status = $1, generation_external_id = $2, generation_error = $3, broll_url = $4
		WHERE id = $5
	`
	_, err := r.db.Exec(ctx, q, s.GenerationStatus, s.GenerationExternalID, s.GenerationError, s.BrollURL, s.ID)
	if err != nil {
		return fmt.Errorf("update broll segment: %w", err)
	}
	return nil
}

func (r *PostgresRepo) CountActiveBrollSegments(ctx context.Context, videoID string) (int, error) {
	const q = `
		SELECT COUNT(*) FROM video_broll_segments
		WHERE video_id = $1 AND generation_status IN ('pending', 'processing')
	`
	var count int
	if err := r.db.QueryRow(ctx, q, videoID).Scan(&count); err != nil {
		return 0, fmt.Errorf("count active broll segments: %w", err)
	}
	return count, nil
}

func (r *PostgresRepo) GetVideoGenerationByVideoID(ctx context.Context, videoID string) (*blogger.VideoGeneration, error) {
	const q = `
		SELECT id, video_id, external_id, status, s3_url, error_message, created_at
		FROM video_generations
		WHERE video_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	var vg blogger.VideoGeneration
	err := r.db.QueryRow(ctx, q, videoID).Scan(
		&vg.ID, &vg.VideoID, &vg.ExternalID, &vg.Status, &vg.S3URL, &vg.ErrorMessage, &vg.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get video generation by video id: %w", err)
	}
	return &vg, nil
}

func (r *PostgresRepo) SaveVideoComposition(ctx context.Context, c *blogger.VideoComposition) error {
	const q = `
		INSERT INTO video_compositions (id, video_id, status, external_id, result_url, error, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, q, c.ID, c.VideoID, c.Status, c.ExternalID, c.ResultURL, c.Error, c.CreatedAt)
	if err != nil {
		return fmt.Errorf("save video composition: %w", err)
	}
	return nil
}

func (r *PostgresRepo) UpdateVideoComposition(ctx context.Context, c *blogger.VideoComposition) error {
	const q = `
		UPDATE video_compositions
		SET status = $1, result_url = $2, error = $3
		WHERE id = $4
	`
	_, err := r.db.Exec(ctx, q, c.Status, c.ResultURL, c.Error, c.ID)
	if err != nil {
		return fmt.Errorf("update video composition: %w", err)
	}
	return nil
}

func (r *PostgresRepo) ListProcessingCompositions(ctx context.Context) ([]*blogger.VideoComposition, error) {
	const q = `
		SELECT id, video_id, status, external_id, result_url, error, created_at
		FROM video_compositions
		WHERE status = 'processing'
		ORDER BY created_at
	`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list processing compositions: %w", err)
	}
	defer rows.Close()

	var result []*blogger.VideoComposition
	for rows.Next() {
		var c blogger.VideoComposition
		if err := rows.Scan(&c.ID, &c.VideoID, &c.Status, &c.ExternalID, &c.ResultURL, &c.Error, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan video composition: %w", err)
		}
		result = append(result, &c)
	}
	return result, rows.Err()
}

func scanBrollSegment(row interface{ Scan(...any) error }) (*blogger.BrollSegment, error) {
	var s blogger.BrollSegment
	err := row.Scan(
		&s.ID, &s.VideoID, &s.Position, &s.StartMS, &s.EndMS,
		&s.Text, &s.BrollPrompt, &s.BrollURL,
		&s.GenerationStatus, &s.GenerationExternalID, &s.GenerationError,
		&s.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan broll segment: %w", err)
	}
	return &s, nil
}

func (r *PostgresRepo) DeleteBrollSegmentsByVideoID(ctx context.Context, videoID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM video_broll_segments WHERE video_id = $1`, videoID)
	return err
}

func (r *PostgresRepo) DeleteCompositionsByVideoID(ctx context.Context, videoID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM video_compositions WHERE video_id = $1`, videoID)
	return err
}

func (r *PostgresRepo) DeleteVideoGenerationsByVideoID(ctx context.Context, videoID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM video_generations WHERE video_id = $1`, videoID)
	return err
}

func (r *PostgresRepo) DeleteVideoWatchersByVideoID(ctx context.Context, videoID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM video_watchers WHERE video_id = $1`, videoID)
	return err
}
