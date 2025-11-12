package service

import (
	"context"
	"fmt"

	"github.com/MitrofanK/Test-Menti/internal/models"
)

type Repository interface {
	Create(ctx context.Context, car models.Car) (int, error)
	GetByID(ctx context.Context, id int) (models.Car, error)
	GetAll(ctx context.Context) ([]models.Car, error)
	Delete(ctx context.Context, id int) error
}

type CarService struct {
	repo Repository
}

func NewService(repo Repository) *CarService {
	return &CarService{
		repo: repo,
	}
}

func (s *CarService) Create(ctx context.Context, car models.Car) (int, error) {
	id, err := s.repo.Create(ctx, car)
	if err != nil {
		return 0, fmt.Errorf("erorr creating car: %w", err)
	}
	return id, nil
}

func (s *CarService) GetByID(ctx context.Context, id int) (models.Car, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CarService) GetAll(ctx context.Context) ([]models.Car, error) {
	return s.repo.GetAll(ctx)
}

func (s *CarService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
