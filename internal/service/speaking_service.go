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
	groqClient         *client.GroqClient
	whisperClient      *client.WhisperClient
	unrealSpeechClient *client.UnrealSpeechTTSClient // 1. Renamed the field
}

// 2. Updated the function signature to accept the new client type
func NewSpeakingService(groq *client.GroqClient, whisper *client.WhisperClient, unreal *client.UnrealSpeechTTSClient) SpeakingService {
	return &speakingServiceImpl{
		groqClient:         groq,
		whisperClient:      whisper,
		unrealSpeechClient: unreal, // 3. Updated the assignment
	}
}


func (s *speakingServiceImpl) ProcessAudioFeedback(ctx context.Context, audioData []byte) ([]byte, error) {
	// 1️⃣ STT: Convert audio to text
	text, err := s.whisperClient.Transcribe(ctx, audioData)
	if err != nil {
		return nil, fmt.Errorf("STT error: %w", err)
	}
	if text == "" {
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
	// 4. Use the new client to generate audio. The method call is the same!
	ttsAudio, err := s.unrealSpeechClient.GenerateAudio(cleanedResponse)
	if err != nil {
		return nil, fmt.Errorf("TTS error: %w", err)
	}

	return ttsAudio, nil
}