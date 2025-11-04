package repository

import (
	"context"

	"github.com/MitrofanK/Test-Menti/internal/models"
)

type CarRepository interface {
	Create(ctx context.Context, car models.Car) (int, error)
	GetByID(ctx context.Context, id int) (models.Car, error)
	GetAll(ctx context.Context) ([]models.Car, error)
	Delete(ctx context.Context, id int) error
}

type CarPostgres struct {
	db *pgxpool.Pool
}

func NewCarPostgres(db *pgxpool.Pool) *CarRepository {
	return &CarPostgres{
		db: db
	}
}

func (r *CarPostgres) Create(ctx context.Context, car models.Car) (int, error) {
	query := `INSERT INTO cars (mark, model, owner_count, price, currency, options) 
	          VALUES ($1, $2, $3, $4, $5, $6) 
	          RETURNING id`
	row := r.db.QueryRow(ctx, query, car.Mark, car.Model, car.OwnerCount, car.Price, car.Currency, car.Options)
	var newID int
	err := row.Scan(&newID)
	if err != nil {
		return 0, err 
	}
	return newID, nil
}

func (r *CarPostgres) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM cars WHERE id = $1`
	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("car not found")
	}
	return nil
}