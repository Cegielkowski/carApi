package usecase_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"carApi/entity"
	"carApi/mocks"
	"carApi/transport/request"
	"carApi/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var ctxTimeout = 60 * time.Second

func TestCarUC_Create(t *testing.T) {
	mockRedisRepo := new(mocks.RedisRepository)
	mockCarRepo := new(mocks.CarRepository)
	createCarReq := request.CreateCarReq{
		Make:           "make",
		Model:          "model",
		Package:        "package",
		Color:          "color",
		Year:           0,
		Category:       "category",
		Mileage:        0,
		Price:          0,
		Identification: "identification",
	}

	t.Run("success", func(t *testing.T) {
		mockCarRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Car")).Return(nil).Once()

		carUsecase := usecase.NewCarUsecase(mockCarRepo, mockRedisRepo, ctxTimeout)
		err := carUsecase.Create(context.TODO(), &createCarReq)

		assert.NoError(t, err)
		mockRedisRepo.AssertExpectations(t)
		mockCarRepo.AssertExpectations(t)
	})

	t.Run("error-db", func(t *testing.T) {
		mockCarRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Car")).Return(errors.New("Unexpected Error")).Once()

		carUsecase := usecase.NewCarUsecase(mockCarRepo, mockRedisRepo, ctxTimeout)
		err := carUsecase.Create(context.TODO(), &createCarReq)

		assert.NotNil(t, err)
		mockRedisRepo.AssertExpectations(t)
		mockCarRepo.AssertExpectations(t)
	})
}

func TestCarUC_GetByID(t *testing.T) {
	mockRedisRepo := new(mocks.RedisRepository)
	mockCarRepo := new(mocks.CarRepository)
	mockCar := entity.Car{
		ID:             1,
		Make:           "make",
		Model:          "model",
		Package:        "package",
		Color:          "color",
		Year:           0,
		Category:       "category",
		Mileage:        0,
		Price:          0,
		Identification: "identification",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockCarRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockCar, nil).Once()

		carUsecase := usecase.NewCarUsecase(mockCarRepo, mockRedisRepo, ctxTimeout)
		car, err := carUsecase.GetByID(context.TODO(), mockCar.ID)

		assert.NoError(t, err)
		assert.NotNil(t, car)
		assert.Equal(t, car.ID, mockCar.ID)
		mockRedisRepo.AssertExpectations(t)
		mockCarRepo.AssertExpectations(t)
	})

	t.Run("car-not-exist", func(t *testing.T) {
		mockCarRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(entity.Car{}, sql.ErrNoRows).Once()

		carUsecase := usecase.NewCarUsecase(mockCarRepo, mockRedisRepo, ctxTimeout)
		car, err := carUsecase.GetByID(context.TODO(), mockCar.ID)

		assert.NotNil(t, err)
		assert.Equal(t, car, entity.Car{})
		mockRedisRepo.AssertExpectations(t)
		mockCarRepo.AssertExpectations(t)
	})

	t.Run("error-db", func(t *testing.T) {
		mockCarRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(entity.Car{}, errors.New("Unexpected Error")).Once()

		carUsecase := usecase.NewCarUsecase(mockCarRepo, mockRedisRepo, ctxTimeout)
		car, err := carUsecase.GetByID(context.TODO(), mockCar.ID)

		assert.NotNil(t, err)
		assert.Equal(t, car, entity.Car{})
		mockRedisRepo.AssertExpectations(t)
		mockCarRepo.AssertExpectations(t)
	})
}

