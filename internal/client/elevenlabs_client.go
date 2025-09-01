package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

type ElevenLabsTTSClient struct {
	apiKey  string
	voiceID string
	client  *http.Client
}

type ttsRequest struct {
	Text          string         `json:"text"`
	ModelID       string         `json:"model_id"`
	VoiceSettings *voiceSettings `json:"voice_settings"`
}

type voiceSettings struct {
	Stability       float64 `json:"stability"`
	SimilarityBoost float64 `json:"similarity_boost"`
}

func NewElevenLabsTTSClient(apiKey, voiceID string) *ElevenLabsTTSClient {
	return &ElevenLabsTTSClient{
		apiKey:  apiKey,
		voiceID: voiceID,
		client:  &http.Client{},
	}
}

func (c *ElevenLabsTTSClient) GenerateAudio(text string) ([]byte, error) {
	url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s", c.voiceID)
	payload := ttsRequest{
		Text:    text,
		ModelID: "eleven_multilingual_v2",
		VoiceSettings: &voiceSettings{
			Stability:       0.5,
			SimilarityBoost: 0.75,
		},
	}
	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("xi-api-key", c.apiKey)
	req.Header.Set("Accept", "audio/mpeg")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("ElevenLabs API error: %s", string(b))
	}

	return ioutil.ReadAll(resp.Body)
}

func (c *ElevenLabsTTSClient) PlayAudio(audioData []byte) error {
	tmpFile, err := ioutil.TempFile("", "tts_*.mp3")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.Write(audioData)
	tmpFile.Close()

	cmd := exec.Command("ffplay", "-autoexit", "-nodisp", tmpFile.Name())
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
