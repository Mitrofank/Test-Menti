package car

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MitrofanK/Test-Menti/internal/errorsx"
	"github.com/MitrofanK/Test-Menti/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

type FacadeService interface {
	GetCarWithConversion(ctx context.Context, id int, targetCurrency string) (models.Car, error)
	CreateCar(ctx context.Context, car models.Car) (int, error)
	GetAllCar(ctx context.Context) ([]models.Car, error)
	DeleteCar(ctx context.Context, id int) error
}

type Handler struct {
	facade FacadeService
	log    *log.Logger
}

func NewHandler(facade FacadeService, log *log.Logger) *Handler {
	return &Handler{
		facade: facade,
		log:    log,
	}
}

func (h *Handler) Create(c *gin.Context) {
	var input models.Car

	if err := c.ShouldBindJSON(&input); err != nil {
		h.log.Warning(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("read error from request body: %w", err)})
		return
	}

	if input.Make == "" || input.Model == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "make and model are required fields"})
		return
	}

	if input.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "the price must be greater than zero"})
		return
	}
	// добавить кастомные типы
	if input.Currency != "RUB" && input.Currency != "USD" && input.Currency != "EUR" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The currency must be the ruble, US dollar, or euro."})
		return
	}

	id, err := h.facade.CreateCar(c.Request.Context(), input)
	if err != nil {
		h.log.Error("create car error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("create car error: %w", err)})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) DeleteCar(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid identifier format"})
		return
	}

	if err := h.facade.DeleteCar(c.Request.Context(), id); err != nil {
		h.log.Error("error deleting car")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("error deleting car: %w", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) GetByIDCar(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid identifier format"})
		return
	}

	targetCurrency := c.Query("currency")

	car, err := h.facade.GetCarWithConversion(c.Request.Context(), id, targetCurrency)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "car not found"})
			return
		}

		if errors.Is(err, errorsx.ErrCurrencyNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		h.log.Errorf("error getting car with conversion: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, car)
}

func (h *Handler) GetAllCar(c *gin.Context) {
	cars, err := h.facade.GetAllCar(c.Request.Context())
	if err != nil {
		h.log.Error("error receiving all id")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("error getting all vehicles: %w", err)})
		return
	}
	c.JSON(http.StatusOK, cars)
}
