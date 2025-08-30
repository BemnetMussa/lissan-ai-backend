package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type GroqClient struct {
	apiKey string
	client *http.Client
}

type GroqRequest struct {
	Model    string        `json:"model"`
	Messages []GroqMessage `json:"messages"`
}
type GroqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type GroqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewGroqClient(apiKey string) *GroqClient {
	return &GroqClient{apiKey: apiKey, client: &http.Client{}}
}

func (c *GroqClient) GenerateContent(ctx context.Context, prompt string) (string, error) {
	reqBody := GroqRequest{
		Model:    "llama3-8b-8192",
		Messages: []GroqMessage{{Role: "user", Content: prompt}},
	}
	data, _ := json.Marshal(reqBody)
	req, _ := http.NewRequestWithContext(ctx, "POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var gr GroqResponse
	if err := json.NewDecoder(resp.Body).Decode(&gr); err != nil {
		return "", err
	}
	if len(gr.Choices) > 0 {
		return gr.Choices[0].Message.Content, nil
	}
	return "", errors.New("no content returned")
}
