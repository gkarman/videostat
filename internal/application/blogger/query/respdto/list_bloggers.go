package respdto

import "github.com/gkarman/demo/internal/application/blogger/query/view"

type ListBloggers struct {
	Items []*view.Blogger
}
