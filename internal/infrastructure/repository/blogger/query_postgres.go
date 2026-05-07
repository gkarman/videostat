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

func (r *QueryPostgres) ListBloggers(ctx context.Context) ([]blogger.BloggerRow, error) {
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

func (r *QueryPostgres) ListVideos(ctx context.Context) ([]blogger.VideoRow, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
		    v.id,
		    p.name            AS platform_name,
		    b.url             AS blogger_url,
		    v.url,
		    v.title,
		    v.views,
		    v.likes,
		    v.comments,
		    v.published_at,
		    v.created_at,
		    v.status,
		    v.error_stage,
		    v.error_message,
		    va.provider       AS analysis_provider,
		    va.raw_payload::text,
		    vp.brief,
		    vp.script,
		    vp.llm_provider,
		    vp.llm_model,
		    vg.platform        AS generation_platform,
		    vg.external_id     AS generation_external_id,
		    vg.status          AS generation_status,
		    vg.s3_url          AS generation_s3_url,
		    vg.error_message   AS generation_error_message
		FROM videos v
		JOIN bloggers b ON b.id = v.blogger_id
		JOIN platforms p ON p.id = b.platform_id
		LEFT JOIN LATERAL (
		    SELECT provider, raw_payload
		    FROM video_analysis
		    WHERE video_id = v.id
		    ORDER BY created_at DESC
		    LIMIT 1
		) va ON true
		LEFT JOIN LATERAL (
		    SELECT brief, script, llm_provider, llm_model
		    FROM video_prompts
		    WHERE video_id = v.id
		    ORDER BY created_at DESC
		    LIMIT 1
		) vp ON true
		LEFT JOIN LATERAL (
		    SELECT platform, external_id, status, s3_url, error_message
		    FROM video_generations
		    WHERE video_id = v.id
		    ORDER BY created_at DESC
		    LIMIT 1
		) vg ON true
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []blogger.VideoRow

	for rows.Next() {
		var v blogger.VideoRow
		if err := rows.Scan(
			&v.ID,
			&v.Platform,
			&v.BloggerURL,
			&v.URL,
			&v.Title,
			&v.Views,
			&v.Likes,
			&v.Comments,
			&v.PublishedAt,
			&v.CreatedAt,
			&v.Status,
			&v.ErrorStage,
			&v.ErrorMessage,
			&v.AnalysisProvider,
			&v.RawPayload,
			&v.Brief,
			&v.Script,
			&v.LLMProvider,
			&v.LLMModel,
			&v.GenerationPlatform,
			&v.GenerationExternalID,
			&v.GenerationStatus,
			&v.GenerationS3URL,
			&v.GenerationErrorMessage,
		); err != nil {
			return nil, err
		}
		result = append(result, v)
	}

	return result, rows.Err()
}