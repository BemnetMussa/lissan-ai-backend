package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lissanai.com/backend/internal/usecase"
)

type GrammarCheck struct {
	grammer_usecase usecase.GrammarUsecase
}

func NewGrammarHandler(grammer_usecase usecase.GrammarUsecase) *GrammarCheck {
	return &GrammarCheck{
		grammer_usecase: grammer_usecase,
	}
}

func (h *GrammarCheck) GrammarCheck(c *gin.Context) {
	var request struct {
		Text string `json:"text" binding:"required"`
	}

	// Validate request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Call usecase/service
	grammarResp, err := h.grammer_usecase.CheckGrammer(request.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return JSON response
	c.JSON(http.StatusOK, grammarResp)
}
