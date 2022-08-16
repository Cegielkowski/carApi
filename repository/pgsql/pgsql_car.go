package pgsql

import (
	"context"
	"database/sql"
	"fmt"

	"carApi/entity"
)

// CarRepository represent the car's repository contract
type CarRepository interface {
	Create(ctx context.Context, car *entity.Car) error
	GetByID(ctx context.Context, id int64) (entity.Car, error)
	Fetch(ctx context.Context) ([]entity.Car, error)
	Update(ctx context.Context, car *entity.Car) error
	Delete(ctx context.Context, id int64) error
}

type pgsqlCarRepository struct {
	db *sql.DB
}

// NewPgsqlCarRepository NewCarRepository will create new an carRepository object representation of CarRepository interface
func NewPgsqlCarRepository(db *sql.DB) CarRepository {
	return &pgsqlCarRepository{
		db: db,
	}
}

func (r *pgsqlCarRepository) Create(ctx context.Context, car *entity.Car) (err error) {
	query := "INSERT INTO cars (make, model, package, color, mileage, price, category, year, identification, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"
	_, err = r.db.ExecContext(ctx, query, car.Make, car.Model, car.Package, car.Color, car.Mileage, car.Price, car.Category, car.Year, car.Identification, car.CreatedAt, car.UpdatedAt)
	return
}

func (r *pgsqlCarRepository) GetByID(ctx context.Context, id int64) (car entity.Car, err error) {
	query := "SELECT id, make, model, package, color, mileage, price, category, year, identification, created_at, updated_at FROM cars WHERE id = $1"
	err = r.db.QueryRowContext(ctx, query, id).Scan(&car.ID, &car.Make, &car.Model, &car.Package, &car.Color, &car.Mileage, &car.Price, &car.Category, &car.Year, &car.Identification, &car.CreatedAt, &car.UpdatedAt)
	return
}

func (r *pgsqlCarRepository) Fetch(ctx context.Context) (cars []entity.Car, err error) {
	query := "SELECT id, make, model, package, color, mileage, price, category, year, identification, created_at, updated_at FROM cars"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return cars, err
	}

	defer rows.Close()

	for rows.Next() {
		var car entity.Car
		err := rows.Scan(&car.ID, &car.Make, &car.Model, &car.Package, &car.Color, &car.Mileage, &car.Price, &car.Category, &car.Year, &car.Identification, &car.CreatedAt, &car.UpdatedAt)
		if err != nil {
			return cars, err
		}

		cars = append(cars, car)
	}

	return cars, nil
}

func (r *pgsqlCarRepository) Update(ctx context.Context, car *entity.Car) (err error) {
	//make, model, package, color, mileage, price, category, year, identification
	query := "UPDATE cars SET make = $1, model = $2,package = $3,color = $4, mileage = $5, price = $6, category = $7, year = $8, identification = $9 , updated_at = $10 WHERE id = $11"
	res, err := r.db.ExecContext(ctx, query, car.Make, car.Model, car.Package, car.Color, car.Mileage, car.Price, car.Category, car.Year, car.Identification, car.UpdatedAt, car.ID)
	if err != nil {
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return
	}

	if affect != 1 {
		err = fmt.Errorf("weird behavior, total affected: %d", affect)
	}

	return
}

func (r *pgsqlCarRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM cars WHERE id = $1"
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return
	}

	if affect != 1 {
		err = fmt.Errorf("weird behavior, total affected: %d", affect)
	}

	return
}
