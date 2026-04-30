package videoanalyzer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL        = "https://api.assemblyai.com/v2"
	pollInterval   = 5 * time.Second
	analyzeTimeout = 10 * time.Minute
)

type AssemblyAIAnalyzer struct {
	apiKey string
	client *http.Client
}

func NewAssemblyAIAnalyzer(apiKey string) *AssemblyAIAnalyzer {
	return &AssemblyAIAnalyzer{
		apiKey: apiKey,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (a *AssemblyAIAnalyzer) ProviderName() string {
	return "assemblyai"
}

func (a *AssemblyAIAnalyzer) Analyze(ctx context.Context, videoURL string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, analyzeTimeout)
	defer cancel()

	id, err := a.submit(ctx, videoURL)
	if err != nil {
		return nil, err
	}

	return a.wait(ctx, id)
}

func (a *AssemblyAIAnalyzer) submit(ctx context.Context, audioURL string) (string, error) {
	payload, err := json.Marshal(map[string]any{
		"audio_url":          audioURL,
		"speech_models":      []string{"universal-2"},
		"auto_highlights":    true,
		"sentiment_analysis": true,
		"entity_detection":   true,
		"summarization":      true,
		"summary_model":      "informative",
		"summary_type":       "bullets",
	})
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/transcript", bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", a.apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("submit transcript: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if res.StatusCode >= 300 {
		return "", fmt.Errorf("assemblyai submit error %d: %s", res.StatusCode, string(body))
	}

	var resp struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("unmarshal submit response: %w", err)
	}
	if resp.ID == "" {
		return "", fmt.Errorf("assemblyai returned empty transcript id")
	}

	return resp.ID, nil
}

func (a *AssemblyAIAnalyzer) wait(ctx context.Context, id string) ([]byte, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(pollInterval):
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/transcript/"+id, nil)
		if err != nil {
			return nil, fmt.Errorf("create poll request: %w", err)
		}
		req.Header.Set("Authorization", a.apiKey)

		res, err := a.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("poll transcript: %w", err)
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("read poll response: %w", err)
		}

		if res.StatusCode >= 300 {
			return nil, fmt.Errorf("assemblyai poll error %d: %s", res.StatusCode, string(body))
		}

		var resp struct {
			Status string `json:"status"`
			Error  string `json:"error"`
		}
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, fmt.Errorf("unmarshal poll response: %w", err)
		}

		switch resp.Status {
		case "completed":
			return body, nil
		case "error":
			return nil, fmt.Errorf("assemblyai transcription failed: %s", resp.Error)
		}
	}
}
