package service

import (
	"context"

	"github.com/MitrofanK/Test-Menti/internal/models"
)

type CarService interface {
	Create(ctx context.Context, car models.Car) (int, error) 
	GetByID(ctx context.Context, id int) (models.Car, error)
	GetAll(ctx context.Context) ([]models.Car, error)
	Delete(ctx context.Context, id int) error
}

type CarServiceImpl struct {
	repo repository.CarRepository 
}

func NewCarService(repo repository.CarRepository) *CarServiceImpl {
	return &CarServiceImpl{
		repo: repo
	}
}

func (s *CarServiceImpl) Create(ctx context.Context, car models.Car) (int, error) {
	if car.Mark == "" || car.Model == "" {
		return 0, errors.New("mark and model are required fields")
	}
	if car.Price <= 0 {
		return 0, errors.New("price must be greater than zero")
	}
	if car.Currency != "RUB" && car.Currency != "USD" && car.Currency != "EUR" {
		return 0, errors.New("currency must be RUB, USD or EUR")
	}
	id, err := s.repo.Create(ctx, car)
	if err != nil {
		return 0, err
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