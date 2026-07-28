package dtos

import "github.com/google/uuid"

type UserResponse struct {
	UserID   uuid.UUID `json:"userId"`
	Username string    `json:"username"`
}
