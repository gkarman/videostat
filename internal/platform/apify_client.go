package platform

import (

	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/infrastructure/videosearcher/apify"
)

func NewApifyClient(cfg *config.Config) *apify.Client {
	confApify := apify.Config{
		Token: cfg.Apify.Token,
		Host:  cfg.Apify.Host,
	}

	return apify.NewClient(confApify)
}