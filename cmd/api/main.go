package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/warmdev17/Wodo-App/internal/config"
	"github.com/warmdev17/Wodo-App/internal/controllers"
	"github.com/warmdev17/Wodo-App/internal/repositories"
	"github.com/warmdev17/Wodo-App/internal/services"
	"github.com/warmdev17/Wodo-App/pkg/jwt"

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
	authSvc := services.NewAuthService(repo, jwtSvc)
	authCtrl := controllers.NewAuthController(authSvc)

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", authCtrl.Register)
		api.POST("/auth/login", authCtrl.Login)
	}

	err = r.Run(":8080")

	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
