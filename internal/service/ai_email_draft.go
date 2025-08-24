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

// rawEmailResponse matches the AI's JSON output
type rawEmailResponse struct {
	Subject        string `json:"subject"`
	GeneratedEmail string `json:"generated_email"`
}

// aiEmailService is the concrete implementation of EmailService
type aiEmailService struct {
	client *genai.Client
	model  string
}

// NewAIEmailService initializes the AI email service
func NewAIEmailService(apiKey string, model string) (interfaces.EmailService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		return nil, err
	}
	return &aiEmailService{client: client, model: model}, nil
}

// GenerateEmail returns a structured AI output
func (s *aiEmailService) GenerateEmail(ctx context.Context, req *entities.EmailRequest) (*entities.EmailResponse, error) {
	// Build system prompt
	prompt := fmt.Sprintf(`
assume you are an email writing assistant.
the user sends either text about writing an email or a drafted email for correction.
your task is to generate a professional email based on the input.
if the user input is in amharic translate it to english first then generate the email.
if the user input is about generating new email use the prompt to generate the email.
if the user input is a drafted email correct it and make it more professional based on the tone and template type.
we use this for ethiopian users who need help with email drafting and creating.
The user may provide input in English or Amharic. 
1. First, translate if needed so you understand it.
2. Consider the tone: %s (options: polite, friendly, formal).
3. Consider the template type: %s (e.g. job_application, application_followup, complaint, etc.).
4. Generate a professional email in English only.

The JSON object must have exactly two keys: "subject" and "generated_email".

Here is the user input: %s
`, req.Tone, req.TemplateType, req.Prompt)

	result, err := s.client.Models.GenerateContent(ctx, s.model, genai.Text(prompt), nil)
	if err != nil {
		return nil, err
	}

	// Clean up raw AI output
	text := strings.TrimSpace(result.Text())
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	// Parse into raw struct
	var rawResp rawEmailResponse
	if err := json.Unmarshal([]byte(text), &rawResp); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w\nRaw output: %s", err, text)
	}

	// Convert escaped "\n" into real line breaks
	cleanBody := strings.ReplaceAll(rawResp.GeneratedEmail, "\\n", "\n")

	// Build final struct
	emailResp := &entities.EmailResponse{
		Subject:        rawResp.Subject,
		GeneratedEmail: cleanBody,
	}

	return emailResp, nil
}
