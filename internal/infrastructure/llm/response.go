package llm

import (
	"encoding/json"
	"fmt"
)

type GeneratedContent struct {
	Brief  json.RawMessage `json:"brief"`
	Script string          `json:"script"`
}

func ParseGeneratedContent(raw string) (*GeneratedContent, error) {
	var content GeneratedContent
	if err := json.Unmarshal([]byte(raw), &content); err != nil {
		return nil, fmt.Errorf("parse llm response: %w", err)
	}
	if len(content.Brief) == 0 {
		return nil, fmt.Errorf("llm response missing brief")
	}
	if content.Script == "" {
		return nil, fmt.Errorf("llm response missing script")
	}
	return &content, nil
}
