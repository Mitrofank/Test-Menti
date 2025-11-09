package service

import (
	"context"
	"fmt"

	"github.com/MitrofanK/Test-Menti/internal/models"
	"github.com/MitrofanK/Test-Menti/internal/repository"
)

type Service interface {
	Create(ctx context.Context, car models.Car) (int, error)
	GetByID(ctx context.Context, id int) (models.Car, error)
	GetAll(ctx context.Context) ([]models.Car, error)
	Delete(ctx context.Context, id int) error
}

type CarServiceImpl struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *CarServiceImpl {
	return &CarServiceImpl{
		repo: repo,
	}
}

func (s *CarServiceImpl) Create(ctx context.Context, car models.Car) (int, error) {
	id, err := s.repo.Create(ctx, car)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания машины: %w", err)
	}
	return id, nil
}

func (s *CarServiceImpl) GetByID(ctx context.Context, id int) (models.Car, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CarServiceImpl) GetAll(ctx context.Context) ([]models.Car, error) {
	return s.repo.GetAll(ctx)
}

func (s *CarServiceImpl) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
