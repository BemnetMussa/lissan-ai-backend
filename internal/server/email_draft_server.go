// internal/server/email_draft_server.go
package server

import (
	"os"

	"github.com/gin-gonic/gin"
	"lissanai.com/backend/internal/domain/interfaces"
	"lissanai.com/backend/internal/handler"
	"lissanai.com/backend/internal/service"
	"lissanai.com/backend/internal/usecase"
)

// SetupEmailRoutes sets up the email generation endpoint
func SetupEmailRoutes(router *gin.RouterGroup) interfaces.EmailUsecase {
	// 1. Initialize the AI email service
	emailService, err := service.NewAIEmailService(os.Getenv("GEMINI_API_KEY"), "gemini-2.5-flash")
	if err != nil {
		panic(err) // fail fast if API key or client setup fails
	}

	// 2. Initialize the usecase
	emailUC := usecase.NewEmailUsecase(emailService)

	// 3. Initialize the controller
	emailController := handler.NewEmailController(emailUC)

	// 4. Define the route
	router.POST("/generate-email", emailController.GenerateEmailHandler)

	return emailUC
}
