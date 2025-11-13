package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	errorsx "github.com/MitrofanK/Test-Menti/internal/errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type UserService interface {
	SignUp(ctx context.Context, email, password string) (int, error)
	SignIn(ctx context.Context, email, password string) (string, error)
}

type Handler struct {
	userService UserService
	log         *log.Logger
}

func NewHandler(userService UserService, log *log.Logger) *Handler {
	return &Handler{
		userService: userService,
		log:         log,
	}
}

type signUpInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) SignUp(c *gin.Context) {
	var input signUpInput

	if err := c.ShouldBindJSON(&input); err != nil {
		h.log.Error("read error from request body")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("read error from request body: %w", err)})
		return
	}

	id, err := h.userService.SignUp(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		if errors.Is(err, errorsx.ErrUserExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, errorsx.ErrPasBeEmpty) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, errorsx.ErrPasLength) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, errorsx.ErrPasAndLoginSame) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		h.log.Error("create user error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) SignIn(c *gin.Context) {
	var input signUpInput

	if err := c.ShouldBindJSON(&input); err != nil {
		h.log.Error("read error from request body")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("read error from request body: %w", err)})
		return
	}

	token, err := h.userService.SignIn(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		if errors.Is(err, errorsx.ErrIncorLoginOrPas) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		h.log.Error("login error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
