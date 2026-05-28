package platform

import (
	"fmt"

	"github.com/gkarman/demo/internal/application"
	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/infrastructure/kling"
	"github.com/gkarman/demo/internal/infrastructure/llm/claude"
	"github.com/gkarman/demo/internal/infrastructure/llm/openai"
	"github.com/gkarman/demo/internal/infrastructure/llm/openrouter"
	"github.com/gkarman/demo/internal/infrastructure/shotstack"
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

func NewBrollGenerator(cfg *config.Config) (application.BrollGenerator, error) {
	switch cfg.LLM.Provider {
	case "openai":
		return openai.NewBrollGenerator(cfg.OpenAI.Token, cfg.OpenAI.Model), nil
	default:
		return nil, fmt.Errorf("broll generator: unsupported LLM provider %q, only openai is supported", cfg.LLM.Provider)
	}
}

func NewKlingClient(cfg *config.Config) application.BrollVideoGenerator {
	return kling.NewClient(cfg.Kling.AccessKeyID, cfg.Kling.SecretKey, cfg.Kling.Model)
}

func NewShotstackClient(cfg *config.Config) application.VideoComposer {
	return shotstack.NewClient(cfg.Shotstack.APIKey, cfg.Shotstack.BaseURL)
}
