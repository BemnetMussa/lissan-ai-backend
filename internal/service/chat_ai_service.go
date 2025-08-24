package service

import "fmt"

type ChatAIService struct {
	// could store API key or model context
}

func NewChatAIService() *ChatAIService {
	return &ChatAIService{}
}

func (s *ChatAIService) NewInterviewSession() string {
	// return session ID (could be UUID)
	return "session-123"
}

func (s *ChatAIService) AskQuestion(sessionID string) string {
	// return next question
	return "Tell me about a project you worked on recently."
}

func (s *ChatAIService) EvaluateAnswer(sessionID, answer string) string {
	if answer == "" {
		return "I didn't catch that. Could you try again?"
	}
	return fmt.Sprintf("Feedback: '%s' is a good start, but structure it with intro, example, conclusion.", answer)
}
