package blogger

import "time"

type VideoPrompt struct {
	ID          string
	VideoID     string
	LLMProvider string
	LLMModel    string
	Brief       string
	Script      string
	CreatedAt   time.Time
}
