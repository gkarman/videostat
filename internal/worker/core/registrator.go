package core

import (
	"log/slog"

	"github.com/gkarman/demo/internal/application/blogger/command"
	"github.com/gkarman/demo/internal/infrastructure/contracts/events"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	"github.com/gkarman/demo/internal/infrastructure/repository/blogger"
	apifysearcher "github.com/gkarman/demo/internal/infrastructure/videosearcher/apify"
	apifysource "github.com/gkarman/demo/internal/infrastructure/videosourcesearcher/apify"
	"github.com/gkarman/demo/internal/worker"
	"github.com/gkarman/demo/internal/worker/core/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouterWithHandlers(log *slog.Logger, db *pgxpool.Pool, apifyClient *apifysearcher.Client) *worker.Router {
	r := worker.NewRouter(log)

	bRepo := blogger.NewPostgres(db)
	disp := dispatcher.New()

	vSearcher := apifysearcher.NewVideoSearcher(apifyClient)
	cmdFVideos := command.NewFetchBloggerVideos(bRepo, vSearcher)

	vSourceSearcher := apifysource.NewVideoSourceSearcher(apifyClient)
	cmdFSources := command.NewFetchVideoSources(bRepo, vSourceSearcher, disp)

	bloggerCreatedHandler := handlers.NewBloggerCreatedHandler(log, cmdFVideos)
	videoProcessingStartedHandler := handlers.NewVideoProcessingStarted(log, cmdFSources)

	r.Register(events.EventBloggerCreatedV1, bloggerCreatedHandler.Handle)
	r.Register(events.EventVideoProcessingStartedV1, videoProcessingStartedHandler.Handle)
	return r
}