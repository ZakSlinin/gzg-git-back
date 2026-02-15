package main

import (
	"database/sql"
	"github.com/ZakSlinin/gzg-git-back/internal/auth/handler"
	"github.com/ZakSlinin/gzg-git-back/internal/auth/repository"
	"github.com/ZakSlinin/gzg-git-back/internal/auth/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"os"
	"time"
)

func main() {
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081", "http://localhost:8082"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// DATABASE
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL env variable not set")
	}

	db, err := sql.Open("pgx", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	//JWT
	jwtManager, err := service.NewJWTManager()

	// Auth
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo, jwtManager)
	authHandler := handler.NewAuthHandler(authService)

	apiAuth := r.Group("/api/auth")
	{
		apiAuth.POST("/login", authHandler.Login)
		apiAuth.POST("/register", authHandler.CreateUser)
	}
}
