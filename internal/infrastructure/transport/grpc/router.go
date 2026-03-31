package grpc

import (
	"log/slog"

	carv1 "github.com/gkarman/demo/api/gen/go/car/v1"
	"github.com/gkarman/demo/internal/application/car/service"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	carrepository "github.com/gkarman/demo/internal/infrastructure/repository/car"
	carhandler "github.com/gkarman/demo/internal/infrastructure/transport/grpc/handler/car"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterServices(server *Server, log *slog.Logger, db *pgxpool.Pool, d *dispatcher.Dispatcher) {
	registerCarService(server, log, db, d)
}

func registerCarService(server *Server, log *slog.Logger, db *pgxpool.Pool, d *dispatcher.Dispatcher) {
	repo := carrepository.NewPostgresRepo(db)
	getSvc := service.NewGet(repo)
	listSvc := service.NewList(repo)
	createSvc := service.NewCreate(repo, d)
	updateSvc := service.NewUpdate(repo)
	deleteSvc := service.NewDelete(repo)

	handler := carhandler.New(log, getSvc, listSvc, createSvc, updateSvc, deleteSvc)
	carv1.RegisterCarServer(server.Registrar(), handler)
}
