// internal/service/speaking_service.go

package service

import (
	"context"
	"fmt"
	"strings"

	"lissanai.com/backend/internal/client"
	// If you moved shared interfaces/types to internal/common, import it here:
	// "lissanai.com/backend/internal/common"
)

type SpeakingService interface {
	ProcessAudioFeedback(ctx context.Context, audioData []byte) ([]byte, error)
}


type speakingServiceImpl struct {
	groqClient       *client.GroqClient
	whisperClient    *client.WhisperClient
	elevenLabsClient *client.ElevenLabsTTSClient
}

func NewSpeakingService(groq *client.GroqClient, whisper *client.WhisperClient, eleven *client.ElevenLabsTTSClient) SpeakingService {
	return &speakingServiceImpl{
		groqClient:       groq,
		whisperClient:    whisper,
		elevenLabsClient: eleven,
	}
}


func (s *speakingServiceImpl) ProcessAudioFeedback(ctx context.Context, audioData []byte) ([]byte, error) {
	// 1️⃣ STT: Convert audio to text
	// Pass the context to your client for potential cancellation.
	text, err := s.whisperClient.Transcribe(ctx, audioData)
	if err != nil {
		return nil, fmt.Errorf("STT error: %w", err)
	}
	if text == "" {
		// Handle cases of silence gracefully
		return nil, fmt.Errorf("no speech detected in audio")
	}
	fmt.Println("transcripted text: ", text)

	//2️⃣ LLM: Generate response
	response, err := s.groqClient.GenerateContent(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("LLM error: %w", err)
	}
	

	// 🧹 Clean and trim response before TTS
	cleanedResponse := strings.TrimSpace(response)
	if len(cleanedResponse) == 0 {
		return nil, fmt.Errorf("LLM generated an empty response")
	}

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