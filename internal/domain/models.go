// internal/domain/models.go
package domain

type SignUpRequest struct {
	Username string `json:"username" binding:"required" example:"testuser"`
	Email    string `json:"email" binding:"required,email" example:"test@lissanai.com"`
	Password string `json:"password" binding:"required,min=8" example:"strongpassword123"`
}

type User struct {
	ID       string `json:"id" example:"user_123abc"`
	Username string `json:"username" example:"testuser"`
	Email    string `json:"email" example:"test@lissanai.com"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}