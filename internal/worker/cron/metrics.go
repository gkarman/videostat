package cron

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var jobDurationBuckets = []float64{0.1, 0.5, 1, 2.5, 5, 10, 30, 60, 120, 300, 600, 1800}

var (
	jobRunsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cron_job_runs_total",
			Help: "Количество запусков cron-задач.",
		},
		[]string{"job", "result"},
	)

	jobDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cron_job_duration_seconds",
			Help:    "Длительность выполнения cron-задачи в секундах.",
			Buckets: jobDurationBuckets,
		},
		[]string{"job"},
	)

	// Gauge с unix-временем последнего успеха. Позволяет поймать задачу,
	// которая молча перестала запускаться: time() - метрика растёт без остановки.
	jobLastSuccess = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cron_job_last_success_timestamp_seconds",
			Help: "Unix-время последнего успешного выполнения cron-задачи.",
		},
		[]string{"job"},
	)
)
