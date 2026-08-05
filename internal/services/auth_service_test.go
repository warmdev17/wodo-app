package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/warmdev17/Wodo-App/internal/config"
	"github.com/warmdev17/Wodo-App/internal/dtos"
	"github.com/warmdev17/Wodo-App/internal/repositories"
	"github.com/warmdev17/Wodo-App/pkg/jwt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Identifier struct {
	Username string
	Email    string
	Password string
}

func generateIdentifier(username string, email string) Identifier {
	timestamp := time.Now().UnixNano()

	if username == "" {
		username = fmt.Sprintf("user_%v", timestamp)
	}
	if email == "" {
		email = fmt.Sprintf("user_%v@test.com", timestamp)
	}
	return Identifier{
		Username: username,
		Email:    email,
		Password: fmt.Sprintf("password_%v", timestamp),
	}
}

func TestRegisterUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cfg := config.Load()

	connStr := "postgres://postgres:maiphuong@127.0.0.1:5432/wodo_test_db?sslmode=disable"
	testDB, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer testDB.Close()

	repo := repositories.New(testDB)
	jwtSvc := jwt.NewService(cfg.JWTSecret)

	authSvc := NewAuthService(repo, jwtSvc, cfg.AccessTokenExpiration, cfg.RefreshTokenExpiration)

	base := generateIdentifier("", "")
	dupUser := generateIdentifier("warmdev", "")
	dupEmail := generateIdentifier("", "warmdevofficial@gmail.com")

	reqSuccess := dtos.RegisterRequest{
		Username: base.Username,
		Email:    base.Email,
		Password: base.Password,
	}

	reqDupUser := dtos.RegisterRequest{
		Username: dupUser.Username,
		Email:    dupUser.Email,
		Password: dupUser.Password,
	}

	reqDupEmail := dtos.RegisterRequest{
		Username: dupEmail.Username,
		Email:    dupEmail.Email,
		Password: dupEmail.Password,
	}

	tests := []struct {
		name         string
		req          dtos.RegisterRequest
		wantUsername string
		wantEmail    string
		wantErr      error
	}{
		{
			name:         "Success - Register",
			req:          reqSuccess,
			wantUsername: reqSuccess.Username,
			wantEmail:    reqSuccess.Email,
			wantErr:      nil,
		},
		{
			name:         "Failed - Username already exists",
			req:          reqDupUser,
			wantUsername: reqDupUser.Username,
			wantEmail:    reqDupUser.Email,
			wantErr:      ErrUsernameTaken,
		},
		{
			name:         "Failed - Email already exists",
			req:          reqDupEmail,
			wantUsername: reqDupEmail.Username,
			wantEmail:    reqDupEmail.Email,
			wantErr:      ErrEmailTaken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := authSvc.RegisterUser(ctx, tt.req)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("RegisterUser() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if tt.wantErr != nil {
				return
			}

			if user.UserID == uuid.Nil {
				t.Errorf("RegisterUser() got empty UserID")
			}

			if user.Username != tt.wantUsername {
				t.Errorf("got username = %v, want = %v", user.Username, tt.wantUsername)
			}

			if user.Email != tt.wantEmail {
				t.Errorf("got email = %v, want = %v", user.Email, tt.wantEmail)
			}
		})
	}
}

func TestLoginUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cfg := config.Load()

	connStr := "postgres://postgres:maiphuong@127.0.0.1:5432/wodo_test_db?sslmode=disable"
	testDB, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer testDB.Close()

	repo := repositories.New(testDB)
	jwtSvc := jwt.NewService(cfg.JWTSecret)
	authSvc := NewAuthService(repo, jwtSvc, cfg.AccessTokenExpiration, cfg.RefreshTokenExpiration)

	reqSucess := dtos.LoginRequest{
		Identifier: "warmdev",
		Password:   "maiphuong",
	}
	reqWrongIdentifier := dtos.LoginRequest{
		Identifier: "nhomnhom",
		Password:   "maiphuong",
	}
	reqWrongPassword := dtos.LoginRequest{
		Identifier: "warmdev",
		Password:   "nhomnhom",
	}

	tests := []struct {
		name    string
		req     dtos.LoginRequest
		wantErr error
	}{
		{
			name:    "Success - Login",
			req:     reqSucess,
			wantErr: nil,
		},
		{
			name:    "Failed - Wrong Identifier",
			req:     reqWrongIdentifier,
			wantErr: ErrInvalidCredentials,
		},
		{
			name:    "Failed - Wrong Password",
			req:     reqWrongPassword,
			wantErr: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := authSvc.LoginUser(ctx, tt.req)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("LoginUser() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if tt.wantErr != nil {
				return
			}

			if res.AccessToken == "" {
				t.Errorf("LoginUser() got empty AccessToken")
			}

			if res.RefreshToken == "" {
				t.Errorf("LoginUser() got empty RefreshToken")
			}
		})
	}
}
