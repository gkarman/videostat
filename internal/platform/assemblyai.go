package platform

import (
	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/infrastructure/videoanalyzer"
)

func NewAssemblyAIAnalyzer(cfg *config.Config) *videoanalyzer.AssemblyAIAnalyzer {
	return videoanalyzer.NewAssemblyAIAnalyzer(cfg.Assemblyai.Token)
}
