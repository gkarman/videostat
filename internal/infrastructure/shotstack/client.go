package shotstack

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

type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client
}

func NewClient(apiKey, baseURL string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		http:    &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) Compose(ctx context.Context, req application.CompositionRequest) (string, error) {
	timeline := c.buildTimeline(req)

	body, err := json.Marshal(map[string]any{
		"timeline": timeline,
		"output": map[string]any{
			"format": "mp4",
			"size": map[string]any{
				"width":  720,
				"height": 1280,
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url("/render"), bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("shotstack api call: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("shotstack api error %d: %s", resp.StatusCode, respBody)
	}

	var result struct {
		Success  bool `json:"success"`
		Response struct {
			ID string `json:"id"`
		} `json:"response"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}
	if !result.Success {
		return "", fmt.Errorf("shotstack returned success=false")
	}
	if result.Response.ID == "" {
		return "", fmt.Errorf("shotstack returned empty render ID, body: %s", respBody)
	}

	return result.Response.ID, nil
}

func (c *Client) GetStatus(ctx context.Context, externalID string) (*application.CompositionStatus, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url("/render/"+externalID), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("x-api-key", c.apiKey)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("shotstack api call: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("shotstack api error %d: %s", resp.StatusCode, respBody)
	}

	var result struct {
		Response struct {
			Status string `json:"status"`
			URL    string `json:"url"`
			Error  string `json:"error"`
		} `json:"response"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	r := result.Response
	status := &application.CompositionStatus{}

	switch r.Status {
	case "done":
		status.Status = "done"
		status.ResultURL = r.URL
	case "failed":
		status.Status = "failed"
		status.Error = r.Error
	default:
		status.Status = "processing"
	}

	return status, nil
}

func (c *Client) url(path string) string {
	return c.baseURL + path
}

func (c *Client) buildTimeline(req application.CompositionRequest) map[string]any {
	// Track 0: B-roll clips (background, muted)
	brollClips := make([]any, 0, len(req.BrollClips))
	for _, clip := range req.BrollClips {
		brollClips = append(brollClips, map[string]any{
			"asset":  map[string]any{"type": "video", "src": clip.URL, "volume": 0},
			"start":  clip.StartSec,
			"length": clip.LenSec,
		})
	}

	// Track 1: Avatar video (bottom, contain — no side cropping)
	avatarClip := map[string]any{
		"asset":    map[string]any{"type": "video", "src": req.AvatarURL},
		"start":    0,
		"length":   req.TotalSec,
		"fit":      "contain",
		"position": "bottom",
	}

	return map[string]any{
		"background": "#000000",
		"tracks": []any{
			map[string]any{"clips": []any{avatarClip}},
			map[string]any{"clips": brollClips},
		},
	}
}
