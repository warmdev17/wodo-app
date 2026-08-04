package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/warmdev17/Wodo-App/internal/dtos"
	"github.com/warmdev17/Wodo-App/internal/services"
	res "github.com/warmdev17/Wodo-App/pkg/response"
)

type AuthController struct {
	authService  *services.AuthService
	isProduction bool
}

func NewAuthController(authService *services.AuthService, isProduction bool) *AuthController {
	return &AuthController{authService: authService, isProduction: isProduction}
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

func (c *AuthController) SetRefreshTokenCookie(ctx *gin.Context, token string) {
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("refresh_token", token, int(c.authService.RefreshTTL().Seconds()), "/api/v1/auth", "", c.isProduction, true)
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
	c.SetRefreshTokenCookie(ctx, user.RefreshToken)
	res.Success(ctx, "Login successful", user)
}

func (c *AuthController) Refresh(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		res.Unauthorized(ctx, "Missing refresh token cookie", err.Error())
		return
	}

	req := dtos.RefreshTokenRequest{RefreshToken: refreshToken}
	result, err := c.authService.RefreshToken(ctx.Request.Context(), req)

	if err != nil {
		if errors.Is(err, services.ErrTokenExpired) || errors.Is(err, services.ErrInvalidToken) {
			res.Unauthorized(ctx, "Invalid or expired refresh token", err.Error())
			return
		}
		res.InternalError(ctx, "Error from server side", err.Error())
		return
	}

	c.SetRefreshTokenCookie(ctx, result.RefreshToken)
	res.Success(ctx, "Refresh new token successful", result)
}

func (c *AuthController) Logout(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		res.Unauthorized(ctx, "Missing refresh token cookie", err.Error())
		return
	}

	err = c.authService.LogoutUser(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, services.ErrInternalServer) {
			res.InternalError(ctx, "Error from server side", err.Error())
			return
		}
		res.Unauthorized(ctx, "Token has invalid or expired", err.Error())
		return
	}

	ctx.SetCookie("refresh_token", "", -1, "/api/v1/auth", "", c.isProduction, true)
	res.SuccessNoData(ctx, "Logout successful")
}
