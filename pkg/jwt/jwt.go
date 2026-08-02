package jwt

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwtlib.RegisteredClaims
	UserID uuid.UUID `json:"userId"`
}

type Service struct {
	secret []byte
}

var ErrTokenExpired = jwtlib.ErrTokenExpired

func NewService(secret string) *Service {
	return &Service{secret: []byte(secret)}
}

func (s *Service) GenerateToken(userID uuid.UUID, duration time.Duration) (string, time.Time, error) {
	now := time.Now()
	expiredAt := time.Now().Add(duration)
	claims := Claims{
		RegisteredClaims: jwtlib.RegisteredClaims{
			IssuedAt:  jwtlib.NewNumericDate(now),
			ExpiresAt: jwtlib.NewNumericDate(expiredAt),
			Subject:   userID.String(),
		},
		UserID: userID,
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, expiredAt, nil
}

func (s *Service) ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwtlib.ParseWithClaims(tokenString, claims, func(t *jwtlib.Token) (any, error) {
		if _, ok := t.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, jwtlib.ErrSignatureInvalid
		}

		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwtlib.ErrSignatureInvalid
	}

	return claims, nil
}
