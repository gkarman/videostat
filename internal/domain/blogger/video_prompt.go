package blogger

import "time"

type VideoPrompt struct {
	ID          string
	VideoID     string
	LLMProvider string
	LLMModel    string
	Prompt      string
	CreatedAt   time.Time
}
