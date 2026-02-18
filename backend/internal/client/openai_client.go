package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	domainerr "S.P.A.R.T.A/backend/internal/domain/errors"
)

type OpenAIClient struct {
	apiKey string
	model  string
	base   string
	http   *http.Client
}

func NewOpenAIClient(key, model, baseURL string) *OpenAIClient {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "https://api.openai.com"
	}
	return &OpenAIClient{
		apiKey: key,
		model:  model,
		base:   strings.TrimRight(baseURL, "/"),
		http: &http.Client{
			Timeout: 25 * time.Second,
		},
	}
}

type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

func (c *OpenAIClient) Generate(ctx context.Context, prompt string) (string, error) {
	if strings.TrimSpace(c.apiKey) == "" {
		return "", domainerr.ErrInternal
	}

	reqBody := ChatRequest{
		Model: c.model,
		Messages: []ChatMessage{
			{Role: "system", Content: "Return ONLY valid JSON. No markdown, no code fences."},
			{Role: "user", Content: prompt},
		},
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", domainerr.ErrInternal
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+"/v1/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return "", domainerr.ErrInternal
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", domainerr.ErrInternal
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", domainerr.ErrInternal
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", domainerr.ErrInternal
	}

	var decoded ChatResponse
	if err := json.Unmarshal(body, &decoded); err != nil {
		return "", domainerr.ErrInternal
	}
	if len(decoded.Choices) == 0 {
		return "", domainerr.ErrInternal
	}

	return decoded.Choices[0].Message.Content, nil
}
