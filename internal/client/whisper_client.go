// internal/client/whisper_client.go

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
		client: &http.Client{
			// Setting a reasonable timeout is always a good practice.
			Timeout: 30 * time.Second,
		},
	}
}

func (c *WhisperClient) Transcribe(ctx context.Context, audioData []byte) (string, error) {
	url := "https://api-inference.huggingface.co/models/openai/whisper-large-v3"

	// This retry logic for a loading model is smart. Let's make it robust.
	for i := 0; i < 3; i++ {
		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(audioData))
		if err != nil {
			return "", fmt.Errorf("failed to create whisper request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Content-Type", "audio/ogg")

		resp, err := c.client.Do(req)
		if err != nil {
			return "", fmt.Errorf("failed to call whisper api: %w", err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			resp.Body.Close() // Still try to close it
			return "", fmt.Errorf("failed to read whisper response body: %w", err)
		}

		// The successful case
		if resp.StatusCode == http.StatusOK {
			resp.Body.Close() // We are done with the body, so we close it.
			var wr WhisperResponse
			if err := json.Unmarshal(body, &wr); err != nil {
				return "", fmt.Errorf("failed to parse successful whisper response: %w, body: %s", err, string(body))
			}
			return wr.Text, nil
		}

		// The error case
		var we WhisperErrorResponse
		if err := json.Unmarshal(body, &we); err == nil && strings.Contains(strings.ToLower(we.Error), "loading") {
			// ==========================================================
			// THE FIX: Manually close the body BEFORE sleeping and continuing.
			// ==========================================================
			resp.Body.Close()
			// ==========================================================
			waitTime := time.Duration(we.EstimatedTime+2) * time.Second
			time.Sleep(waitTime)
			continue // Retry the loop
		}

		// Any other unexpected error
		resp.Body.Close() // We are done, so close it.
		return "", fmt.Errorf("unexpected whisper response, status: %s, body: %s", resp.Status, string(body))
	}

	return "", errors.New("model did not load in time after 3 retries")
}