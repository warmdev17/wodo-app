// Package services
package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/warmdev17/Wodo-App/internal/dtos"
	"github.com/warmdev17/Wodo-App/internal/repositories"
	"github.com/warmdev17/Wodo-App/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo       repositories.Querier
	jwtSvc     *jwt.Service
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthService(repo repositories.Querier, jwtSvc *jwt.Service, accessTTL time.Duration, refreshTTL time.Duration) *AuthService {
	return &AuthService{repo: repo, jwtSvc: jwtSvc, accessTTL: accessTTL, refreshTTL: refreshTTL}
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
		Email:    user.Email,
	}

	return userDTO, nil
}

func (s *AuthService) LoginUser(ctx context.Context, req dtos.LoginRequest) (dtos.AuthResponse, error) {
	user, err := s.repo.GetUserByEmailOrUsername(ctx, req.Identifier)
	if err != nil {
		return dtos.AuthResponse{}, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(req.Password))
	if err != nil {
		return dtos.AuthResponse{}, ErrInvalidCredentials
	}

	accessToken, _, err := s.jwtSvc.GenerateToken(user.ID, s.accessTTL)
	if err != nil {
		return dtos.AuthResponse{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, expiredAt, err := s.jwtSvc.GenerateToken(user.ID, s.refreshTTL)
	if err != nil {
		return dtos.AuthResponse{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	args := repositories.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: expiredAt,
	}
	_, err = s.repo.CreateRefreshToken(ctx, args)
	if err != nil {
		return dtos.AuthResponse{}, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return dtos.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         dtos.ToUserResponse(user),
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req dtos.RefreshTokenRequest) (dtos.RefreshTokenResponse, error) {
	refreshTokenStr := req.RefreshToken
	claims, err := s.jwtSvc.ParseToken(refreshTokenStr)
	if err != nil {
		return dtos.RefreshTokenResponse{}, ErrInvalidToken
	}

	tokenRecord, err := s.repo.GetRefreshTokenByToken(ctx, refreshTokenStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dtos.RefreshTokenResponse{}, ErrTokenNotExists
		}
		return dtos.RefreshTokenResponse{}, ErrInternalServer
	}

	if tokenRecord.IsRevoked {
		return dtos.RefreshTokenResponse{}, ErrTokenRevoked
	}
	if time.Now().After(tokenRecord.ExpiresAt) {
		return dtos.RefreshTokenResponse{}, ErrTokenExpired
	}

	err = s.repo.RevokeToken(ctx, tokenRecord.Token)
	if err != nil {
		return dtos.RefreshTokenResponse{}, ErrInternalServer
	}

	newAccessToken, _, err := s.jwtSvc.GenerateToken(claims.UserID, s.accessTTL)
	if err != nil {
		return dtos.RefreshTokenResponse{}, fmt.Errorf("failed to generate access token: %w", err)
	}
	newRefreshToken, expiredAt, err := s.jwtSvc.GenerateToken(claims.UserID, s.refreshTTL)
	if err != nil {
		return dtos.RefreshTokenResponse{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	args := repositories.CreateRefreshTokenParams{
		UserID:    claims.UserID,
		Token:     newRefreshToken,
		ExpiresAt: expiredAt,
	}

	_, err = s.repo.CreateRefreshToken(ctx, args)
	if err != nil {
		return dtos.RefreshTokenResponse{}, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return dtos.RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
