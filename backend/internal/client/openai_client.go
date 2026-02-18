package client

import (
	"context"
)

type OpenAIClient struct {
	apiKey string
}

func NewOpenAIClient(key string) *OpenAIClient {
	return &OpenAIClient{apiKey: key}
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
	// implement OpenAI call later
	return "", nil
}