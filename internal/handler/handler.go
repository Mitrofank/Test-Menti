package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MitrofanK/Test-Menti/internal/models"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	Create(ctx context.Context, car models.Car) (int, error)
	GetByID(ctx context.Context, id int) (models.Car, error)
	GetAll(ctx context.Context) ([]models.Car, error)
	Delete(ctx context.Context, id int) error
}

type HandlerImpl struct {
	handler Handler
}

func NewHandler(h Handler) *HandlerImpl {
	return &HandlerImpl{handler: h}
}

func (h *HandlerImpl) Create(c *gin.Context) {
	var input models.Car

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if input.Mark == "" || input.Model == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "марка и модель являются обязательными полями"})
		return
	}
	if input.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "цена должна быть больше нуля"})
		return
	}
	if input.Currency != "RUB" && input.Currency != "USD" && input.Currency != "EUR" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "валюта должна быть рубль, доллар США или евро"})
		return
	}

	id, err := h.handler.Create(c.Request.Context(), input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("ошибка создания машины: %w", err)})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *HandlerImpl) Delete(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат идентификатора"})
		return
	}
	err = h.handler.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("ошибка удаления машины %w", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *HandlerImpl) GetByID(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат идентификатора"})
		return
	}
	car, err := h.handler.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("ошибка получения машины %w", err)})
		return
	}
	c.JSON(http.StatusOK, car)
}

func (h *HandlerImpl) GetAll(c *gin.Context) {
	cars, err := h.handler.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("ошибка получения всех машин %w", err)})
		return
	}
	c.JSON(http.StatusOK, cars)
}
