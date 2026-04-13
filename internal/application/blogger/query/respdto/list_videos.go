package respdto

import "github.com/gkarman/demo/internal/application/blogger/query/view"

type ListVideos struct {
	Items []*view.Video
}
