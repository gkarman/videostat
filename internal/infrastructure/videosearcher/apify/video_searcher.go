package apify

import (
	"context"
	"encoding/json"
	"fmt"
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
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Likes       int64     `json:"likes"`
	Views       int64     `json:"views"`
	Comments    int64     `json:"comments"`
	PublishedAt time.Time `json:"createdAt"`
}

type instaVideo struct {
	ID          string    `json:"id"`
	URL         string    `json:"link"`
	Title       string    `json:"caption"`
	Likes       int64     `json:"likeCount"`
	Views       int64     `json:"viewCount"`
	Comments    int64     `json:"commentsCount"`
	PublishedAt time.Time `json:"takenAt"`
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
	log := logger.FromContext(ctx)

	log.Debug("находим данные для блогера", "id", b.ID, "url", b.URL)
	input := map[string]any{
		"channels": []string{b.URL},
		"maxResultsShorts": 10,
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

	days := 40
	return s.youtubeToVideos(b, items, days), nil
}

func (s *VideoSearcher) searchTikTok(ctx context.Context, b *blogger.Blogger) ([]*blogger.Video, error) {
	input := map[string]any{
		"startUrls":  []map[string]string{{"url": b.URL}},
		"maxResults": 10,
	}
	raw, err := s.client.RunActorSync(ctx, "clockworks~tiktok-scraper", input)
	if err != nil {
		return nil, err
	}

	var items []tiktokVideo
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, err
	}

	days := 10
	return s.tiktokToVideos(b, items, days), nil
}

func (s *VideoSearcher) searchInstagram(ctx context.Context, b *blogger.Blogger) ([]*blogger.Video, error) {
	input := map[string]any{
		"startUrls":  []map[string]string{{"url": b.URL}},
		"maxResults": 10,
	}
	raw, err := s.client.RunActorSync(ctx, "apify~instagram-scraper", input)
	if err != nil {
		return nil, err
	}

	var items []instaVideo
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, err
	}

	days := 10
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