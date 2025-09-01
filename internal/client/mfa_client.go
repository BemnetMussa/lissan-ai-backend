package client

import (
	"context"

	"lissanai.com/backend/internal/domain/entities"
	"lissanai.com/backend/internal/domain/interfaces"
)

// MockMFAClient is a fake client for testing purposes. It does NOT make network calls.
type MockMFAClient struct{}

// NewMockMFAClient is the constructor for our fake client.
// It returns the MFAClient INTERFACE.
func NewMockMFAClient() interfaces.MFAClient {
	return &MockMFAClient{}
}

// GetPhoneticAlignment pretends to process the audio and returns hardcoded data.
func (c *MockMFAClient) GetPhoneticAlignment(ctx context.Context, audioData []byte, targetText string) (*entities.MFAAlignmentResponse, error) {
	// This simulates the MFA service correctly analyzing the phrase "She sells seashells".
	// We are intentionally simulating a mistake for the word "sells" (SH instead of S).
	mockResponse := &entities.MFAAlignmentResponse{
		Words: []entities.MFAWord{
			{Word: "she", Phonemes: []string{"SH", "IY"}},
			{Word: "sells", Phonemes: []string{"SH", "EH", "L", "Z"}}, // Simulated mistake
			{Word: "seashells", Phonemes: []string{"S", "IY", "SH", "EH", "L", "Z"}},
		},
	}
	return mockResponse, nil
}
