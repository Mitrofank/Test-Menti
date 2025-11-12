package main

import (
	"context"
	"os"

	"github.com/MitrofanK/Test-Menti/internal/handler"
	"github.com/MitrofanK/Test-Menti/internal/repository"
	"github.com/MitrofanK/Test-Menti/internal/service"
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

	repo := repository.NewRepository(dbpool)
	service := service.NewService(repo)
	handler := handler.NewHandler(service, log.New())
	router := gin.Default()
	api := router.Group("/api/v1")
	{
		cars := api.Group("/cars")
		{
			cars.POST("/add", handler.Create)
			cars.GET("/:id", handler.GetByID)
			cars.GET("", handler.GetAll)
			cars.DELETE("/:id", handler.Delete)
		}
	}

	log.Info("Starting the server on port 8080...")

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