func TestCarUC_Fetch(t *testing.T) {
	mockRedisRepo := new(mocks.RedisRepository)
	mockCarRepo := new(mocks.CarRepository)
	mockCar := entity.Car{
		ID:             1,
		Make:           "make",
		Model:          "model",
		Package:        "package",
		Color:          "color",
		Year:           0,
		Category:       "category",
		Mileage:        0,
		Price:          0,
		Identification: "identification",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	mockListCar := make([]entity.Car, 0)
	mockListCar = append(mockListCar, mockCar)

	t.Run("success", func(t *testing.T) {
		mockRedisRepo.On("Get", mock.AnythingOfType("string")).Return("", errors.New("Unexpected Error")).Once()
		mockCarRepo.On("Fetch", mock.Anything).Return(mockListCar, nil).Once()
		mockRedisRepo.On("Set", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).Return(nil).Once()

		carUsecase := usecase.NewCarUsecase(mockCarRepo, mockRedisRepo, ctxTimeout)
		cars, err := carUsecase.Fetch(context.TODO())

		assert.NoError(t, err)
		assert.Len(t, cars, len(mockListCar))
		mockRedisRepo.AssertExpectations(t)
		mockCarRepo.AssertExpectations(t)
	})

	t.Run("success-get-from-cache", func(t *testing.T) {
		mockListCarByte, _ := json.Marshal(mockListCar)
		mockRedisRepo.On("Get", mock.AnythingOfType("string")).Return(string(mockListCarByte), nil).Once()

		carUsecase := usecase.NewCarUsecase(mockCarRepo, mockRedisRepo, ctxTimeout)
		cars, err := carUsecase.Fetch(context.TODO())

		assert.NoError(t, err)
		assert.Len(t, cars, len(mockListCar))
		mockRedisRepo.AssertExpectations(t)
		mockCarRepo.AssertExpectations(t)

	})

	t.Run("error-db", func(t *testing.T) {
		mockRedisRepo.On("Get", mock.AnythingOfType("string")).Return("", errors.New("Unexpected Error")).Once()
		mockCarRepo.On("Fetch", mock.Anything).Return([]entity.Car{}, errors.New("Unexpected Error")).Once()

		carUsecase := usecase.NewCarUsecase(mockCarRepo, mockRedisRepo, ctxTimeout)
		cars, err := carUsecase.Fetch(context.TODO())

		assert.NotNil(t, err)
		assert.Len(t, cars, 0)
		mockRedisRepo.AssertExpectations(t)
		mockCarRepo.AssertExpectations(t)
	})
}

func TestCarUC_Update(t *testing.T) {
	mockRedisRepo := new(mocks.RedisRepository)
	mockCarRepo := new(mocks.CarRepository)
	mockCar := entity.Car{
		ID:             1,
		Make:           "make",
		Model:          "model",
		Package:        "package",
		Color:          "color",
		Year:           0,
		Category:       "category",
		Mileage:        0,
		Price:          0,
		Identification: "identification",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	updateCarReq := request.UpdateCarReq{
		Make:           "make",
		Model:          "model",
		Package:        "package",
		Color:          "color",
		Year:           0,
		Category:       "category",
		Mileage:        0,
		Price:          0,
		Identification: "identification",
	}

	t.Run("success", func(t *testing.T) {
		mockCarRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockCar, nil).Once()
		mockCarRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Car")).Return(nil).Once()

		carUsecase := usecase.NewCarUsecase(mockCarRepo, mockRedisRepo, ctxTimeout)
		err := carUsecase.Update(context.TODO(), mockCar.ID, &updateCarReq)

		assert.NoError(t, err)
		mockRedisRepo.AssertExpectations(t)
		mockCarRepo.AssertExpectations(t)
	})

	t.Run("car-not-exist", func(t *testing.T) {
		mockCarRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(entity.Car{}, sql.ErrNoRows).Once()

		carUsecase := usecase.NewCarUsecase(mockCarRepo, mockRedisRepo, ctxTimeout)
		err := carUsecase.Update(context.TODO(), mockCar.ID, &updateCarReq)

		assert.NotNil(t, err)
		mockRedisRepo.AssertExpectations(t)
		mockCarRepo.AssertExpectations(t)
	})

	t.Run("error-db", func(t *testing.T) {
		mockCarRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockCar, nil).Once()
		mockCarRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Car")).Return(errors.New("Unexpected Error")).Once()

		carUsecase := usecase.NewCarUsecase(mockCarRepo, mockRedisRepo, ctxTimeout)
		err := carUsecase.Update(context.TODO(), mockCar.ID, &updateCarReq)

		assert.NotNil(t, err)
		mockRedisRepo.AssertExpectations(t)
		mockCarRepo.AssertExpectations(t)
	})
}

func TestCarUC_Delete(t *testing.T) {
	mockRedisRepo := new(mocks.RedisRepository)
	mockCarRepo := new(mocks.CarRepository)
	mockCar := entity.Car{
		ID:             1,
		Make:           "make",
		Model:          "model",
		Package:        "package",
		Color:          "color",
		Year:           0,
		Category:       "category",
		Mileage:        0,
		Price:          0,
		Identification: "identification",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mockCarRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockCar, nil).Once()
		mockCarRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()

		carRepository := usecase.NewCarUsecase(mockCarRepo, mockRedisRepo, ctxTimeout)
		err := carRepository.Delete(context.TODO(), mockCar.ID)

		assert.NoError(t, err)
		mockRedisRepo.AssertExpectations(t)
		mockCarRepo.AssertExpectations(t)
	})

	t.Run("car-not-exist", func(t *testing.T) {
		mockCarRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(entity.Car{}, sql.ErrNoRows).Once()

		carRepository := usecase.NewCarUsecase(mockCarRepo, mockRedisRepo, ctxTimeout)
		err := carRepository.Delete(context.TODO(), mockCar.ID)

		assert.NotNil(t, err)
		mockRedisRepo.AssertExpectations(t)
		mockCarRepo.AssertExpectations(t)
	})

	t.Run("error-db", func(t *testing.T) {
		mockCarRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockCar, nil).Once()
		mockCarRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(errors.New("Unexpected Error")).Once()

		carRepository := usecase.NewCarUsecase(mockCarRepo, mockRedisRepo, ctxTimeout)
		err := carRepository.Delete(context.TODO(), mockCar.ID)

		assert.NotNil(t, err)
		mockRedisRepo.AssertExpectations(t)
		mockCarRepo.AssertExpectations(t)
	})
}
