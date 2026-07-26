// Package services
package services

import (
	"context"

	"github.com/warmdev17/Wodo-App/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repositories.Querier
}

func NewUserService(repo repositories.Querier) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(ctx context.Context, username string, email string, password string) (repositories.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return repositories.User{}, err
	}

	args := repositories.CreateUserParams{
		Username:     username,
		Email:        email,
		HashPassword: string(hashedPassword),
	}

	user, err := s.repo.CreateUser(ctx, args)
	if err != nil {
		return repositories.User{}, nil
	}

	return user, nil
}
