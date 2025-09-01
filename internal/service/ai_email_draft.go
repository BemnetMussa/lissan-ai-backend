package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/genai"
	"lissanai.com/backend/internal/domain/entities"
	"lissanai.com/backend/internal/domain/interfaces"
)

// aiEmailService is the private implementation of the EmailService interface.
type aiEmailService struct {
	client *genai.Client
	model  string
}

// NewAIEmailService is the public constructor.
func NewAIEmailService(apiKey string, model string) (interfaces.EmailService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		return nil, err
	}
	return &aiEmailService{client: client, model: model}, nil
}

// GenerateEmailFromPrompt handles the logic for creating a new email.
func (s *aiEmailService) GenerateEmailFromPrompt(ctx context.Context, req *entities.GenerateEmailRequest) (*entities.EmailResponse, error) {
	prompt := fmt.Sprintf(`
Your task is to generate a new, complete, professional English email.
The user's request might be in English or Amharic.
Consider the tone: %s and the template type: %s.
Your response MUST be a single, minified JSON object with two keys: "subject" and "body".
Do not include any introductory text or code fences.
User's Request: %s`,
		req.Tone, req.TemplateType, req.Prompt)

	return s.callAIAndParseResponse(ctx, prompt)
}

// EditEmailDraft handles the logic for correcting an existing email.
func (s *aiEmailService) EditEmailDraft(ctx context.Context, req *entities.EditEmailRequest) (*entities.EmailResponse, error) {
	prompt := fmt.Sprintf(`
Your task is to correct and improve an existing email draft to make it more professional.
Fix all grammatical errors, improve the tone, and enhance clarity.
Consider the desired tone: %s and template type: %s.
Your response MUST be a single, minified JSON object with two keys: "subject" and "body".
Do not include any introductory text or code fences.
User's Email Draft: %s`,
		req.Tone, req.TemplateType, req.Draft)

	return s.callAIAndParseResponse(ctx, prompt)
}

// callAIAndParseResponse is a private helper to avoid duplicating code.
func (s *aiEmailService) callAIAndParseResponse(ctx context.Context, prompt string) (*entities.EmailResponse, error) {
	result, err := s.client.Models.GenerateContent(ctx, s.model, genai.Text(prompt), nil)
	if err != nil {
		return nil, err
	}

	text := strings.TrimSpace(result.Text())
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimSuffix(text, "```")

	var emailResp entities.EmailResponse
	if err := json.Unmarshal([]byte(text), &emailResp); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w\nRaw output: %s", err, text)
	}

	return &emailResp, nil
}
