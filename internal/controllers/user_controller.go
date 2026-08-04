package controllers

import (
	"github.com/gin-gonic/gin"
	res "github.com/warmdev17/Wodo-App/pkg/response"
)

type UserController struct {
}

func NewUserController() *UserController {
	return &UserController{}
}

func (c *UserController) GetMe(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		res.Unauthorized(ctx, "Unauthorized", "User context not found")
		return
	}

	res.Success(ctx, "Get profile successful", gin.H{
		"userId": userID,
	})
}
