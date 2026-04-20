package api

import (
	"context"

	api_v1 "github.com/gkarman/demo/api/gen/go/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GetBloggerList(ctx context.Context, _ *api_v1.GetBloggerListRequest) (*api_v1.GetBloggerListResponse, error) {
	resp, err := h.listBloggers.Run(ctx)
	if err != nil {
		h.log.Error("get blogger list failed", "error", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	items := make([]*api_v1.BloggerInfo, 0, len(resp.Items))
	for _, b := range resp.Items {
		items = append(items, &api_v1.BloggerInfo{
			Id:       b.ID,
			Planform: b.Platform,
			Url:      b.URL,
		})
	}

	return &api_v1.GetBloggerListResponse{
		Bloggers: items,
	}, nil
}
