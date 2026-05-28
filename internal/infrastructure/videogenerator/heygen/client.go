package heygen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gkarman/demo/internal/application"
)

const (
	baseURL         = "https://api.heygen.com"
	videoWidth      = 1920
	videoHeight     = 1080
	httpTimeout     = 30 * time.Second
)

type Client struct {
	apiKey   string
	avatarID string
	voiceID  string
	http     *http.Client
}

func NewClient(apiKey, avatarID, voiceID string) *Client {
	return &Client{
		apiKey:   apiKey,
		avatarID: avatarID,
		voiceID:  voiceID,
		http:     &http.Client{Timeout: httpTimeout},
	}
}

func (c *Client) Platform() string { return "heygen" }

func (c *Client) Submit(ctx context.Context, script string) (string, error) {
	body, err := json.Marshal(map[string]any{
		"video_inputs": []map[string]any{
			{
				"character": map[string]any{
					"type":      "avatar",
					"avatar_id": c.avatarID,
					"scale":     1,
				},
				"voice": map[string]any{
					"type":       "text",
					"voice_id":   c.voiceID,
					"input_text": script,
				},
				"background": map[string]any{
					"type":  "color",
					"value": "#000000",
				},
			},
		},
		"dimension": map[string]int{
			"width":  videoWidth,
			"height": videoHeight,
		},
	})
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/v2/video/generate", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("heygen api call: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("heygen submit error %d: %s", resp.StatusCode, respBody)
	}

	var result struct {
		Error *string `json:"error"`
		Data  struct {
			VideoID string `json:"video_id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("parse submit response: %w", err)
	}
	if result.Error != nil {
		return "", fmt.Errorf("heygen error: %s", *result.Error)
	}
	if result.Data.VideoID == "" {
		return "", fmt.Errorf("heygen returned empty video_id")
	}

	return result.Data.VideoID, nil
}

func (c *Client) GetStatus(ctx context.Context, externalID string) (*application.GenerationStatus, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/v1/video_status.get", nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("x-api-key", c.apiKey)
	q := req.URL.Query()
	q.Set("video_id", externalID)
	req.URL.RawQuery = q.Encode()

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("heygen api call: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("heygen status error %d: %s", resp.StatusCode, respBody)
	}

	var result struct {
		Data struct {
			Status   string `json:"status"`
			VideoURL string `json:"video_url"`
			Error    *struct {
				Code   string `json:"code"`
				Detail string `json:"detail"`
			} `json:"error"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse status response: %w", err)
	}

	s := &application.GenerationStatus{
		Status:   result.Data.Status,
		VideoURL: result.Data.VideoURL,
	}
	if result.Data.Error != nil {
		s.ErrorMessage = result.Data.Error.Detail
	}
	return s, nil
}
