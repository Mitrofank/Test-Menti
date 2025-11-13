package main

import (
	"context"
	"os"
	"time"

	carhandler "github.com/MitrofanK/Test-Menti/internal/api/handler/car"
	userhandler "github.com/MitrofanK/Test-Menti/internal/api/handler/user"
	"github.com/MitrofanK/Test-Menti/internal/api/middleware"
	"github.com/MitrofanK/Test-Menti/internal/auth"
	"github.com/MitrofanK/Test-Menti/internal/repository"
	carservice "github.com/MitrofanK/Test-Menti/internal/service/car"
	userservice "github.com/MitrofanK/Test-Menti/internal/service/user"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

func main() {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatal("Environment variable DATABASE_URL is not set")
	}

	dbpool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalf("Not able to connect to database: %v\n", err)
	}

	defer dbpool.Close()

	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("Not able to ping database: %v\n", err)
	}

	log.Info("Successful connection to the database")

	signingKey := os.Getenv("SIGNING_KEY")
	if signingKey == "" {
		log.Fatal("SIGNING_KEY environment variable is not set")
	}

	authMiddleware := middleware.NewMiddleware(signingKey)

	tokenTTL := time.Hour * 12
	tokenManager, err := auth.NewTokenManager(signingKey, tokenTTL)

	if err != nil {
		log.Fatalf("Error creating token manager: %v\n", err)
	}

	repo := repository.NewRepository(dbpool)

	userService := userservice.NewService(repo, tokenManager)
	carService := carservice.NewService(repo)

	carHandler := carhandler.NewHandler(carService, log.New())
	userHandler := userhandler.NewHandler(userService, log.New())

	router := gin.Default()

	api := router.Group("/api/v1")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/sign-up", userHandler.SignUp)
			authGroup.POST("/sign-in", userHandler.SignIn)
		}

		protected := api.Group("/")
		protected.Use(authMiddleware.UserIdentity)
		{
			cars := protected.Group("/cars")
			{
				cars.POST("/add", carHandler.Create)
				cars.GET("/:id", carHandler.GetByID)
				cars.GET("", carHandler.GetAll)
				cars.DELETE("/:id", carHandler.Delete)
			}
		}
	}

	log.Info("Starting the server on port 8080...")

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
