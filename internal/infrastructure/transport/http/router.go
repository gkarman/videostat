package http

import (
	"log/slog"

	"github.com/gkarman/demo/internal/application/car/service"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	"github.com/gkarman/demo/internal/infrastructure/repository/car"
	"github.com/gkarman/demo/internal/infrastructure/transport/http/handler"
	car2 "github.com/gkarman/demo/internal/infrastructure/transport/http/handler/car"
	middleware2 "github.com/gkarman/demo/internal/infrastructure/transport/http/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(log *slog.Logger, db *pgxpool.Pool, d *dispatcher.Dispatcher) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware2.Logger(log))
	r.Use(middleware2.Recovery())
	registerHomeRoutes(r)
	registerCarRoutes(r, db, d)
	return r
}

func registerHomeRoutes(r *chi.Mux) {
	homeHandler := handler.NewHomeHandler()
	r.Get("/", homeHandler.Home)
}

func registerCarRoutes(r *chi.Mux, db *pgxpool.Pool, d *dispatcher.Dispatcher) {
	repo := car.NewPostgresRepo(db)

	listSvc := service.NewList(repo)
	listHandler := car2.NewList(listSvc)

	getCarSvc := service.NewGet(repo)
	getCarHandler := car2.NewGetCarHandler(getCarSvc)

	createCarSvc := service.NewCreate(repo, d)
	createHandler := car2.NewCreate(createCarSvc)

	updateCarSvc := service.NewUpdate(repo)
	updateHandler := car2.NewUpdate(updateCarSvc)

	deleteCarSvc := service.NewDelete(repo)
	deleteHandler := car2.NewDelete(deleteCarSvc)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/cars", createHandler.Handle)
		r.Get("/cars", listHandler.Handle)
		r.Get("/cars/{id}", getCarHandler.Handle)
		r.Put("/cars/{id}", updateHandler.Handle)
		r.Delete("/cars/{id}", deleteHandler.Handle)
	})
}
