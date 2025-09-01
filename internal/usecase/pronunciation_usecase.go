package usecase

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"lissanai.com/backend/internal/domain/entities"
	"lissanai.com/backend/internal/domain/interfaces"
)

type pronunciationUsecase struct {
	mfaClient  interfaces.MFAClient
	pronunDict map[string][]string
}

// NewPronunciationUsecase is the constructor. It accepts the interface and loads the dictionary.
func NewPronunciationUsecase(mfaClient interfaces.MFAClient) (interfaces.PronunciationUsecase, error) {
	dict, err := loadPronunciationDictionary("cmudict.dict") // Assumes cmudict.dict is in the project root
	if err != nil {
		return nil, fmt.Errorf("could not load pronunciation dictionary: %w", err)
	}
	return &pronunciationUsecase{
		mfaClient:  mfaClient,
		pronunDict: dict,
	}, nil
}

// GetPracticeSentences returns the static list of sentences and their audio URLs.
func (uc *pronunciationUsecase) GetPracticeSentences() []*entities.PracticeSentence {
	return []*entities.PracticeSentence{
		{ID: "ps_001", Text: "She sells seashells by the seashore.", CorrectAudioURL: "https://your-cdn.com/audio/sentence1.mp3"},
		{ID: "ps_002", Text: "The quick brown fox jumps over the lazy dog.", CorrectAudioURL: "https://your-cdn.com/audio/sentence2.mp3"},
	}
}

// ... (The struct definition, NewPronunciationUsecase, GetPracticeSentences,
//      and loadPronunciationDictionary functions are all the same)

// AssessPronunciation now orchestrates the analysis and builds the SIMPLER response.
func (uc *pronunciationUsecase) AssessPronunciation(ctx context.Context, targetText string, audioData []byte) (*entities.PronunciationFeedback, error) {
	producedAlignment, err := uc.mfaClient.GetPhoneticAlignment(ctx, audioData, targetText)
	if err != nil {
		return nil, fmt.Errorf("MFA alignment failed: %w", err)
	}

	targetWords := strings.Fields(strings.ToLower(strings.Trim(targetText, ".?!")))

	// Initialize our new response struct
	feedback := &entities.PronunciationFeedback{
		MispronouncedWords: make([]entities.MispronouncedWord, 0),
	}

	// var totalScore float32 = 0
	var correctWords int = 0
	var feedbackSummary strings.Builder
	feedbackSummary.WriteString("Good effort! ")

	producedMap := make(map[string][]string)
	for _, wordData := range producedAlignment.Words {
		producedMap[strings.ToLower(wordData.Word)] = wordData.Phonemes
	}

	for _, wordStr := range targetWords {
		expectedPhonemes, ok := uc.pronunDict[wordStr]
		if !ok { // If word is not in dictionary, assume it's correct for now.
			correctWords++
			continue
		}

		producedPhonemes, ok := producedMap[wordStr]
		if !ok { // If MFA didn't detect the word, it's a mistake.
			continue
		}

		// --- CORE LOGIC CHANGE ---
		// We will now only check for the first mistake in a word.
		isCorrect := true
		for i := 0; i < len(expectedPhonemes) && i < len(producedPhonemes); i++ {
			if expectedPhonemes[i] != producedPhonemes[i] {
				isCorrect = false
				tip := fmt.Sprintf("For the word '%s', the sound for '%s' was a bit unclear. Try to make a clear '%s' sound.", wordStr, producedPhonemes[i], expectedPhonemes[i])

				// Add the word to our new, simpler list
				feedback.MispronouncedWords = append(feedback.MispronouncedWords, entities.MispronouncedWord{
					Word:        wordStr,
					FeedbackTip: tip,
				})

				feedbackSummary.WriteString(tip + " ")
				break // Stop checking this word after the first error
			}
		}

		if isCorrect && len(expectedPhonemes) == len(producedPhonemes) {
			correctWords++
		}
	}

	// Calculate a simpler overall score.
	if len(targetWords) > 0 {
		feedback.OverallAccuracyScore = (float32(correctWords) / float32(len(targetWords))) * 100.0
	}

	if len(feedback.MispronouncedWords) == 0 {
		feedbackSummary.WriteString("Everything sounded great!")
	}

	feedback.FullFeedbackSummary = feedbackSummary.String()

	return feedback, nil
}

// loadPronunciationDictionary reads the CMUDict file into a map.
func loadPronunciationDictionary(filepath string) (map[string][]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	dict := make(map[string][]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, ";;;") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		word := strings.ToLower(parts[0])
		phonemes := parts[1:]
		dict[word] = phonemes
	}
	log.Printf("Successfully loaded pronunciation dictionary with %d words.", len(dict))
	return dict, scanner.Err()
}
