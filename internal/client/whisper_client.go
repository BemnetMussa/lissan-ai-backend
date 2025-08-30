package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type WhisperClient struct {
	apiKey string
	client *http.Client
}

type WhisperResponse struct {
	Text string `json:"text"`
}
type WhisperErrorResponse struct {
	Error         string  `json:"error"`
	EstimatedTime float64 `json:"estimated_time,omitempty"`
}

func NewWhisperClient(apiKey string) *WhisperClient {
	return &WhisperClient{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

func (c *WhisperClient) Transcribe(ctx context.Context, audioData []byte) (string, error) {
	url := "https://api-inference.huggingface.co/models/openai/whisper-large-v3"

	for i := 0; i < 3; i++ {
		req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(audioData))
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Content-Type", "audio/ogg")

		resp, err := c.client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode == 200 {
			var wr WhisperResponse
			if err := json.Unmarshal(body, &wr); err != nil {
				return "", fmt.Errorf("parse error: %w, body: %s", err, string(body))
			}
			return wr.Text, nil
		}

		var we WhisperErrorResponse
		if err := json.Unmarshal(body, &we); err == nil && strings.Contains(strings.ToLower(we.Error), "loading") {
			waitTime := time.Duration(we.EstimatedTime+2) * time.Second
			time.Sleep(waitTime)
			continue
		}
		return "", errors.New("unexpected Whisper response")
	}
	return "", errors.New("model did not load in time")
}
