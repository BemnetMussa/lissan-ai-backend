package interfaces

import (
	"context"

	"lissanai.com/backend/internal/domain/entities"
)

// MFAClient defines the contract for any client that performs phonetic alignment.
// Both our real HTTP client and our mock client will implement this.
type MFAClient interface {
	GetPhoneticAlignment(ctx context.Context, audioData []byte, targetText string) (*entities.MFAAlignmentResponse, error)
}

// PronunciationUsecase defines the contract for the pronunciation business logic.
type PronunciationUsecase interface {
	GetPracticeSentences() []*entities.PracticeSentence
	AssessPronunciation(ctx context.Context, targetText string, audioData []byte) (*entities.PronunciationFeedback, error)
}
