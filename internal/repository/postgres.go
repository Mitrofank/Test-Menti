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

func NewPostgresRepository(db *pgxpool.Pool) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (r *Postgres) CreateCar(ctx context.Context, car models.Car) (int, error) {
	query := `insert into cars (make, model, year, owner_id, previous_owners_count, currency, price, options)
	          values ($1, $2, $3, $4, $5, $6, $7, $8) 
	          returning id;`
	row := r.db.QueryRow(ctx, query, car.Make, car.Model, car.Year, car.OwnerID, car.PreviousOwnersCount, car.Currency, car.Price, car.Options)
	var newID int
	err := row.Scan(&newID)

	if err != nil {
		return 0, fmt.Errorf("error creating car: %w", err)
	}

	return newID, nil
}

func (r *Postgres) GetByIDCar(ctx context.Context, id int) (models.Car, error) {
	query := `select id, make, model, year, owner_id, previous_owners_count, currency, price, options
			  from cars 
			  where id = $1;`
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

func (r *Postgres) GetAllCar(ctx context.Context) ([]models.Car, error) {
	query := `select id, make, model, year, owner_id, previous_owners_count, currency, price, options 
			  from cars;`

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

func (r *Postgres) DeleteCar(ctx context.Context, id int) error {
	query := `delete from cars where id = $1;`
	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting car: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("car not found")
	}
	return nil
}

func (r *Postgres) CreateUser(ctx context.Context, user models.User) (int, error) {
	query := `insert into users (email, password_hash, role_id)
			  values ($1, $2, $3)
			  returning id;`

	row := r.db.QueryRow(ctx, query, user.Email, user.PasswordHash, user.RoleID)
	var newID int
	err := row.Scan(&newID)

	if err != nil {
		return 0, fmt.Errorf("error creating user: %w", err)
	}

	return newID, nil
}

func (r *Postgres) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	query := `select id, email, password_hash, role_id 
			  from users
			  where email = $1;`
	row := r.db.QueryRow(ctx, query, email)
	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.RoleID,
	)

	if err != nil {
		return models.User{}, fmt.Errorf("error getting user^ %w", err)
	}

	return user, nil
}
