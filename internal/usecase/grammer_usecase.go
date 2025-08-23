package usecase

import (
	"lissanai.com/backend/internal/domain/interfaces"
	"lissanai.com/backend/internal/domain/models"
)

type GrammarUsecase struct {
	ai_service interfaces.AiServiceInterface
}

func NewGrammerUsecase(ai_service interfaces.AiServiceInterface) *GrammarUsecase {
	return &GrammarUsecase{ai_service: ai_service}
}

func (au *GrammarUsecase) CheckGrammer(text string) (*models.GrammarResponse, error) {
	response, err := au.ai_service.CheckGrammar(text)
	if err != nil {
		return nil, err
	}
	return response, nil
}
