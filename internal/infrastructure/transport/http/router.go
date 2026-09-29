package http

import (
	"log/slog"

	blogger_cmd "github.com/gkarman/demo/internal/application/blogger/command"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	blogger_repo "github.com/gkarman/demo/internal/infrastructure/repository/blogger"
	"github.com/gkarman/demo/internal/infrastructure/transport/http/handler"
	blogger_handler "github.com/gkarman/demo/internal/infrastructure/transport/http/handler/blogger"
	middleware2 "github.com/gkarman/demo/internal/infrastructure/transport/http/middleware"
	sharedapify "github.com/gkarman/demo/internal/infrastructure/apify"
	videoapify "github.com/gkarman/demo/internal/infrastructure/videosearcher/apify"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(log *slog.Logger, db *pgxpool.Pool, d *dispatcher.Dispatcher, apify *sharedapify.Client) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware2.Logger(log))
	r.Use(middleware2.Recovery())
	registerHomeRoutes(r)
	registerMetricsRoutes(r)
	registerVideoRoutes(r, db, d, apify)
	return r
}

func registerHomeRoutes(r *chi.Mux) {
	homeHandler := handler.NewHomeHandler()
	r.Get("/", homeHandler.Home)
}

// registerMetricsRoutes отдаёт метрики в формате Prometheus: он сам приходит сюда за ними.
func registerMetricsRoutes(r *chi.Mux) {
	r.Handle("/metrics", promhttp.Handler())
}

func registerVideoRoutes(r *chi.Mux, db *pgxpool.Pool, _ *dispatcher.Dispatcher, a *sharedapify.Client) {
	repoBlogger := blogger_repo.NewPostgres(db)
	videoSearcher := videoapify.NewVideoSearcher(a)
	fetchVideoCmd := blogger_cmd.NewFetchBloggerVideos(repoBlogger, videoSearcher)
	fetchVideoHandler := blogger_handler.NewGetCarHandler(fetchVideoCmd)

	r.Route("/test", func(r chi.Router) {
		r.Get("/videos-apify/{id}", fetchVideoHandler.Handle)
	})
}
