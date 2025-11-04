package handler

import (
	"net/http"

	"github.com/MitrofanK/Test-Menti/internal/models"
	"github.com/MitrofanK/Test-Menti/internal/service"
	"github.com/gin-gonic/gin"
)

type CarHandler struct {
	service service.CarService
}

func NewCarHandler(s service.CarService) *CarHandler {
	return &CarHandler{service: s}
}

func (h *CarHandler) Create(c *gin.Context) {
	var input models.Car
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.service.Create(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
