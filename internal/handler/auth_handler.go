// internal/handler/auth_handler.go
package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"lissanai.com/backend/internal/domain"
	"lissanai.com/backend/internal/usecase"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

// SignUp godoc
// @Summary      Register a new user
// @Description  Creates a new user account. All fields are required.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        userInfo body domain.SignUpRequest true "User Signup Information"
// @Success      201 {object} domain.User
// @Failure 400 {object} domain.ErrorResponse
// @Router       /auth/signup [post]
func (h *AuthHandler) SignUp(c *gin.Context) {
	var req domain.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authUsecase.SignUp(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}