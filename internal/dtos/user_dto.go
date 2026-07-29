package dtos

import (
	"github.com/google/uuid"
	"github.com/warmdev17/Wodo-App/internal/repositories"
)

type UserResponse struct {
	UserID   uuid.UUID `json:"userId"`
	Username string    `json:"username"`
}

func ToUserResponse(u repositories.User) UserResponse {
	return UserResponse{
		UserID:   u.ID,
		Username: u.Username,
	}
}
