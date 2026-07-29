package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/warmdev17/Wodo-App/internal/dtos"
	"github.com/warmdev17/Wodo-App/internal/services"
	res "github.com/warmdev17/Wodo-App/pkg/response"
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
		res.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	user, err := c.authService.RegisterUser(ctx.Request.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrEmailTaken) {
			res.Conflict(ctx, "Cannot register", services.ErrEmailTaken.Error())
			return
		} else if errors.Is(err, services.ErrUsernameTaken) {
			res.Conflict(ctx, "Cannot register", services.ErrUsernameTaken.Error())
			return
		}
		res.InternalError(ctx, "Cannot register", err.Error())
		return
	}

	res.Created(ctx, "Create new account successful", user)
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req dtos.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	user, err := c.authService.LoginUser(ctx.Request.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			res.Unauthorized(ctx, "Failed to login", services.ErrInvalidCredentials.Error())
			return
		}

		res.InternalError(ctx, "Error from server side", err.Error())
		return
	}
	res.Success(ctx, "Login successful", user)
}
