package videoanalyzer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://api.assemblyai.com/v2"

type AssemblyAIAnalyzer struct {
	apiKey string
	client *http.Client
}

func NewAssemblyAIAnalyzer(apiKey string) *AssemblyAIAnalyzer {
	return &AssemblyAIAnalyzer{
		apiKey: apiKey,
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

func (a *AssemblyAIAnalyzer) Analyze(ctx context.Context, videoURL string) ([]byte, error) {
	id, err := a.submit(ctx, videoURL)
	if err != nil {
		return nil, err
	}

	raw, err := a.wait(ctx, id)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func (a *AssemblyAIAnalyzer) submit(ctx context.Context, url string) (string, error) {
	body := []byte(fmt.Sprintf(`{
		"audio_url": "%s",
		"auto_highlights": true
	}`, url))

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/transcript", bytes.NewReader(body))
	req.Header.Set("authorization", a.apiKey)
	req.Header.Set("content-type", "application/json")

	res, err := a.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	b, _ := io.ReadAll(res.Body)

	id := extractID(b)
	return id, nil
}

func (a *AssemblyAIAnalyzer) wait(ctx context.Context, id string) ([]byte, error) {
	for {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/transcript/"+id, nil)
		req.Header.Set("authorization", a.apiKey)

		res, err := a.client.Do(req)
		if err != nil {
			return nil, err
		}

		raw, _ := io.ReadAll(res.Body)
		res.Body.Close()

		if isCompleted(raw) {
			return raw, nil
		}

		if isError(raw) {
			return nil, fmt.Errorf("assemblyai transcription error")
		}

		time.Sleep(3 * time.Second)
	}
}

func extractID(b []byte) string {
	const key = `"id":"`
	i := bytes.Index(b, []byte(key))
	if i == -1 {
		return ""
	}
	start := i + len(key)
	end := bytes.IndexByte(b[start:], '"')
	return string(b[start : start+end])
}

func isCompleted(b []byte) bool {
	return bytes.Contains(b, []byte(`"status":"completed"`))
}

func isError(b []byte) bool {
	return bytes.Contains(b, []byte(`"status":"error"`))
}
