package server

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"lissanai.com/backend/internal/client"
	"lissanai.com/backend/internal/handler"
	"lissanai.com/backend/internal/service"
)

func SetupSpeakingRoutes(router *gin.RouterGroup) {
	godotenv.Load() // optional .env

	groqAPIKey := os.Getenv("GROQ_API_KEY")
	hfAPIKey := os.Getenv("HF_API_KEY")
	elevenLabsKey := os.Getenv("ELEVENLABS_API_KEY")
	voiceID := os.Getenv("ELEVENLABS_VOICE_ID")

	if groqAPIKey == "" || hfAPIKey == "" || elevenLabsKey == "" || voiceID == "" {
		log.Fatal("Missing one or more API keys or voice ID")
	}

	groqClient := client.NewGroqClient(groqAPIKey)
	whisperClient := client.NewWhisperClient(hfAPIKey)
	elevenLabsClient := client.NewElevenLabsTTSClient(elevenLabsKey, voiceID)

    // This now returns the interface type, which is what the handler expects.
	speakingService := service.NewSpeakingService(groqClient, whisperClient, elevenLabsClient)
	conversationHandler := handler.NewConversationHandler(speakingService)

	router.GET("/ws/conversation", conversationHandler.HandleConversation)
}