package platform

import (
	"fmt"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/infrastructure/llm/claude"
	"github.com/gkarman/demo/internal/infrastructure/llm/openai"
	"github.com/gkarman/demo/internal/infrastructure/llm/openrouter"
)

func NewVideoPromptGenerator(cfg *config.Config) (application.VideoPromptGenerator, error) {
	switch cfg.LLM.Provider {
	case "anthropic":
		return claude.NewPromptGenerator(cfg.Anthropic.Token, cfg.Anthropic.Model), nil
	case "openrouter":
		return openrouter.NewPromptGenerator(cfg.OpenRouter.Token, cfg.OpenRouter.Model), nil
	case "openai":
		return openai.NewPromptGenerator(cfg.OpenAI.Token, cfg.OpenAI.Model), nil
	default:
		return nil, fmt.Errorf("unknown LLM provider %q: supported values are anthropic, openrouter, openai", cfg.LLM.Provider)
	}
}
