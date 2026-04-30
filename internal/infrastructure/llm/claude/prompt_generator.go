package claude

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

const systemPrompt = `You are an expert content strategist and video scriptwriter specializing in talking head videos.

Given a video transcript with analysis data (highlights, chapters, sentiment, entities, summary),
your task is to generate a detailed creative prompt for producing a new talking head video
that captures the same core message, energy, and persuasive hooks.

The prompt must include:
1. Core message — the single most important idea the video conveys
2. Target audience — who this content is for
3. Opening hook — the first 10 seconds that will stop the scroll
4. Key talking points — 3-5 main arguments or insights to cover
5. Emotional arc — how the energy and tone should shift throughout
6. Memorable phrases — specific lines or formulations worth keeping or adapting
7. Delivery style — pace, tone, body language cues for the speaker
8. Closing CTA — how to end and what action to drive

Return only the prompt text, ready to use for video production. No preamble, no explanations.`

type analysisSummary struct {
	Transcript     string          `json:"transcript"`
	Summary        any             `json:"summary,omitempty"`
	Chapters       any             `json:"chapters,omitempty"`
	KeyHighlights  any             `json:"key_highlights,omitempty"`
	SentimentPeaks any             `json:"sentiment_peaks,omitempty"`
	Entities       any             `json:"entities,omitempty"`
}

type PromptGenerator struct {
	client anthropic.Client
	model  string
}

func NewPromptGenerator(apiKey, model string) *PromptGenerator {
	return &PromptGenerator{
		client: anthropic.NewClient(option.WithAPIKey(apiKey)),
		model:  model,
	}
}

func (g *PromptGenerator) ProviderName() string { return "anthropic" }
func (g *PromptGenerator) ModelName() string    { return g.model }

func (g *PromptGenerator) Generate(ctx context.Context, rawPayload []byte) (string, error) {
	summary, err := extractSummary(rawPayload)
	if err != nil {
		return "", fmt.Errorf("extract analysis summary: %w", err)
	}

	summaryJSON, err := json.Marshal(summary)
	if err != nil {
		return "", fmt.Errorf("marshal summary: %w", err)
	}

	msg, err := g.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(g.model),
		MaxTokens: 2048,
		System: []anthropic.TextBlockParam{
			{Text: systemPrompt},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewTextBlock(string(summaryJSON)),
			),
		},
	})
	if err != nil {
		return "", fmt.Errorf("claude api call: %w", err)
	}

	if len(msg.Content) == 0 {
		return "", fmt.Errorf("claude returned empty response")
	}

	block := msg.Content[0]
	if block.Type != "text" {
		return "", fmt.Errorf("unexpected claude response content type: %s", block.Type)
	}

	return block.Text, nil
}

func extractSummary(raw []byte) (*analysisSummary, error) {
	var payload struct {
		Text    string `json:"text"`
		Summary string `json:"summary"`
		Chapters []struct {
			Start    int    `json:"start"`
			End      int    `json:"end"`
			Headline string `json:"headline"`
			Gist     string `json:"gist"`
			Summary  string `json:"summary"`
		} `json:"chapters"`
		AutoHighlightsResult struct {
			Results []struct {
				Rank       float64 `json:"rank"`
				Text       string  `json:"text"`
				Count      int     `json:"count"`
			} `json:"results"`
		} `json:"auto_highlights_result"`
		SentimentAnalysisResults []struct {
			Text      string  `json:"text"`
			Sentiment string  `json:"sentiment"`
			Confidence float64 `json:"confidence"`
			Start     int     `json:"start"`
			End       int     `json:"end"`
		} `json:"sentiment_analysis_results"`
		Entities []struct {
			Text        string `json:"text"`
			EntityType  string `json:"entity_type"`
		} `json:"entities"`
	}

	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}

	s := &analysisSummary{
		Transcript: payload.Text,
	}

	if payload.Summary != "" {
		s.Summary = payload.Summary
	}
	if len(payload.Chapters) > 0 {
		s.Chapters = payload.Chapters
	}
	if len(payload.AutoHighlightsResult.Results) > 0 {
		s.KeyHighlights = payload.AutoHighlightsResult.Results
	}
	if len(payload.Entities) > 0 {
		s.Entities = payload.Entities
	}

	// Only keep high-confidence sentiment peaks (POSITIVE/NEGATIVE with conf > 0.8)
	var peaks []any
	for _, r := range payload.SentimentAnalysisResults {
		if r.Sentiment != "NEUTRAL" && r.Confidence > 0.8 {
			peaks = append(peaks, r)
		}
	}
	if len(peaks) > 0 {
		s.SentimentPeaks = peaks
	}

	return s, nil
}
