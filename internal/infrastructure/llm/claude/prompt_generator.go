package claude

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/gkarman/demo/internal/infrastructure/llm"
)

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
	summary, err := llm.ExtractSummary(rawPayload)
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
			{Text: llm.SystemPrompt},
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
