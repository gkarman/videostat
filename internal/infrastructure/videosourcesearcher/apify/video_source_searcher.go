package apify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gkarman/demo/internal/domain/blogger"
	"github.com/gkarman/demo/internal/infrastructure/logger"
	apifyclient "github.com/gkarman/demo/internal/infrastructure/videosearcher/apify"
)

type VideoSourceSearcher struct {
	client *apifyclient.Client
}

func NewVideoSourceSearcher(client *apifyclient.Client) *VideoSourceSearcher {
	return &VideoSourceSearcher{client: client}
}

type youtubeDownloadResult struct {
	DownloadedFileURL string `json:"downloadedFileUrl"`
}

func (s *VideoSourceSearcher) SearchUrl(ctx context.Context, v *blogger.Video) (string, error) {
	log := logger.FromContext(ctx).With("component", "VideoSourceSearcher", "videoID", v.ID)

	switch {
	case isYouTube(v.URL):
		log.Debug("searching youtube source", "url", v.URL)
		return s.searchYouTube(ctx, v.URL)
	case isTikTok(v.URL):
		log.Debug("tiktok source search not implemented")
		return "", errors.New("tiktok source search not implemented")
	case isInstagram(v.URL):
		log.Debug("instagram source search not implemented")
		return "", errors.New("instagram source search not implemented")
	default:
		return "", fmt.Errorf("unsupported video url: %s", v.URL)
	}
}

func (s *VideoSourceSearcher) searchYouTube(ctx context.Context, videoURL string) (string, error) {
	input := map[string]any{
		"storeInKVStore": true,
		"videos": []map[string]string{
			{"url": videoURL},
		},
	}

	raw, err := s.client.RunActorSync(ctx, "streamers~youtube-video-downloader", input)
	if err != nil {
		return "", fmt.Errorf("run actor: %w", err)
	}

	var items []youtubeDownloadResult
	if err := json.Unmarshal(raw, &items); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if len(items) == 0 {
		return "", errors.New("no download url found")
	}

	if items[0].DownloadedFileURL == "" {
		return "", errors.New("download url is empty")
	}

	return items[0].DownloadedFileURL, nil
}

func isYouTube(url string) bool {
	return strings.Contains(url, "youtube.com") || strings.Contains(url, "youtu.be")
}

func isTikTok(url string) bool {
	return strings.Contains(url, "tiktok.com")
}

func isInstagram(url string) bool {
	return strings.Contains(url, "instagram.com")
}