package metrics

import (
	"context"
	"log/slog"
	"time"

	"github.com/gkarman/demo/internal/domain/blogger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
)

const collectTimeout = 5 * time.Second

var (
	videosDesc = prometheus.NewDesc(
		"videostat_videos",
		"Количество видео в каждом статусе.",
		[]string{"status"}, nil,
	)
	videosFailedDesc = prometheus.NewDesc(
		"videostat_videos_failed",
		"Количество видео, остановившихся с ошибкой, по этапу ошибки.",
		[]string{"stage"}, nil,
	)
)

var (
	knownStatuses = []blogger.VideoStatus{
		blogger.VideoStatusCreated,
		blogger.VideoStatusProcessing,
		blogger.VideoStatusGenerationProcessing,
		blogger.VideoStatusGenerationFailed,
		blogger.VideoStatusReady,
		blogger.VideoStatusFailed,
	}
	knownStages = []blogger.VideoErrorStage{
		blogger.ErrorStageFileFetch,
		blogger.ErrorStageAnalysis,
		blogger.ErrorStageInsights,
		blogger.ErrorStageGeneration,
	}
)

// VideoCollector считает видео прямо в БД в момент, когда Prometheus забирает /metrics.
// В отличие от счётчиков в коде, значения всегда совпадают с БД и не сбрасываются при рестарте.
// Регистрировать только в одном процессе, иначе одни и те же цифры придут от нескольких job.
type VideoCollector struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewVideoCollector(db *pgxpool.Pool, log *slog.Logger) *VideoCollector {
	return &VideoCollector{db: db, log: log}
}

func (c *VideoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- videosDesc
	ch <- videosFailedDesc
}

func (c *VideoCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), collectTimeout)
	defer cancel()

	// Если запрос упал — не отдаём метрику вовсе: лучше «нет данных», чем ложный ноль.
	if byStatus, err := c.countBy(ctx, `SELECT status, count(*) FROM videos GROUP BY status`); err != nil {
		c.log.Error("collect videos by status", "error", err)
	} else {
		// Заранее известные статусы отдаём и с нулём, чтобы на графике был 0, а не пустота.
		for _, s := range knownStatuses {
			ch <- prometheus.MustNewConstMetric(videosDesc, prometheus.GaugeValue, byStatus[string(s)], string(s))
		}
	}

	if byStage, err := c.countBy(ctx, `SELECT error_stage, count(*) FROM videos WHERE error_stage IS NOT NULL GROUP BY error_stage`); err != nil {
		c.log.Error("collect failed videos by stage", "error", err)
	} else {
		for _, s := range knownStages {
			ch <- prometheus.MustNewConstMetric(videosFailedDesc, prometheus.GaugeValue, byStage[string(s)], string(s))
		}
	}
}

func (c *VideoCollector) countBy(ctx context.Context, query string) (map[string]float64, error) {
	rows, err := c.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]float64)
	for rows.Next() {
		var key string
		var count int64
		if err := rows.Scan(&key, &count); err != nil {
			return nil, err
		}
		result[key] = float64(count)
	}
	return result, rows.Err()
}
