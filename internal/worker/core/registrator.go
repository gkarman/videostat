package core

import (
	"log/slog"

	"github.com/gkarman/demo/internal/application/blogger/command"
	"github.com/gkarman/demo/internal/infrastructure/contracts/events"
	"github.com/gkarman/demo/internal/infrastructure/repository/blogger"
	"github.com/gkarman/demo/internal/infrastructure/videosearcher/apify"
	"github.com/gkarman/demo/internal/worker"
	"github.com/gkarman/demo/internal/worker/core/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouterWithHandlers(log *slog.Logger, db *pgxpool.Pool, apifyClient *apify.Client) *worker.Router {
	r := worker.NewRouter(log)

	vSearcher := apify.NewVideoSearcher(apifyClient)
	bRepo := blogger.NewPostgres(db)
	cmdFVideos := command.NewFetchBloggerVideos(bRepo, vSearcher)

	bloggerCreatedHandler := handlers.NewBloggerCreatedHandler(log, cmdFVideos)

	r.Register(events.EventBloggerCreatedV1, bloggerCreatedHandler.Handle)
	return r
}