package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io" // Use io instead of the deprecated ioutil
	"net/http"
	"os"
	"os/exec"
)

// UnrealSpeechTTSClient holds the configuration for the Unreal Speech API client.
type UnrealSpeechTTSClient struct {
	apiKey  string
	voiceID string
	client  *http.Client
}

// unrealSpeechRequest defines the JSON body for the Unreal Speech API request.
// Note the JSON tags match the API's required capitalization.
type unrealSpeechRequest struct {
	Text    string `json:"Text"`
	VoiceId string `json:"VoiceId"`
	Bitrate string `json:"Bitrate,omitempty"` // omitempty means it won't be included if it's empty
	Speed   string `json:"Speed,omitempty"`
}

// NewUnrealSpeechTTSClient creates a new client for the Unreal Speech API.
func NewUnrealSpeechTTSClient(apiKey, voiceID string) *UnrealSpeechTTSClient {
	return &UnrealSpeechTTSClient{
		apiKey:  apiKey,
		voiceID: voiceID,
		client:  &http.Client{},
	}
}

// GenerateAudio connects to the Unreal Speech API and returns the audio data as bytes.
func (c *UnrealSpeechTTSClient) GenerateAudio(text string) ([]byte, error) {
	// The new, correct V8 API endpoint
	url := "https://api.v8.unrealspeech.com/stream"

	// Create the request payload with the new structure
	payload := unrealSpeechRequest{
		Text:    text,
		VoiceId: c.voiceID, // The voice ID is now part of the body
		Bitrate: "192k",    // A sensible default
		Speed:   "0",
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set the new, correct headers for Unreal Speech
	// It uses a Bearer token instead of 'xi-api-key'
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request to Unreal Speech API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { // Check for 200 OK
		// Try to read the error message from the API response
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Unreal Speech API error: received status code %d - %s", resp.StatusCode, string(bodyBytes))
	}

	// Use io.ReadAll which is the modern standard
	return io.ReadAll(resp.Body)
}

// PlayAudio does not need to change at all!
// It already accepts a byte slice of audio, so it will work perfectly with the
// output from our new GenerateAudio function.
func (c *UnrealSpeechTTSClient) PlayAudio(audioData []byte) error {
	tmpFile, err := os.CreateTemp("", "tts_*.mp3")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(audioData); err != nil {
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}

	// This command works for any MP3 file, regardless of origin
	cmd := exec.Command("ffplay", "-autoexit", "-nodisp", tmpFile.Name())
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}