package platform

import (
	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/infrastructure/llm/claude"
)

func NewClaudePromptGenerator(cfg *config.Config) *claude.PromptGenerator {
	return claude.NewPromptGenerator(cfg.Anthropic.Token, cfg.Anthropic.Model)
}
