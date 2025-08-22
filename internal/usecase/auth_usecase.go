// internal/usecase/auth_usecase.go
package usecase

import (
	"lissanai.com/backend/internal/domain"
	"lissanai.com/backend/internal/repository"
)

type AuthUsecase interface {
	SignUp(req *domain.SignUpRequest) (*domain.User, error)
}

type authUsecase struct {
	userRepo repository.UserRepository
}

func NewAuthUsecase(userRepo repository.UserRepository) AuthUsecase {
	return &authUsecase{userRepo: userRepo}
}

func (u *authUsecase) SignUp(req *domain.SignUpRequest) (*domain.User, error) {
	// The Auth Team will add password hashing and validation here.
	
	newUser := &domain.User{
		Username: req.Username,
		Email:    req.Email,
	}

	return u.userRepo.CreateUser(newUser)
}