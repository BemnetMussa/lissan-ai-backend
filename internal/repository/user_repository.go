// internal/repository/user_repository.go
package repository

import "lissanai.com/backend/internal/domain"

// UserRepository is the contract. The Auth Team will build the real MongoDB version.
type UserRepository interface {
	CreateUser(user *domain.User) (*domain.User, error)
}

// For our demo, we create a fake "in-memory" version.
type inMemoryUserRepository struct {}

func NewInMemoryUserRepository() UserRepository {
	return &inMemoryUserRepository{}
}

// The Auth Team will replace this with real MongoDB logic.
func (r *inMemoryUserRepository) CreateUser(user *domain.User) (*domain.User, error) {
	user.ID = "user_123abc" // Fake a database ID
	return user, nil
}