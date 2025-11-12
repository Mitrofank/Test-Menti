package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MitrofanK/Test-Menti/internal/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type CarService interface {
	Create(ctx context.Context, car models.Car) (int, error)
	GetByID(ctx context.Context, id int) (models.Car, error)
	GetAll(ctx context.Context) ([]models.Car, error)
	Delete(ctx context.Context, id int) error
}

type Handler struct {
	carService CarService
	log        *log.Logger
}

func NewHandler(carService CarService, log *log.Logger) *Handler {
	return &Handler{
		carService: carService,
		log:        log,
	}
}

func (h *Handler) Create(c *gin.Context) {
	var input models.Car

	if err := c.ShouldBindJSON(&input); err != nil {
		h.log.Error("read error from request body")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("read error from request body: %w", err)})
		return
	}

	if input.Mark == "" || input.Model == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "make and model are required fields"})
		return
	}

	if input.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "the price must be greater than zero"})
		return
	}

	if input.Currency != "RUB" && input.Currency != "USD" && input.Currency != "EUR" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The currency must be the ruble, US dollar, or euro."})
		return
	}

	id, err := h.carService.Create(c.Request.Context(), input)
	if err != nil {
		h.log.Error("create car error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("create car error: %w", err)})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) Delete(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid identifier format"})
		return
	}

	if err := h.carService.Delete(c.Request.Context(), id); err != nil {
		h.log.Error("error deleting car")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("error deleting car: %w", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) GetByID(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid identifier format"})
		return
	}

	car, err := h.carService.GetByID(c.Request.Context(), id)
	if err != nil {
		h.log.Error("error receiving id")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("error getting vehicle id: %w", err)})
		return
	}
	c.JSON(http.StatusOK, car)
}

func (h *Handler) GetAll(c *gin.Context) {
	cars, err := h.carService.GetAll(c.Request.Context())
	if err != nil {
		h.log.Error("error receiving all id")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("error getting all vehicles: %w", err)})
		return
	}
	c.JSON(http.StatusOK, cars)
}
