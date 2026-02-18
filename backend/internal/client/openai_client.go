package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
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

	reqBody := ChatRequest{
		Model: "gpt-4.1-mini",
		Messages: []ChatMessage{
			{Role: "user", Content: prompt},
		},
	}

	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://api.openai.com/v1/chat/completions",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if len(res.Choices) == 0 {
		return "", nil
	}

	return res.Choices[0].Message.Content, nil
}
