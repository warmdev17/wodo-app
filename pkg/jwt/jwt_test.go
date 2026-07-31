package jwt

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/warmdev17/Wodo-App/internal/config"
)

func TestGenerateToken(t *testing.T) {
	cfg := config.Load()
	secretKey := cfg.JWTSecret
	jwtSvc := NewService(secretKey)

	tests := []struct {
		name     string
		userID   uuid.UUID
		duration time.Duration
		wantErr  bool
	}{
		{
			name:     "Success - Generate valid access token (15 mins)",
			userID:   uuid.New(),
			duration: cfg.AccessTokenExpiration,
			wantErr:  false,
		},
		{
			name:     "Success - Generate valid refresh token (30 days)",
			userID:   uuid.New(),
			duration: cfg.RefreshTokenExpiration,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := jwtSvc.GenerateToken(tt.userID, tt.duration)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GenerateToken() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if !tt.wantErr && gotToken == "" {
				t.Errorf("GenerateToken() returned empty token string")
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	cfg := config.Load()
	correctSecret := cfg.JWTSecret
	wrongSecret := "hihihahahuhu"

	jwtSvc := NewService(correctSecret)
	wrongJwtSvc := NewService(wrongSecret)

	testUserID := uuid.New()

	validToken, _ := jwtSvc.GenerateToken(testUserID, 15*time.Minute)
	expiredToken, _ := jwtSvc.GenerateToken(testUserID, -1*time.Minute)
	wrongToken, _ := wrongJwtSvc.GenerateToken(testUserID, 15*time.Minute)

	tests := []struct {
		name        string
		tokenString string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Success - Valid Token",
			tokenString: validToken,
			wantUserID:  testUserID,
			wantErr:     false,
		},
		{
			name:        "Failed - Expired Token",
			tokenString: expiredToken,
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Failed - Token Signed With Wrong Secret",
			tokenString: wrongToken,
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Failed - Malformed Token String",
			tokenString: "invalid.token.structure",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Failed - Empty Token String",
			tokenString: "",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := jwtSvc.ParseToken(tt.tokenString)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if claims.UserID != tt.wantUserID {
					t.Errorf("ParseToken() got UserID = %v, want %v", claims.UserID, tt.wantUserID)
				}
			}
		})
	}
}
