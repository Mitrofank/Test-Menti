package main

import (
	"context"

	carhandler "github.com/MitrofanK/Test-Menti/internal/api/handler/car"
	userhandler "github.com/MitrofanK/Test-Menti/internal/api/handler/user"
	"github.com/MitrofanK/Test-Menti/internal/api/middleware"
	"github.com/MitrofanK/Test-Menti/internal/auth"
	"github.com/MitrofanK/Test-Menti/internal/config"
	"github.com/MitrofanK/Test-Menti/internal/repository"
	carservice "github.com/MitrofanK/Test-Menti/internal/service/car"
	userservice "github.com/MitrofanK/Test-Menti/internal/service/user"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	dbpool, err := pgxpool.New(context.Background(), cfg.Postgres.URL)
	if err != nil {
		log.Fatalf("Not able to connect to database: %v\n", err)
	}

	defer dbpool.Close()

	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("Not able to ping database: %v\n", err)
	}

	log.Info("Successful connection to the database")

	authMiddleware := middleware.NewMiddleware(cfg.JWT.SigningKey)

	tokenManager, err := auth.NewTokenManager(cfg.JWT.SigningKey, cfg.JWT.TokenTTL)

	if err != nil {
		log.Fatalf("Error creating token manager: %v\n", err)
	}

	repo := repository.NewPostgresRepository(dbpool)

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

		publicCars := api.Group("/cars")
		{
			publicCars.GET("", carHandler.GetAll)
			publicCars.GET("/:id", carHandler.GetByID)
		}

		privateCars := api.Group("/cars")
		privateCars.Use(authMiddleware.UserIdentity)
		{
			privateCars.POST("/add", carHandler.Create)
			privateCars.DELETE("/:id", carHandler.Delete)
		}
	}

	log.Info("Starting the server on port 8080...")

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
