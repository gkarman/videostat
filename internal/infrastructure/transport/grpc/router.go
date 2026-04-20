package grpc

import (
	"log/slog"

	api_v1 "github.com/gkarman/demo/api/gen/go/v1"
	"github.com/gkarman/demo/internal/application/blogger/query"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	"github.com/gkarman/demo/internal/infrastructure/repository/blogger"
	"github.com/gkarman/demo/internal/infrastructure/transport/grpc/handler/api"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterServices(server *Server, log *slog.Logger, db *pgxpool.Pool, d *dispatcher.Dispatcher) {
	registerBloggerService(server, log, db, d)
}

func registerBloggerService(server *Server, log *slog.Logger, db *pgxpool.Pool, _ *dispatcher.Dispatcher) {
	queryRepo := blogger.NewQueryPostgres(db)
	listBloggersQuery := query.NewListBloggers(queryRepo)

	handler := api.NewHandler(log, listBloggersQuery)
	api_v1.RegisterAPIServer(server.Registrar(), handler)
}
