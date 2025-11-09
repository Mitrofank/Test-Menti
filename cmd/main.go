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
		log.Fatal("Переменная среды DATABASE_URL не установлена")
	}

	dbpool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalf("Не удалось создать пул соединений: %v\n", err)
	}

	defer dbpool.Close()

	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("Невозможно выполнить ping базы данных: %v\n", err)
	}
	log.Info("Успешное подключение к базе данных")

	Repo := repository.NewRepository(dbpool)
	Service := service.NewService(Repo)
	Handler := handler.NewHandler(Service)
	router := gin.Default()
	api := router.Group("/api/v1")
	{
		cars := api.Group("/cars")
		{
			cars.POST("/add", Handler.Create)
			cars.GET("/:id", Handler.GetByID)
			cars.GET("", Handler.GetAll)
			cars.DELETE("/:id", Handler.Delete)
		}
	}
	log.Info("Запуск сервера на порту 8080...")

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
