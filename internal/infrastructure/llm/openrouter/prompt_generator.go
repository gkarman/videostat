package openrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gkarman/demo/internal/infrastructure/llm"
)

const baseURL = "https://openrouter.ai/api/v1/chat/completions"

type PromptGenerator struct {
	apiKey string
	model  string
	http   *http.Client
}

func NewPromptGenerator(apiKey, model string) *PromptGenerator {
	return &PromptGenerator{
		apiKey: apiKey,
		model:  model,
		http:   &http.Client{},
	}
}

func (g *PromptGenerator) ProviderName() string { return "openrouter" }
func (g *PromptGenerator) ModelName() string    { return g.model }

func (g *PromptGenerator) Generate(ctx context.Context, rawPayload []byte) (string, error) {
	summary, err := llm.ExtractSummary(rawPayload)
	if err != nil {
		return "", fmt.Errorf("extract analysis summary: %w", err)
	}

	summaryJSON, err := json.Marshal(summary)
	if err != nil {
		return "", fmt.Errorf("marshal summary: %w", err)
	}

	body, err := json.Marshal(map[string]any{
		"model": g.model,
		"messages": []map[string]string{
			{"role": "system", "content": llm.SystemPrompt},
			{"role": "user", "content": string(summaryJSON)},
		},
		"max_tokens": 2048,
	})
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	const maxAttempts = 4
	delay := 5 * time.Second

	for attempt := range maxAttempts {
		result, err := g.do(ctx, body)
		if err == nil {
			return result, nil
		}
		if attempt == maxAttempts-1 {
			return "", err
		}
		var apiErr *apiError
		if isRateLimit(err, &apiErr) {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(delay):
				delay *= 2
				continue
			}
		}
		return "", err
	}
	panic("unreachable")
}

type apiError struct {
	status int
	body   []byte
}

func (e *apiError) Error() string {
	return fmt.Sprintf("openrouter api error %d: %s", e.status, e.body)
}

func isRateLimit(err error, out **apiError) bool {
	if e, ok := err.(*apiError); ok && e.status == http.StatusTooManyRequests {
		*out = e
		return true
	}
	return false
}

func (g *PromptGenerator) do(ctx context.Context, body []byte) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+g.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("openrouter api call: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", &apiError{status: resp.StatusCode, body: respBody}
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	if len(result.Choices) == 0 || result.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("openrouter returned empty response")
	}

	return result.Choices[0].Message.Content, nil
}
