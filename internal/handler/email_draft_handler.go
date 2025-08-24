package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"lissanai.com/backend/internal/domain/entities"
	"lissanai.com/backend/internal/domain/interfaces"
)

type EmailController struct {
	emailUC interfaces.EmailUsecase
}

func NewEmailController(emailUC interfaces.EmailUsecase) *EmailController {
	return &EmailController{emailUC: emailUC}
}

func (c *EmailController) GenerateEmailHandler(ctx *gin.Context) {
	var req entities.EmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emailResp, err := c.emailUC.GenerateProfessionalEmail(context.Background(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, emailResp)
}
