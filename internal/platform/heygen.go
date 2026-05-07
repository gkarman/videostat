package platform

import (
	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/infrastructure/videogenerator/heygen"
)

func NewHeyGenClient(cfg *config.Config) *heygen.Client {
	return heygen.NewClient(cfg.HeyGen.APIKey, cfg.HeyGen.AvatarID, cfg.HeyGen.VoiceID)
}
