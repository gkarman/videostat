package http

import (
	"log/slog"

	blogger_cmd "github.com/gkarman/demo/internal/application/blogger/command"
	"github.com/gkarman/demo/internal/application/car/service"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	blogger_repo "github.com/gkarman/demo/internal/infrastructure/repository/blogger"
	"github.com/gkarman/demo/internal/infrastructure/repository/car"
	"github.com/gkarman/demo/internal/infrastructure/transport/http/handler"
	blogger_handler "github.com/gkarman/demo/internal/infrastructure/transport/http/handler/blogger"
	car_handler "github.com/gkarman/demo/internal/infrastructure/transport/http/handler/car"
	middleware2 "github.com/gkarman/demo/internal/infrastructure/transport/http/middleware"
	"github.com/gkarman/demo/internal/infrastructure/videosearcher/apify"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(log *slog.Logger, db *pgxpool.Pool, d *dispatcher.Dispatcher) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware2.Logger(log))
	r.Use(middleware2.Recovery())
	registerHomeRoutes(r)
	registerCarRoutes(r, db, d)
	registerVideoRoutes(r, db, d)
	return r
}

func registerHomeRoutes(r *chi.Mux) {
	homeHandler := handler.NewHomeHandler()
	r.Get("/", homeHandler.Home)
}

func registerCarRoutes(r *chi.Mux, db *pgxpool.Pool, d *dispatcher.Dispatcher) {
	repo := car.NewPostgresRepo(db)

	listSvc := service.NewList(repo)
	listHandler := car_handler.NewList(listSvc)

	getCarSvc := service.NewGet(repo)
	getCarHandler := car_handler.NewGetCarHandler(getCarSvc)

	createCarSvc := service.NewCreate(repo, d)
	createHandler := car_handler.NewCreate(createCarSvc)

	updateCarSvc := service.NewUpdate(repo)
	updateHandler := car_handler.NewUpdate(updateCarSvc)

	deleteCarSvc := service.NewDelete(repo)
	deleteHandler := car_handler.NewDelete(deleteCarSvc)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/cars", createHandler.Handle)
		r.Get("/cars", listHandler.Handle)
		r.Get("/cars/{id}", getCarHandler.Handle)
		r.Put("/cars/{id}", updateHandler.Handle)
		r.Delete("/cars/{id}", deleteHandler.Handle)
	})
}

func registerVideoRoutes(r *chi.Mux, db *pgxpool.Pool, d *dispatcher.Dispatcher) {

	repoBlogger := blogger_repo.NewPostgres(db)
	apifyClient := apify.NewClient("apify_api_hToHCZfuFv6sMfR0Z20rVpQ2W1b5Fx3uKK7W")
	videoSearcher := apify.NewVideoSearcher(apifyClient)
	fetchVideoCmd := blogger_cmd.NewFetchBloggerVideos(repoBlogger, videoSearcher)
	fetchVideoHandler := blogger_handler.NewGetCarHandler(fetchVideoCmd)

	r.Route("/test", func(r chi.Router) {
		r.Get("/videos-apify/{id}", fetchVideoHandler.Handle)
	})
}
