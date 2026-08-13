package middlewares

import (
	"errors"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/warmdev17/Wodo-App/pkg/jwt"
	res "github.com/warmdev17/Wodo-App/pkg/response"
)

const UserIDKey string = "userID"

func AuthMiddleware(jwtSvc *jwt.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		log.Println(authHeader)
		if authHeader == "" {
			res.Unauthorized(ctx, "Authorization header is required", "Missing Authorization header")
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			res.Unauthorized(ctx, "Invalid Authorization header format", "Format must be Bearer <token>")
			ctx.Abort()
			return
		}

		claims, err := jwtSvc.ParseToken(parts[1])
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				res.Unauthorized(ctx, "Token has expired", err.Error())
			} else {
				res.Unauthorized(ctx, "Invalid token", err.Error())
			}
			ctx.Abort()
			return
		}

		uid, _ := uuid.Parse(claims.ID)
		ctx.Set(UserIDKey, uid)
		ctx.Next()
	}
}
