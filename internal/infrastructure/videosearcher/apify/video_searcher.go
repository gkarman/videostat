package apify

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gkarman/demo/internal/domain/blogger"
	"github.com/gkarman/demo/internal/infrastructure/logger"
	"github.com/google/uuid"
)

type VideoSearcher struct {
	client *Client
}

func NewVideoSearcher(client *Client) *VideoSearcher {
	return &VideoSearcher{client: client}
}

type youtubeVideo struct {
	VideoID     string    `json:"id"`
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Views       int64     `json:"viewCount"`
	Likes       int64     `json:"likes"`
	Comments    int64     `json:"commentsCount"`
	PublishedAt time.Time `json:"date"`
}

type tiktokVideo struct {
	ID          string    `json:"id"`
	URL         string    `json:"webVideoUrl"`
	Title       string    `json:"text"`
	Likes       int64     `json:"diggCount"`
	Views       int64     `json:"playCount"`
	Comments    int64     `json:"commentCount"`
	PublishedAt time.Time `json:"createTimeISO"`
}

type instaVideo struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	Title       string    `json:"caption"`
	Likes       int64     `json:"likesCount"`
	Views       int64     `json:"videoViewCount"`
	Comments    int64     `json:"commentsCount"`
	PublishedAt time.Time `json:"timestamp"`
}

func (s *VideoSearcher) Search(ctx context.Context, b *blogger.Blogger) ([]*blogger.Video, error) {
	switch b.PlatformID {
	case 1:
		return s.searchYouTube(ctx, b)
	case 2:
		return s.searchTikTok(ctx, b)
	case 3:
		return s.searchInstagram(ctx, b)
	default:
		return nil, fmt.Errorf("unsupported platform: %d", b.PlatformID)
	}
}

func (s *VideoSearcher) searchYouTube(ctx context.Context, b *blogger.Blogger) ([]*blogger.Video, error) {
	var (
		maxResultsShorts = s.client.cfg.Limits.YouTubeLimits.MaxVideos
		days             = s.client.cfg.Limits.YouTubeLimits.Days
	)

	log := logger.FromContext(ctx)

	log.Debug("находим данные для блогера", "id", b.ID, "url", b.URL)
	input := map[string]any{
		"channels":         []string{b.URL},
		"maxResultsShorts": maxResultsShorts,
	}

	raw, err := s.client.RunActorSync(ctx, "streamers~youtube-shorts-scraper", input)
	log.Debug("Получили ответ", raw)
	if err != nil {
		return nil, err
	}

	var items []youtubeVideo
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, err
	}
	log.Debug("проверка", "кол-во элементов", len(items))
	log.Debug("проверка", "элементы", items)

	return s.youtubeToVideos(b, items, days), nil
}

func (s *VideoSearcher) searchTikTok(ctx context.Context, b *blogger.Blogger) ([]*blogger.Video, error) {
	var (
		maxItems       = s.client.cfg.Limits.TikTokLimits.MaxVideos
		resultsPerPage = s.client.cfg.Limits.TikTokLimits.MaxVideos
		days           = s.client.cfg.Limits.TikTokLimits.Days
	)

	log := logger.FromContext(ctx)
	log.Debug("находим данные для блогера (tiktok)", "id", b.ID, "url", b.URL)

	input := map[string]any{
		"profiles":              []string{extractTikTokUsername(b.URL)},
		"maxItems":              maxItems,
		"resultsPerPage":        resultsPerPage,
		"profileScrapeSections": []string{"videos"},
	}
	raw, err := s.client.RunActorSync(ctx, "clockworks~tiktok-scraper", input)
	if err != nil {
		return nil, err
	}
	log.Debug("tiktok", "raw", raw)

	var items []tiktokVideo
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, err
	}

	log.Debug("tiktok", "items", items)
	return s.tiktokToVideos(b, items, days), nil
}

func (s *VideoSearcher) searchInstagram(ctx context.Context, b *blogger.Blogger) ([]*blogger.Video, error) {
	var (
		maxResults = s.client.cfg.Limits.InstagramLimits.MaxVideos
		days       = s.client.cfg.Limits.InstagramLimits.Days
	)

	log := logger.FromContext(ctx)
	log.Debug("находим данные для блогера (Instagram)", "id", b.ID, "url", b.URL)

	input := map[string]any{
		"directUrls":    []string{b.URL},
		"addParentData": false,
		"resultsLimit":  maxResults,
		"searchLimit":   maxResults,
		"resultsType":   "reels",
		"searchType":    "hashtag",
	}
	raw, err := s.client.RunActorSync(ctx, "apify~instagram-scraper", input)
	if err != nil {
		return nil, err
	}
	log.Debug("instagram", "raw", raw)

	var items []instaVideo
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, err
	}
	log.Debug("instagram", "items", items)

	return s.instaToVideos(b, items, days), nil
}

func (s *VideoSearcher) youtubeToVideos(b *blogger.Blogger, items []youtubeVideo, days int) []*blogger.Video {
	now := time.Now()
	minDate := now.AddDate(0, 0, -days)
	var result []*blogger.Video
	for _, it := range items {
		if it.PublishedAt.Before(minDate) {
			continue
		}
		result = append(result, &blogger.Video{
			ID:          uuid.NewString(),
			BloggerID:   b.ID,
			ExternalID:  it.VideoID,
			URL:         it.URL,
			Title:       it.Title,
			Views:       it.Views,
			Likes:       it.Likes,
			Comments:    it.Comments,
			PublishedAt: it.PublishedAt,
			CreatedAt:   now,
		})
	}
	return result
}

func (s *VideoSearcher) tiktokToVideos(b *blogger.Blogger, items []tiktokVideo, days int) []*blogger.Video {
	var result []*blogger.Video

	for _, it := range items {
		result = append(result, &blogger.Video{
			ID:          uuid.NewString(),
			BloggerID:   b.ID,
			ExternalID:  it.ID,
			URL:         it.URL,
			Title:       it.Title,
			Views:       it.Views,
			Likes:       it.Likes,
			Comments:    it.Comments,
			PublishedAt: it.PublishedAt,
			CreatedAt:   time.Now(),
		})
	}

	return result
}

func (s *VideoSearcher) instaToVideos(b *blogger.Blogger, items []instaVideo, days int) []*blogger.Video {
	now := time.Now()
	minDate := now.AddDate(0, 0, -days)
	var result []*blogger.Video
	for _, it := range items {
		if it.PublishedAt.Before(minDate) {
			continue
		}
		result = append(result, &blogger.Video{
			ID:          uuid.NewString(),
			BloggerID:   b.ID,
			ExternalID:  it.ID,
			URL:         it.URL,
			Title:       it.Title,
			Views:       it.Views,
			Likes:       it.Likes,
			Comments:    it.Comments,
			PublishedAt: it.PublishedAt,
			CreatedAt:   now,
		})
	}
	return result
}

func extractTikTokUsername(url string) string {
	// https://www.tiktok.com/@humphreytalks -> humphreytalks
	parts := strings.Split(url, "@")
	if len(parts) < 2 {
		return ""
	}
	return strings.Split(parts[1], "?")[0]
}
