package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/warmdev17/Wodo-App/internal/config"
	"github.com/warmdev17/Wodo-App/internal/controllers"
	"github.com/warmdev17/Wodo-App/internal/middlewares"
	"github.com/warmdev17/Wodo-App/internal/repositories"
	"github.com/warmdev17/Wodo-App/internal/services"
	"github.com/warmdev17/Wodo-App/pkg/jwt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	connStr := cfg.DatabaseURL
	db, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	db.Close()

	repo := repositories.New(db)
	jwtSvc := jwt.NewService(cfg.JWTSecret)
	authSvc := services.NewAuthService(repo, jwtSvc, cfg.AccessTokenExpiration, cfg.RefreshTokenExpiration)
	authCtrl := controllers.NewAuthController(authSvc, cfg.IsProduction)
	userCtrl := controllers.NewUserController()

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
		protected.GET("/users/me", userCtrl.GetMe)
		protected.POST("/auth/logout", authCtrl.Logout)
	}

	err = r.Run(":8080")

	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
