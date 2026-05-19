package core

import (
	"log/slog"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/application/blogger/command"
	blogger_handlers "github.com/gkarman/demo/internal/application/blogger/handlers"
	"github.com/gkarman/demo/internal/domain/blogger"
	sharedapify "github.com/gkarman/demo/internal/infrastructure/apify"
	"github.com/gkarman/demo/internal/infrastructure/contracts/events"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	blogger_repo "github.com/gkarman/demo/internal/infrastructure/repository/blogger"
	apifysearcher "github.com/gkarman/demo/internal/infrastructure/videosearcher/apify"
	apifysource "github.com/gkarman/demo/internal/infrastructure/videosourcesearcher/apify"
	"github.com/gkarman/demo/internal/worker"
	"github.com/gkarman/demo/internal/worker/core/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouterWithHandlers(
	log *slog.Logger,
	db *pgxpool.Pool,
	apifyClient *sharedapify.Client,
	analyzer application.VideoAnalyzer,
	videoGenerator application.VideoGenerator,
	publisher application.Publisher,
) *worker.Router {
	r := worker.NewRouter(log)

	bRepo := blogger_repo.NewPostgres(db)

	disp := dispatcher.New()
	disp.Register(&blogger.VideoSourceFound{}, blogger_handlers.VideoSourceFoundToRabbitHandler(publisher, log))
	disp.Register(&blogger.VideoAnalyzeDone{}, blogger_handlers.VideoAnalyzeDoneToRabbitHandler(publisher, log))

	vSearcher := apifysearcher.NewVideoSearcher(apifyClient)
	cmdFVideos := command.NewFetchBloggerVideos(bRepo, vSearcher)

	vSourceSearcher := apifysource.NewVideoSourceSearcher(apifyClient)
	cmdFSources := command.NewFetchVideoSources(bRepo, vSourceSearcher, disp)

	cmdAnalyze := command.NewAnalyzeVideo(bRepo, analyzer, disp)
	cmdSubmitGeneration := command.NewSubmitVideoGeneration(bRepo, videoGenerator)

	bloggerCreatedHandler := handlers.NewBloggerCreatedHandler(log, cmdFVideos)
	videoProcessingStartedHandler := handlers.NewVideoProcessingStarted(log, cmdFSources)
	videoSourceFoundHandler := handlers.NewVideoSourceFoundHandler(log, cmdAnalyze)
	videoAnalyzeDoneHandler := handlers.NewVideoAnalyzeDoneHandler(log, cmdSubmitGeneration)

	r.Register(events.EventBloggerCreatedV1, bloggerCreatedHandler.Handle)
	r.Register(events.EventVideoProcessingStartedV1, videoProcessingStartedHandler.Handle)
	r.Register(events.EventVideoSourceFoundV1, videoSourceFoundHandler.Handle)
	r.Register(events.EventVideoAnalyzeDoneV1, videoAnalyzeDoneHandler.Handle)
	return r
}
