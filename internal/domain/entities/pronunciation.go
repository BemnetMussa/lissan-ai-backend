package entities

// PracticeSentence is the data for the "Listen" feature, providing pre-generated audio.
type PracticeSentence struct {
	ID              string `json:"id"`
	Text            string `json:"text"`
	CorrectAudioURL string `json:"correct_audio_url"`
}

type MispronouncedWord struct {
	Word        string `json:"word"`
	FeedbackTip string `json:"feedback_tip"`
}

// PronunciationFeedback is the final, structured response for a user's attempt.
type PronunciationFeedback struct {
	OverallAccuracyScore float32             `json:"overall_accuracy_score"`
	MispronouncedWords   []MispronouncedWord `json:"mispronouncedwords"`
	FullFeedbackSummary  string              `json:"full_feedback_summary"`
}

// WordFeedback provides detailed feedback for a single word.
// type WordFeedback struct {
// 	Word          string        `json:"word"`
// 	AccuracyScore float32       `json:"accuracy_score"`
// 	IsCorrect     bool          `json:"is_correct"`
// 	Errors        []ErrorDetail `json:"errors"`
// }

// // ErrorDetail describes a specific phonetic error.
// type ErrorDetail struct {
// 	ExpectedPhoneme string `json:"expected_phoneme"`
// 	ProducedPhoneme string `json:"produced_phoneme"`
// 	FeedbackTip     string `json:"feedback_tip"`
// }

// MFAAlignmentResponse is the struct that matches the JSON response
// from the (real or mock) Python MFA microservice.
type MFAAlignmentResponse struct {
	Words []MFAWord `json:"words"`
}

type MFAWord struct {
	Word     string   `json:"word"`
	Phonemes []string `json:"phonemes"`
}
