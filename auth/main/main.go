package main

import (
	"github.com/ZakSlinin/gzg-git-back/auth/handler"
	"github.com/ZakSlinin/gzg-git-back/auth/repository"
	"github.com/ZakSlinin/gzg-git-back/auth/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // driver
	"log"
	"os"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	jwtManager, err := service.NewJWTManager()
	if err != nil {
		log.Fatalf("failed to create jwt manager: %v", err)
	}
	fileService := service.NewFileService("./uploads/avatars")

	authRepo := repository.NewAuthRepository(db.DB)
	authService := service.NewAuthService(authRepo, jwtManager)
	authHandler := handler.NewAuthHandler(authService, fileService)

	router := gin.Default()
	router.Static("/uploads", "./uploads/avatar")

	api := router.Group("/api/v1")
	{
		api.POST("api/auth/register", authHandler.CreateUser)
		api.POST("api/auth/login", authHandler.Login)
	}

	router.Run(":" + port)
	log.Println("Server is running on port " + port)
}
