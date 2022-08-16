package http_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	httpDelivery "carApi/delivery/http"
	"carApi/entity"
	"carApi/mocks"
	"carApi/transport/request"
	"carApi/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCarHandler_Create(t *testing.T) {
	mockCarUC := new(mocks.CarUsecase)
	createCarReq := request.CreateCarReq{
		Make:           "Make",
		Model:          "Model",
		Package:        "Package",
		Color:          "Color",
		Year:           0,
		Category:       "Category",
		Mileage:        0,
		Price:          0,
		Identification: "Identification",
	}

	t.Run("success", func(t *testing.T) {
		jsonReq, err := json.Marshal(createCarReq)
		assert.NoError(t, err)

		mockCarUC.On("Create", mock.Anything, mock.AnythingOfType("*request.CreateCarReq")).
			Return(nil).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/api/v1/cars", strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/cars")

		handler := httpDelivery.CarHandler{
			CarUC: mockCarUC,
		}
		err = handler.Create(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockCarUC.AssertExpectations(t)
	})

	t.Run("error-validation", func(t *testing.T) {
		invalidCreateCarReq := request.CreateCarReq{
			Make:           "Make",
			Model:          "Model",
			Package:        "Package",
			Color:          "Color",
			Year:           0,
			Category:       "Category",
			Mileage:        0,
			Price:          0,
			Identification: "Identification",
		}
		jsonReq, err := json.Marshal(invalidCreateCarReq)
		assert.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/api/v1/cars", strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/cars")

		handler := httpDelivery.CarHandler{
			CarUC: mockCarUC,
		}
		err = handler.Create(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockCarUC.AssertExpectations(t)
	})

	t.Run("error-usecase", func(t *testing.T) {
		jsonReq, err := json.Marshal(createCarReq)
		assert.NoError(t, err)

		mockCarUC.On("Create", mock.Anything, mock.AnythingOfType("*request.CreateCarReq")).
			Return(errors.New("Unexpected Error")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/api/v1/cars", strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/cars")

		handler := httpDelivery.CarHandler{
			CarUC: mockCarUC,
		}
		err = handler.Create(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockCarUC.AssertExpectations(t)
	})

}

