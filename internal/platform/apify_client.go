package platform

import (
	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/infrastructure/videosearcher/apify"
)

func NewApifyClient(cfg *config.Config) *apify.Client {
	confApify := apify.Config{
		Token: cfg.Apify.Token,
		Host:  cfg.Apify.Host,
		Limits: apify.Limits{
			YouTubeLimits: apify.YouTubeLimits{
				MaxVideos: cfg.Apify.YoutubeMaxVideos,
				Days:      cfg.Apify.YoutubeDays,
			},
			TikTokLimits: apify.TikTokLimits{
				MaxVideos: cfg.Apify.TiktokMaxVideos,
				Days:      cfg.Apify.TiktokDays,
			},
			InstagramLimits: apify.InstagramLimits{
				MaxVideos: cfg.Apify.InstagramMaxVideos,
				Days:      cfg.Apify.InstagramDays,
			},
		},
	}

	return apify.NewClient(confApify)
}