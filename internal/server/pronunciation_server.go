package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"lissanai.com/backend/internal/client"
	"lissanai.com/backend/internal/handler"
	"lissanai.com/backend/internal/usecase"
)

func SetupPronunciationRoutes(router *gin.RouterGroup) {
	// --- Using the MOCK client for now ---
	// When your friend's service is ready, you will swap this one line.
	// mfaClient := client.NewMFAClient(os.Getenv("MFA_SERVICE_URL"))
	mockMFAClient := client.NewMockMFAClient()

	pronunciationUC, err := usecase.NewPronunciationUsecase(mockMFAClient)
	if err != nil {
		log.Fatalf("FATAL: Could not create pronunciation usecase: %v", err)
	}

	pronunciationHandler := handler.NewPronunciationHandler(pronunciationUC)

	pronunciationRoutes := router.Group("/pronunciation")
	{
		pronunciationRoutes.GET("/sentences", pronunciationHandler.GetSentences)
		pronunciationRoutes.POST("/assess", pronunciationHandler.AssessPronunciation)
	}
}
