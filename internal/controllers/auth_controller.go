package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/warmdev17/Wodo-App/internal/dtos"
	"github.com/warmdev17/Wodo-App/internal/services"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req dtos.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.authService.RegisterUser(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot register " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Create account success",
		"data":    user,
	})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req dtos.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.authService.LoginUser(ctx, req)
	switch err {
	case services.ErrInvalidCredentials:
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Credentials"})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	ctx.JSON(http.StatusOK, gin.H{"data": user})
}
