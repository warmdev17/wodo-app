package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/warmdev17/Wodo-App/internal/config"
	"github.com/warmdev17/Wodo-App/internal/controllers"
	"github.com/warmdev17/Wodo-App/internal/middlewares"
	"github.com/warmdev17/Wodo-App/internal/repositories"
	"github.com/warmdev17/Wodo-App/internal/services"
	"github.com/warmdev17/Wodo-App/pkg/jwt"
	"github.com/warmdev17/Wodo-App/pkg/response"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := config.Load()

	connStr := cfg.DatabaseURL
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer func() {
		err := db.Close()

		if err != nil {
			log.Println("Failed to close connection pool")
		}
	}()

	repo := repositories.New(db)
	jwtSvc := jwt.NewService(cfg.JWTSecret)
	authSvc := services.NewAuthService(repo, jwtSvc, cfg.AccessTokenExpiration, cfg.RefreshTokenExpiration)
	authCtrl := controllers.NewAuthController(authSvc, cfg.IsProduction)

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", authCtrl.Register)
		api.POST("/auth/login", authCtrl.Login)
		api.POST("/auth/refresh", authCtrl.Refresh)
	}

	protected := r.Group("/api/v1")
	protected.Use(middlewares.AuthMiddleware(jwtSvc))
	{
		protected.GET("/users/me", func(ctx *gin.Context) {
			userID, _ := ctx.Get("userID")
			response.Success(ctx, "Get profile successful", gin.H{
				"userId": userID,
			})
		})
		protected.POST("/auth/logout", authCtrl.Logout)
	}

	err = r.Run(":8080")

	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
