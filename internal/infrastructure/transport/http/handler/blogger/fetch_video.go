package blogger

import (
	"net/http"

	"github.com/gkarman/demo/internal/application/blogger/command"
	"github.com/gkarman/demo/internal/application/blogger/reqdto"
	"github.com/gkarman/demo/internal/infrastructure/logger"
	"github.com/go-chi/chi/v5"
)

type FetchVideo struct {
	c *command.FetchBloggerVideos
}

func NewGetCarHandler(c *command.FetchBloggerVideos) *FetchVideo {
	return &FetchVideo{
		c: c,
	}
}

func (h *FetchVideo) Handle(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())
	log.Info("Начало")
	id := chi.URLParam(r, "id")
	log.Info("id", "id", id)
	req := reqdto.FetchBloggerVideos{BloggerID: id}
	ctx := r.Context()

	_ = h.c.Execute(ctx, req)
}