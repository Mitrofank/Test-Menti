package main

import (
	"context"
	"log"
	"os"

	"github.com/MitrofanK/Test-Menti/internal/handler"
	"github.com/MitrofanK/Test-Menti/internal/repository"
	"github.com/MitrofanK/Test-Menti/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}
	dbpool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbpool.Close()
	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}
	log.Println("Successfully connected to the database!")
	carRepo := repository.NewCarPostgres(dbpool)
	carService := service.NewCarService(carRepo)
	carHandler := handler.NewCarHandler(carService)
	router := gin.Default()
	api := router.Group("/api/v1")
	{
		cars := api.Group("/cars")
		{
			cars.POST("/add", carHandler.Create)
			cars.GET("/:id", carHandler.GetByID)
			cars.GET("", carHandler.GetAll)
			cars.DELETE("/:id", carHandler.Delete)
		}
	}
	log.Println("Starting server on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
