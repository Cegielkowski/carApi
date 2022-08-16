package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"carApi/entity"
	"carApi/repository/pgsql"
	"carApi/repository/redis"
	"carApi/transport/request"
	"carApi/utils"
)

// CarUsecase represent the car's usecase contract
type CarUsecase interface {
	Create(ctx context.Context, request *request.CreateCarReq) error
	GetByID(ctx context.Context, id int64) (entity.Car, error)
	Fetch(ctx context.Context) ([]entity.Car, error)
	Update(ctx context.Context, id int64, request *request.UpdateCarReq) error
	Delete(ctx context.Context, id int64) error
}

type carUsecase struct {
	carRepo    pgsql.CarRepository
	redisRepo  redis.RedisRepository
	ctxTimeout time.Duration
}

// NewCarUsecase will create new an carUsecase object representation of CarUsecase interface
func NewCarUsecase(carRepo pgsql.CarRepository, redisRepo redis.RedisRepository, ctxTimeout time.Duration) CarUsecase {
	return &carUsecase{
		carRepo:    carRepo,
		redisRepo:  redisRepo,
		ctxTimeout: ctxTimeout,
	}
}

func (u *carUsecase) Create(c context.Context, request *request.CreateCarReq) (err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	err = u.carRepo.Create(ctx, &entity.Car{
		Make:           request.Make,
		Model:          request.Model,
		Package:        request.Package,
		Color:          request.Color,
		Year:           request.Year,
		Category:       request.Category,
		Mileage:        request.Mileage,
		Price:          request.Price,
		Identification: request.Identification,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	return
}

func (u *carUsecase) GetByID(c context.Context, id int64) (car entity.Car, err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	car, err = u.carRepo.GetByID(ctx, id)
	if err != nil && err == sql.ErrNoRows {
		err = utils.NewNotFoundError("car not found")
		return
	}
	return
}

func (u *carUsecase) Fetch(c context.Context) (cars []entity.Car, err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	carsCached, _ := u.redisRepo.Get("cars")
	if err = json.Unmarshal([]byte(carsCached), &cars); err == nil {
		return
	}

	cars, err = u.carRepo.Fetch(ctx)
	if err != nil {
		return
	}

	carsString, _ := json.Marshal(&cars)
	u.redisRepo.Set("cars", carsString, 30*time.Second)
	return
}

func (u *carUsecase) Update(c context.Context, id int64, request *request.UpdateCarReq) (err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	car, err := u.carRepo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.NewNotFoundError("car not found")
			return
		}
		return
	}

	car.Identification = request.Identification
	car.Price = request.Price
	car.Mileage = request.Mileage
	car.Year = request.Year
	car.Color = request.Color
	car.Make = request.Make
	car.Model = request.Category
	car.Package = request.Package
	car.Category = request.Category
	car.UpdatedAt = time.Now()

	err = u.carRepo.Update(ctx, &car)
	return
}

func (u *carUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	_, err = u.carRepo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.NewNotFoundError("car not found")
			return
		}
		return
	}

	err = u.carRepo.Delete(ctx, id)
	return
}
