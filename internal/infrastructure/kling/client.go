package kling

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gkarman/demo/internal/application"
)

const baseURL = "https://api.klingai.com"

type Client struct {
	accessKeyID string
	secretKey   string
	model       string
	http        *http.Client
}

func NewClient(accessKeyID, secretKey, model string) *Client {
	return &Client{
		accessKeyID: accessKeyID,
		secretKey:   secretKey,
		model:       model,
		http:        &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) Submit(ctx context.Context, prompt string, durationSec int) (string, error) {
	body, err := json.Marshal(map[string]any{
		"model_name":      c.model,
		"prompt":          prompt,
		"negative_prompt": "slow motion, slow-mo, cinematic slow, time lapse, fade in, fade out, smooth dolly, peaceful, calm, static shot, still camera, CGI, animation, 3D render, video game, cartoon, Asian people, Asian faces, Chinese people, Japanese people, Korean people, Asian locations, China, Japan, Korea, Asian city, Asian street, non-English signs, Chinese text, subtitles, captions, watermark, Asian characters, hanzi, kanji",
		"duration":        fmt.Sprintf("%d", durationSec),
		"aspect_ratio":    "9:16",
	})
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/v1/videos/text2video", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	c.setAuth(req)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("kling api call: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("kling api error %d: %s", resp.StatusCode, respBody)
	}

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			TaskID string `json:"task_id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}
	if result.Code != 0 {
		return "", fmt.Errorf("kling error %d: %s", result.Code, result.Message)
	}

	return result.Data.TaskID, nil
}

func (c *Client) GetStatus(ctx context.Context, externalID string) (*application.BrollGenerationStatus, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/v1/videos/text2video/"+externalID, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	c.setAuth(req)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("kling api call: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kling api error %d: %s", resp.StatusCode, respBody)
	}

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"message"`
		Data struct {
			TaskStatus string `json:"task_status"`
			TaskResult *struct {
				Videos []struct {
					URL string `json:"url"`
				} `json:"videos"`
			} `json:"task_result"`
			TaskStatusMsg string `json:"task_status_msg"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	status := &application.BrollGenerationStatus{}
	switch result.Data.TaskStatus {
	case "succeed":
		status.Status = "done"
		if result.Data.TaskResult != nil && len(result.Data.TaskResult.Videos) > 0 {
			status.VideoURL = result.Data.TaskResult.Videos[0].URL
		}
	case "failed":
		status.Status = "failed"
		status.Error = result.Data.TaskStatusMsg
	default:
		status.Status = "processing"
	}

	return status, nil
}

func (c *Client) setAuth(req *http.Request) {
	token := c.buildJWT()
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
}

func (c *Client) buildJWT() string {
	now := time.Now().Unix()

	header := base64url([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload := base64url([]byte(fmt.Sprintf(`{"iss":"%s","exp":%d,"nbf":%d}`, c.accessKeyID, now+1800, now-5)))

	sig := hmacSHA256(c.secretKey, header+"."+payload)
	return header + "." + payload + "." + sig
}

func base64url(data []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
}

func hmacSHA256(secret, data string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return strings.TrimRight(base64.URLEncoding.EncodeToString(h.Sum(nil)), "=")
}
