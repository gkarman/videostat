package api

import (
	"log/slog"

	api_v1 "github.com/gkarman/demo/api/gen/go/v1"
	"github.com/gkarman/demo/internal/application/blogger/query"
)

type Handler struct {
	api_v1.UnimplementedAPIServer
	log          *slog.Logger
	listBloggers *query.ListBloggers
}

func NewHandler(
	log *slog.Logger,
	listBloggers *query.ListBloggers,
) *Handler {
	return &Handler{
		log:          log,
		listBloggers: listBloggers,
	}
}
