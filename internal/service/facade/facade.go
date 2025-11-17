package facade

import (
	"context"
	"strings"

	"github.com/MitrofanK/Test-Menti/internal/models"
)

type CarService interface {
	CreateCar(ctx context.Context, car models.Car) (int, error)
	GetByIDCar(ctx context.Context, id int) (models.Car, error)
	GetAllCar(ctx context.Context) ([]models.Car, error)
	DeleteCar(ctx context.Context, id int) error
}

type CurrencyService interface {
	Convert(ctx context.Context, from, to string, amount float64) (float64, error)
}

type Service struct {
	carService      CarService
	currencyService CurrencyService
}

func NewService(cs CarService, curS CurrencyService) *Service {
	return &Service{
		carService:      cs,
		currencyService: curS,
	}
}

func (s *Service) GetCarWithConversion(ctx context.Context, id int, targetCurrency string) (models.Car, error) {
	car, err := s.carService.GetByIDCar(ctx, id)
	if err != nil {
		return models.Car{}, err
	}

	if targetCurrency != "" && models.CurrencyCode(strings.ToUpper(targetCurrency)) != car.Currency {
		convertedPrice, err := s.currencyService.Convert(
			ctx,
			string(car.Currency),
			targetCurrency,
			float64(car.Price),
		)
		if err != nil {
			return models.Car{}, err
		}

		car.Price = int(convertedPrice)
		car.Currency = models.CurrencyCode(strings.ToUpper(targetCurrency))
	}

	return car, nil
}

func (s *Service) CreateCar(ctx context.Context, car models.Car) (int, error) {
	return s.carService.CreateCar(ctx, car)
}

func (s *Service) GetAllCar(ctx context.Context) ([]models.Car, error) {
	return s.carService.GetAllCar(ctx)
}

func (s *Service) DeleteCar(ctx context.Context, id int) error {
	return s.carService.DeleteCar(ctx, id)
}
