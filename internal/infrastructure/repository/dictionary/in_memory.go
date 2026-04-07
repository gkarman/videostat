package dictionary

import (
	"context"

	"github.com/gkarman/demo/internal/domain/dictionary"
)

const PLATFORM_NAME_YOUTUBE = "youtube"
const PLATFORM_NAME_TIKTOK = "tiktok"

type FakeRepo struct {
	platforms map[string]*dictionary.Platform
}

func NewFake() *FakeRepo {
	platforms := make(map[string]*dictionary.Platform)
	platforms[PLATFORM_NAME_YOUTUBE] = &dictionary.Platform{ID: 1, Name: PLATFORM_NAME_YOUTUBE}
	platforms[PLATFORM_NAME_TIKTOK] = &dictionary.Platform{ID: 2, Name: PLATFORM_NAME_TIKTOK}

	return &FakeRepo{
		platforms: platforms,
	}
}

func (r *FakeRepo) GetPlatformByName(_ context.Context, name string) (*dictionary.Platform, error) {
	if p, ok := r.platforms[name]; ok {
		return p, nil
	}
	return nil, nil
}
