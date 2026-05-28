package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gkarman/demo/internal/domain/blogger"
	"github.com/google/uuid"
)

const brollSystemPrompt = `You are a video producer. Given a transcript with word-level millisecond timestamps, split the words into semantic segments of approximately 5-8 seconds each. For each segment produce:
- start_ms: start timestamp (ms) of the first word
- end_ms: end timestamp (ms) of the last word
- text: the words of this segment joined with spaces
- broll_prompt: a short English prompt (max 15 words) for generating a relevant background B-roll video clip

Return ONLY a valid JSON object in this exact format: {"segments": [...]}`

type BrollGenerator struct {
	apiKey string
	model  string
	http   *http.Client
}

func NewBrollGenerator(apiKey, model string) *BrollGenerator {
	return &BrollGenerator{
		apiKey: apiKey,
		model:  model,
		http:   &http.Client{},
	}
}

func (g *BrollGenerator) GenerateSegments(ctx context.Context, rawPayload []byte) ([]*blogger.BrollSegment, error) {
	words, err := extractWords(rawPayload)
	if err != nil {
		return nil, fmt.Errorf("extract words: %w", err)
	}

	wordsJSON, err := json.Marshal(words)
	if err != nil {
		return nil, fmt.Errorf("marshal words: %w", err)
	}

	body, err := json.Marshal(map[string]any{
		"model": g.model,
		"messages": []map[string]string{
			{"role": "system", "content": brollSystemPrompt},
			{"role": "user", "content": string(wordsJSON)},
		},
		"max_tokens":      2048,
		"response_format": map[string]string{"type": "json_object"},
	})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+g.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("openai api call: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openai api error %d: %s", resp.StatusCode, respBody)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("openai returned empty response")
	}

	var parsed struct {
		Segments []struct {
			StartMS     int    `json:"start_ms"`
			EndMS       int    `json:"end_ms"`
			Text        string `json:"text"`
			BrollPrompt string `json:"broll_prompt"`
		} `json:"segments"`
	}
	if err := json.Unmarshal([]byte(result.Choices[0].Message.Content), &parsed); err != nil {
		return nil, fmt.Errorf("parse segments json: %w", err)
	}
	raw := parsed.Segments

	segments := make([]*blogger.BrollSegment, 0, len(raw))
	for i, s := range raw {
		segments = append(segments, &blogger.BrollSegment{
			ID:          uuid.NewString(),
			Position:    i,
			StartMS:     s.StartMS,
			EndMS:       s.EndMS,
			Text:        s.Text,
			BrollPrompt: s.BrollPrompt,
			CreatedAt:   time.Now(),
		})
	}

	return segments, nil
}

func extractWords(rawPayload []byte) (any, error) {
	var payload struct {
		Words []struct {
			Text  string `json:"text"`
			Start int    `json:"start"`
			End   int    `json:"end"`
		} `json:"words"`
	}
	if err := json.Unmarshal(rawPayload, &payload); err != nil {
		return nil, err
	}
	return payload.Words, nil
}
