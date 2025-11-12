package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/MitrofanK/Test-Menti/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (r *Postgres) Create(ctx context.Context, car models.Car) (int, error) {
	query := `insert into cars (id, make, model, year, OwnerID, PreviousOwnersCount, currency, price, options) 
	          values ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
	          returning id`
	row := r.db.QueryRow(ctx, query, car.Make, car.Model, car.Year, car.OwnerID, car.PreviousOwnersCount, car.Currency, car.Price, car.Options)
	var newID int
	err := row.Scan(&newID)

	if err != nil {
		return 0, fmt.Errorf("error creating car: %w", err)
	}

	return newID, nil
}
func (r *Postgres) GetByID(ctx context.Context, id int) (models.Car, error) {
	query := `select id, make, model, year, OwnerID, PreviousOwnersCount, currency, price, options 
			  from cars 
			  where id = $1`
	row := r.db.QueryRow(ctx, query, id)
	var car models.Car
	err := row.Scan(
		&car.ID,
		&car.Make,
		&car.Model,
		&car.Year,
		&car.OwnerID,
		&car.PreviousOwnersCount,
		&car.Currency,
		&car.Price,
		&car.Options,
	)

	if err != nil {
		return models.Car{}, fmt.Errorf("error getting car: %w", err)
	}

	return car, nil
}

func (r *Postgres) GetAll(ctx context.Context) ([]models.Car, error) {
	query := `select id, make, model, year, OwnerID, PreviousOwnersCount, currency, price, options 
			  from cars`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error getting all cars: %w", err)
	}
	defer rows.Close()

	cars := make([]models.Car, 0)

	for rows.Next() {
		var car models.Car
		if err := rows.Scan(
			&car.ID,
			&car.Make,
			&car.Model,
			&car.Year,
			&car.OwnerID,
			&car.PreviousOwnersCount,
			&car.Currency,
			&car.Price,
			&car.Options,
		); err != nil {
			return nil, fmt.Errorf("error scanning car: %w", err)
		}
		cars = append(cars, car)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cars, nil
}

func (r *Postgres) Delete(ctx context.Context, id int) error {
	query := `delete from cars where id = $1`
	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting car: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("car not found")
	}
	return nil
}
