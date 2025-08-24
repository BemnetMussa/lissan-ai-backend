package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"lissanai.com/backend/internal/usecase"
)

type ChatHandler struct {
	usecase *usecase.ChatUsecase
}

func NewChatHandler(u *usecase.ChatUsecase) *ChatHandler {
	return &ChatHandler{usecase: u}
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func (h *ChatHandler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	sessionID := h.usecase.StartSession()
	firstQuestion := h.usecase.GetNextQuestion(sessionID)
	conn.WriteJSON(map[string]string{"type": "ai_question", "text": firstQuestion})

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		userAnswer := string(msg)

		// async feedback
		go func(answer string) {
			time.Sleep(1 * time.Second) // simulate AI thinking
			feedback := h.usecase.EvaluateAnswer(sessionID, answer)
			conn.WriteJSON(map[string]string{"type": "ai_feedback", "text": feedback})

			nextQ := h.usecase.GetNextQuestion(sessionID)
			conn.WriteJSON(map[string]string{"type": "ai_question", "text": nextQ})
		}(userAnswer)
	}
}
