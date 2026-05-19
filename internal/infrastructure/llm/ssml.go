package llm

import (
	"encoding/json"
	"fmt"
	"strings"
)

func BuildSSML(rawPayload []byte) (string, error) {
	var payload struct {
		Text  string `json:"text"`
		Words []struct {
			Text  string `json:"text"`
			Start int    `json:"start"`
			End   int    `json:"end"`
		} `json:"words"`
	}

	if err := json.Unmarshal(rawPayload, &payload); err != nil {
		return "", fmt.Errorf("parse payload: %w", err)
	}

	if len(payload.Words) == 0 {
		return "<speak>" + payload.Text + "</speak>", nil
	}

	var sb strings.Builder
	sb.WriteString("<speak>")

	for i, word := range payload.Words {
		if i > 0 {
			gap := word.Start - payload.Words[i-1].End
			if gap >= 400 {
				sb.WriteString(fmt.Sprintf(`<break time="%.1fs"/>`, float64(gap)/1000))
			} else {
				sb.WriteString(" ")
			}
		}
		sb.WriteString(word.Text)
	}

	sb.WriteString("</speak>")
	return sb.String(), nil
}
