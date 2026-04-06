package blogger

import "context"

type Repo interface {
	Save(ctx context.Context, blogger *Blogger) error
	ExistByUrl(ctx context.Context, url string) (bool, error)
}