func TestCarHandler_GetByID(t *testing.T) {
	mockCarUC := new(mocks.CarUsecase)
	mockCar := entity.Car{
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

	t.Run("success", func(t *testing.T) {
		mockCarUC.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
			Return(mockCar, nil).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/cars/"+strconv.Itoa(int(mockCar.ID)), strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/cars/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockCar.ID)))

		handler := httpDelivery.CarHandler{
			CarUC: mockCarUC,
		}
		err = handler.GetByID(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockCarUC.AssertExpectations(t)
	})

	t.Run("data-not-exist", func(t *testing.T) {
		mockCarUC.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
			Return(entity.Car{}, utils.NewNotFoundError("car not found")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/cars/"+strconv.Itoa(int(mockCar.ID)), strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/cars/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockCar.ID)))

		handler := httpDelivery.CarHandler{
			CarUC: mockCarUC,
		}
		err = handler.GetByID(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockCarUC.AssertExpectations(t)
	})

	t.Run("error-usecase", func(t *testing.T) {
		mockCarUC.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
			Return(entity.Car{}, errors.New("Unexpected Error")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/cars/"+strconv.Itoa(int(mockCar.ID)), strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/cars/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockCar.ID)))

		handler := httpDelivery.CarHandler{
			CarUC: mockCarUC,
		}
		err = handler.GetByID(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockCarUC.AssertExpectations(t)
	})
}

func TestCarHandler_Fetch(t *testing.T) {
	mockCarUC := new(mocks.CarUsecase)
	mockCar := entity.Car{
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

	mockListCar := make([]entity.Car, 0)
	mockListCar = append(mockListCar, mockCar)

	t.Run("success", func(t *testing.T) {
		mockCarUC.On("Fetch", mock.Anything).Return(mockListCar, nil).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/cars/", strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/cars/")

		handler := httpDelivery.CarHandler{
			CarUC: mockCarUC,
		}
		err = handler.Fetch(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockCarUC.AssertExpectations(t)
	})

	t.Run("error-usecase", func(t *testing.T) {
		mockCarUC.On("Fetch", mock.Anything).Return([]entity.Car{}, errors.New("Unexpected Error")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/api/v1/cars/", strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/cars/")

		handler := httpDelivery.CarHandler{
			CarUC: mockCarUC,
		}
		err = handler.Fetch(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockCarUC.AssertExpectations(t)
	})
}

func TestCarHandler_Update(t *testing.T) {
	mockCarUC := new(mocks.CarUsecase)
	mockCar := entity.Car{
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
	updateCarReq := request.UpdateCarReq{
		Make:           "Make2",
		Model:          "Model2",
		Package:        "Package2",
		Color:          "Color2",
		Year:           0,
		Category:       "Category2",
		Mileage:        0,
		Price:          0,
		Identification: "Identification2",
	}

	t.Run("success", func(t *testing.T) {
		jsonReq, err := json.Marshal(updateCarReq)
		assert.NoError(t, err)

		mockCarUC.On("Update", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("*request.UpdateCarReq")).
			Return(nil).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/cars/"+strconv.Itoa(int(mockCar.ID)), strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/cars/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockCar.ID)))

		handler := httpDelivery.CarHandler{
			CarUC: mockCarUC,
		}
		err = handler.Update(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockCarUC.AssertExpectations(t)
	})

	t.Run("error-validation", func(t *testing.T) {
		invalidUpdateCarReq := request.UpdateCarReq{
			Make: "",
		}
		jsonReq, err := json.Marshal(invalidUpdateCarReq)
		assert.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/cars/"+strconv.Itoa(int(mockCar.ID)), strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/cars/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockCar.ID)))

		handler := httpDelivery.CarHandler{
			CarUC: mockCarUC,
		}
		err = handler.Update(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockCarUC.AssertExpectations(t)
	})

	t.Run("data-not-exist", func(t *testing.T) {
		jsonReq, err := json.Marshal(updateCarReq)
		assert.NoError(t, err)

		mockCarUC.On("Update", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("*request.UpdateCarReq")).
			Return(utils.NewNotFoundError("car not found")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/cars/"+strconv.Itoa(int(mockCar.ID)), strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/cars/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockCar.ID)))

		handler := httpDelivery.CarHandler{
			CarUC: mockCarUC,
		}
		err = handler.Update(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockCarUC.AssertExpectations(t)
	})

	t.Run("error-usecase", func(t *testing.T) {
		jsonReq, err := json.Marshal(updateCarReq)
		assert.NoError(t, err)

		mockCarUC.On("Update", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("*request.UpdateCarReq")).
			Return(errors.New("Unexpected Error")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/api/v1/cars/"+strconv.Itoa(int(mockCar.ID)), strings.NewReader(string(jsonReq)))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/cars/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockCar.ID)))

		handler := httpDelivery.CarHandler{
			CarUC: mockCarUC,
		}
		err = handler.Update(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockCarUC.AssertExpectations(t)
	})

}

func TestCarHandler_Delete(t *testing.T) {
	mockCarUC := new(mocks.CarUsecase)
	mockCar := entity.Car{
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

	t.Run("success", func(t *testing.T) {
		mockCarUC.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.DELETE, "/api/v1/cars/"+strconv.Itoa(int(mockCar.ID)), strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/cars/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockCar.ID)))

		handler := httpDelivery.CarHandler{
			CarUC: mockCarUC,
		}
		err = handler.Delete(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockCarUC.AssertExpectations(t)
	})

	t.Run("data-not-exist", func(t *testing.T) {
		mockCarUC.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(utils.NewNotFoundError("car not found")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.DELETE, "/api/v1/cars/"+strconv.Itoa(int(mockCar.ID)), strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/cars/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockCar.ID)))

		handler := httpDelivery.CarHandler{
			CarUC: mockCarUC,
		}
		err = handler.Delete(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockCarUC.AssertExpectations(t)
	})

	t.Run("error-usecase", func(t *testing.T) {
		mockCarUC.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(errors.New("Unexpected Error")).Once()

		e := echo.New()
		req, err := http.NewRequest(echo.DELETE, "/api/v1/cars/"+strconv.Itoa(int(mockCar.ID)), strings.NewReader(""))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/cars/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(mockCar.ID)))

		handler := httpDelivery.CarHandler{
			CarUC: mockCarUC,
		}
		err = handler.Delete(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockCarUC.AssertExpectations(t)
	})
}
