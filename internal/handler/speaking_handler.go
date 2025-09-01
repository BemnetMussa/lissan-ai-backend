package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"lissanai.com/backend/internal/service"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ControlMessage defines the structure for text-based commands.
type ControlMessage struct {
	Type string `json:"type"`
}

// --- CHANGE #1: THE STRUCT FIELD ---
// The handler now holds the INTERFACE.
type ConversationHandler struct {
	speakingService service.SpeakingService
}

// --- CHANGE #2: THE CONSTRUCTOR PARAMETER ---
// The constructor now accepts the INTERFACE.
func NewConversationHandler(s *service.SpeakingService) *ConversationHandler {
	return &ConversationHandler{speakingService: *s}
}

// The HandleConversation method is correct as you wrote it.
func (h *ConversationHandler) HandleConversation(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close()
	log.Println("Client connected. Waiting for audio stream...")

	var audioBuffer bytes.Buffer

	for {
		msgType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		switch msgType {
		case websocket.BinaryMessage:
			audioBuffer.Write(message)
		case websocket.TextMessage:
			var ctrlMsg ControlMessage
			if err := json.Unmarshal(message, &ctrlMsg); err == nil && ctrlMsg.Type == "end_of_speech" {
				if audioBuffer.Len() == 0 {
					continue
				}
				log.Printf("End-of-speech received. Processing %d bytes.", audioBuffer.Len())

				feedbackAudio, err := h.speakingService.ProcessAudioFeedback(context.Background(), audioBuffer.Bytes())
				if err != nil {
					log.Println("Processing error:", err)
					audioBuffer.Reset()
					continue
				}

				conn.WriteMessage(websocket.BinaryMessage, feedbackAudio)
				audioBuffer.Reset()
			}
		}
	}
}
