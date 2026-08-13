package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/warmdev17/Wodo-App/internal/middlewares"
	"github.com/warmdev17/Wodo-App/internal/services"
)

func GetUserIDFromContext(ctx *gin.Context) (uuid.UUID, error) {
	val, exists := ctx.Get(middlewares.UserIDKey)
	if !exists {
		return uuid.Nil, services.ErrContextNotFound
	}

	userID, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, services.ErrInternalServer
	}

	return userID, nil

}
