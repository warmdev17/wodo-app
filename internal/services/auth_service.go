// Package services
package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/warmdev17/Wodo-App/internal/dtos"
	"github.com/warmdev17/Wodo-App/internal/repositories"
	"github.com/warmdev17/Wodo-App/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailTaken         = errors.New("email already registered")
	ErrUsernameTaken      = errors.New("username already registered")
)

type AuthService struct {
	repo   repositories.Querier
	jwtSvc *jwt.Service
}

func NewAuthService(repo repositories.Querier, jwtSvc *jwt.Service) *AuthService {
	return &AuthService{repo: repo, jwtSvc: jwtSvc}
}

func (s *AuthService) RegisterUser(ctx context.Context, req dtos.RegisterRequest) (dtos.UserResponse, error) {
	var user repositories.User
	var err error

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dtos.UserResponse{}, ErrInvalidInput
	}

	args := repositories.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		HashPassword: string(hashedPassword),
	}

	user, err = s.repo.CreateUser(ctx, args)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_username_key":
				return dtos.UserResponse{}, ErrUsernameTaken
			case "users_email_key":
				return dtos.UserResponse{}, ErrEmailTaken
			}
		}
	}

	userDTO := dtos.UserResponse{
		UserID:   user.ID,
		Username: user.Username,
	}

	return userDTO, nil
}
