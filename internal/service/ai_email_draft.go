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
	GeneratedEmail string `json:"body"`
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
func (s *aiEmailService) ProcessEmail(ctx context.Context, req *entities.EmailRequest) (*entities.EmailResponse, error) {
	var prompt string

	switch strings.ToUpper(req.Type) {
	case "GENERATE":
		prompt = fmt.Sprintf(`
You are an expert email writing assistant for Ethiopian professionals who may use English as a second language.
Your task is to generate a new, complete, professional email based on the user's request.
The user's request might be in English or Amharic; handle both appropriately.
The final email must be in English.

Consider the following if there exist:
- Tone: %s
- Template Type: %s

Respond ONLY with a single, minified JSON object. The JSON object must have exactly two keys: "subject" and "body".

User's Request: %s
`, req.Tone, req.TemplateType, req.Prompt)
	case "EDIT":
		prompt = fmt.Sprintf(`
You are an expert English communication coach for Ethiopian professionals.
Your task is to correct and improve an existing email draft to make it more professional.
Fix all grammatical errors, improve the tone, and enhance clarity.
The final, improved email must be in English.

Consider the following if there exist:
- Tone: %s
- Template Type: %s

Respond ONLY with a single, minified JSON object. The JSON object must have exactly two keys: "subject" and "body".

User's Email Draft: %s
`, req.Tone, req.TemplateType, req.Prompt)

	default:
		// Handle cases where the type is missing or invalid
		return nil, fmt.Errorf("invalid request type: '%s'. Must be 'GENERATE' or 'EDIT'", req.Type)
	}

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

	// Parse into our internal struct
	var rawResp rawEmailResponse
	if err := json.Unmarshal([]byte(text), &rawResp); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w\nRaw output: %s", err, text)
	}

	// Convert escaped "\n" into real line breaks
	cleanBody := strings.ReplaceAll(rawResp.GeneratedEmail, "\\n", "\n")

	// Build the final response struct to send back to the handler
	emailResp := &entities.EmailResponse{
		Subject:        rawResp.Subject,
		GeneratedEmail: cleanBody,
	}

	return emailResp, nil
}
