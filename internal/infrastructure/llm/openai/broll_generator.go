package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gkarman/demo/internal/domain/blogger"
	"github.com/google/uuid"
)

const brollSystemPrompt = `You are a viral TikTok/Reels video editor. You will receive pre-split transcript segments with exact timestamps. Your ONLY job is to write a broll_prompt for each segment.

For each segment write an energetic B-roll prompt (40-60 words):

ENERGY & PACE (CRITICAL):
- ABSOLUTELY NO slow motion, NO slow-mo, NO cinematic slow, NO time-lapse
- ABSOLUTELY NO static shots, NO smooth dolly, NO peaceful or calm scenes
- ABSOLUTELY NO CGI, NO animation, NO 3D render, NO video game look — REAL people, REAL locations only
- Every shot MUST show REAL humans in fast action: running, grabbing, reacting, arguing, working fast
- Camera ALWAYS moving: "whip pan", "rapid handheld shake", "crash zoom", "urgent tracking"
- Music video energy — every second must feel alive and urgent

PEOPLE & SETTING (strictly required):
- ALL people must be American/Western: white, Black, Hispanic — NEVER Asian faces
- REAL people doing REAL things — not poses, not stock-photo stillness
- Contemporary USA 2020s: American diners, US city streets, American offices, US homes
- Modern iPhones, US cars, American brands visible
- NO Asian people, NO Asian locations, NO China/Japan/Korea
- All text/signs in English only

TECHNICAL:
- Lighting: specific ("harsh fluorescent overhead", "cold blue neon glow", "warm golden backlight")
- Style: "hyper-realistic 4K handheld", "raw documentary urgency"
- NO subtitles, NO captions, NO watermarks
- All segments form one cohesive visual story

Return ONLY valid JSON: {"segments": [{"start_ms": int, "end_ms": int, "text": "...", "broll_prompt": "..."}]}`

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
	transcript, words, err := extractWords(rawPayload)
	if err != nil {
		return nil, fmt.Errorf("extract words: %w", err)
	}

	segments := splitIntoSegments(words, 3000)
	if len(segments) == 0 {
		return nil, fmt.Errorf("no segments after split")
	}

	segJSON, err := json.Marshal(segments)
	if err != nil {
		return nil, fmt.Errorf("marshal segments: %w", err)
	}

	userMsg := fmt.Sprintf("FULL TRANSCRIPT:\n%s\n\nSEGMENTS (write broll_prompt for each):\n%s", transcript, string(segJSON))

	body, err := json.Marshal(map[string]any{
		"model": g.model,
		"messages": []map[string]string{
			{"role": "system", "content": brollSystemPrompt},
			{"role": "user", "content": userMsg},
		},
		"max_tokens":      4096,
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

	result2 := make([]*blogger.BrollSegment, 0, len(parsed.Segments))
	for i, s := range parsed.Segments {
		result2 = append(result2, &blogger.BrollSegment{
			ID:          uuid.NewString(),
			Position:    i,
			StartMS:     s.StartMS,
			EndMS:       s.EndMS,
			Text:        s.Text,
			BrollPrompt: s.BrollPrompt,
			CreatedAt:   time.Now(),
		})
	}

	return result2, nil
}

type word struct {
	Text  string `json:"text"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

type segment struct {
	StartMS int    `json:"start_ms"`
	EndMS   int    `json:"end_ms"`
	Text    string `json:"text"`
}

// splitIntoSegments splits words into segments of ~targetMS milliseconds each.
func splitIntoSegments(words []word, targetMS int) []segment {
	if len(words) == 0 {
		return nil
	}

	var segments []segment
	segStart := words[0].Start
	var texts []string

	for i, w := range words {
		texts = append(texts, w.Text)
		duration := w.End - segStart

		isLast := i == len(words)-1
		if duration >= targetMS || isLast {
			segments = append(segments, segment{
				StartMS: segStart,
				EndMS:   w.End,
				Text:    strings.Join(texts, " "),
			})
			if !isLast {
				segStart = words[i+1].Start
				texts = nil
			}
		}
	}

	return segments
}

func extractWords(rawPayload []byte) (string, []word, error) {
	var payload struct {
		Text  string `json:"text"`
		Words []word `json:"words"`
	}
	if err := json.Unmarshal(rawPayload, &payload); err != nil {
		return "", nil, err
	}
	return payload.Text, payload.Words, nil
}
