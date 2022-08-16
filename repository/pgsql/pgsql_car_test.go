package pgsql_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"carApi/entity"
	"carApi/repository/pgsql"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestCarRepo_Create(t *testing.T) {
	car := &entity.Car{
		Make:           "Make",
		Model:          "Model",
		Package:        "Package",
		Color:          "Color",
		Year:           0,
		Category:       "Category",
		Mileage:        0,
		Price:          0,
		Identification: "Identification",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := "INSERT INTO cars"
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(car.Make, car.Model, car.Package, car.Color, car.Mileage, car.Price, car.Category, car.Year, car.Identification, car.CreatedAt, car.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	carRepo := pgsql.NewPgsqlCarRepository(db)
	err = carRepo.Create(context.TODO(), car)
	assert.NoError(t, err)
}

func TestCarRepo_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	carMock := entity.Car{
		ID:             1,
		Make:           "Make",
		Model:          "Model",
		Package:        "Package",
		Color:          "Color",
		Year:           1,
		Category:       "Category",
		Mileage:        1,
		Price:          1,
		Identification: "Identification",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "make", "model", "package", "color", "mileage", "price", "category", "year", "identification", "created_at", "updated_at"}).
		AddRow(carMock.ID, carMock.Make, carMock.Model, carMock.Package, carMock.Color, carMock.Mileage, carMock.Price, carMock.Category, carMock.Year, carMock.Identification, carMock.CreatedAt, carMock.UpdatedAt)

	query := "SELECT id, make, model, package, color, mileage, price, category, year, identification, created_at, updated_at FROM cars WHERE id = $1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(1).
		WillReturnRows(rows)

	carRepo := pgsql.NewPgsqlCarRepository(db)
	car, err := carRepo.GetByID(context.TODO(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, car)
	assert.Equal(t, carMock.ID, car.ID)
}

func TestCarRepo_Fetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockCars := []entity.Car{
		{
			ID:             1,
			Make:           "Make",
			Model:          "Model",
			Package:        "Package",
			Color:          "Color",
			Year:           0,
			Category:       "Category",
			Mileage:        0,
			Price:          0,
			Identification: "Identification",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             2,
			Make:           "Make",
			Model:          "Model2",
			Package:        "Package2",
			Color:          "Color2",
			Year:           0,
			Category:       "Category2",
			Mileage:        0,
			Price:          0,
			Identification: "Identification2",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "make", "model", "package", "color", "mileage", "price", "category", "year", "identification", "created_at", "updated_at"}).
		AddRow(mockCars[0].ID, mockCars[0].Make, mockCars[0].Model, mockCars[0].Package, mockCars[0].Color, mockCars[0].Mileage, mockCars[0].Price, mockCars[0].Category, mockCars[0].Year, mockCars[0].Identification, mockCars[0].CreatedAt, mockCars[0].UpdatedAt).
		AddRow(mockCars[1].ID, mockCars[1].Make, mockCars[1].Model, mockCars[1].Package, mockCars[1].Color, mockCars[1].Mileage, mockCars[1].Price, mockCars[1].Category, mockCars[1].Year, mockCars[1].Identification, mockCars[1].CreatedAt, mockCars[1].UpdatedAt)

	query := "SELECT id, make, model, package, color, mileage, price, category, year, identification, created_at, updated_at FROM cars"
	mock.ExpectQuery(query).WillReturnRows(rows)

	carRepo := pgsql.NewPgsqlCarRepository(db)
	cars, err := carRepo.Fetch(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, cars, 2)
}

func TestCarRepo_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	carMock := &entity.Car{
		ID:             1,
		Make:           "Make",
		Model:          "Model",
		Package:        "Package",
		Color:          "Color",
		Year:           0,
		Category:       "Category",
		Mileage:        0,
		Price:          0,
		Identification: "Identification",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	query := "UPDATE cars SET make = $1, model = $2,package = $3,color = $4, mileage = $5, price = $6, category = $7, year = $8, identification = $9 , updated_at = $10 WHERE id = $11"
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(carMock.Make, carMock.Model, carMock.Package, carMock.Color, carMock.Mileage, carMock.Price, carMock.Category, carMock.Year, carMock.Identification, carMock.UpdatedAt, carMock.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	carRepo := pgsql.NewPgsqlCarRepository(db)
	err = carRepo.Update(context.TODO(), carMock)
	assert.NoError(t, err)
}

func TestCarRepo_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := "DELETE FROM cars WHERE id = $1"
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	carRepo := pgsql.NewPgsqlCarRepository(db)
	err = carRepo.Delete(context.TODO(), 1)
	assert.NoError(t, err)
}
