package service

import (
	"context"
	"fmt"
	"strings"

	"lissanai.com/backend/internal/client"
)

type SpeakingService struct {
	groqClient       *client.GroqClient
	whisperClient    *client.WhisperClient
	elevenLabsClient *client.ElevenLabsTTSClient
}

func NewSpeakingService(groq *client.GroqClient, whisper *client.WhisperClient, eleven *client.ElevenLabsTTSClient) *SpeakingService {
	return &SpeakingService{
		groqClient:       groq,
		whisperClient:    whisper,
		elevenLabsClient: eleven,
	}
}

func (s *SpeakingService) ProcessAudioFeedback(ctx context.Context, audioData []byte) ([]byte, error) {
	// 1️⃣ STT: Convert audio to text
	text, err := s.whisperClient.Transcribe(ctx, audioData)
	if err != nil {
		return nil, fmt.Errorf("STT error: %w", err)
	}

	// 2️⃣ LLM: Generate response
	response, err := s.groqClient.GenerateContent(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("LLM error: %w", err)
	}

	// 🧹 Clean and trim response before TTS
	cleanedResponse := strings.TrimSpace(response)
	if len(cleanedResponse) > 100 {
		cleanedResponse = cleanedResponse[:100]
	}

	// 3️⃣ TTS: Convert response to audio
	ttsAudio, err := s.elevenLabsClient.GenerateAudio(cleanedResponse)
	if err != nil {
		return nil, fmt.Errorf("TTS error: %w", err)
	}

	return ttsAudio, nil
}
