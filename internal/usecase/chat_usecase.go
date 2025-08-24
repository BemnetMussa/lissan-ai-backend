package usecase

import "lissanai.com/backend/internal/service"

type ChatUsecase struct {
	ai *service.ChatAIService
}

func NewChatUsecase(ai *service.ChatAIService) *ChatUsecase {
	return &ChatUsecase{ai: ai}
}

func (u *ChatUsecase) StartSession() string {
	return u.ai.NewInterviewSession()
}

func (u *ChatUsecase) GetNextQuestion(sessionID string) string {
	return u.ai.AskQuestion(sessionID)
}

func (u *ChatUsecase) EvaluateAnswer(sessionID, answer string) string {
	return u.ai.EvaluateAnswer(sessionID, answer)
}
