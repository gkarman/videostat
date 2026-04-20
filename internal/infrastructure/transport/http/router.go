package http

import (
	"log/slog"

	blogger_cmd "github.com/gkarman/demo/internal/application/blogger/command"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	blogger_repo "github.com/gkarman/demo/internal/infrastructure/repository/blogger"
	"github.com/gkarman/demo/internal/infrastructure/transport/http/handler"
	blogger_handler "github.com/gkarman/demo/internal/infrastructure/transport/http/handler/blogger"
	middleware2 "github.com/gkarman/demo/internal/infrastructure/transport/http/middleware"
	"github.com/gkarman/demo/internal/infrastructure/videosearcher/apify"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(log *slog.Logger, db *pgxpool.Pool, d *dispatcher.Dispatcher, apify *apify.Client) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware2.Logger(log))
	r.Use(middleware2.Recovery())
	registerHomeRoutes(r)
	registerVideoRoutes(r, db, d, apify)
	return r
}

func registerHomeRoutes(r *chi.Mux) {
	homeHandler := handler.NewHomeHandler()
	r.Get("/", homeHandler.Home)
}


func registerVideoRoutes(r *chi.Mux, db *pgxpool.Pool, _ *dispatcher.Dispatcher, a *apify.Client) {

	repoBlogger := blogger_repo.NewPostgres(db)
	videoSearcher := apify.NewVideoSearcher(a)
	fetchVideoCmd := blogger_cmd.NewFetchBloggerVideos(repoBlogger, videoSearcher)
	fetchVideoHandler := blogger_handler.NewGetCarHandler(fetchVideoCmd)

	r.Route("/test", func(r chi.Router) {
		r.Get("/videos-apify/{id}", fetchVideoHandler.Handle)
	})
}
