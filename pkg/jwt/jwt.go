package jwt

import (
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

func (s *Service) GenerateToken(userID uuid.UUID, ttlMinutes int) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwtlib.RegisteredClaims{
			IssuedAt:  jwtlib.NewNumericDate(now),
			ExpiresAt: jwtlib.NewNumericDate(now.Add(time.Duration(ttlMinutes) * time.Minute)),
			Subject:   userID.String(),
		},
		UserID: userID,
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)

	return token.SignedString(s.secret)
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
