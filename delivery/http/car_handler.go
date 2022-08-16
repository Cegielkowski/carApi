package http

import (
	"net/http"
	"strconv"

	"carApi/delivery/middleware"
	"carApi/transport/request"
	"carApi/usecase"
	"carApi/utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
)

type CarHandler struct {
	CarUC usecase.CarUsecase
}

// NewCarHandler will initialize the cars / resources endpoint
func NewCarHandler(e *echo.Echo, middleware *middleware.Middleware, carUC usecase.CarUsecase) {
	handler := &CarHandler{
		CarUC: carUC,
	}

	apiV1 := e.Group("/api/v1")
	apiV1.POST("/cars", handler.Create)
	apiV1.GET("/cars/:id", handler.GetByID)
	apiV1.GET("/cars", handler.Fetch)
	apiV1.PUT("/cars/:id", handler.Update)
	apiV1.DELETE("/cars/:id", handler.Delete)
}

func (h *CarHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.CreateCarReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if err := h.CarUC.Create(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "car created",
	})

}

func (h *CarHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewNotFoundError("car not found"))
	}

	car, err := h.CarUC.GetByID(ctx, int64(id))
	if err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": car})
}

func (h *CarHandler) Fetch(c echo.Context) error {
	ctx := c.Request().Context()

	cars, err := h.CarUC.Fetch(ctx)
	if err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": cars})
}

func (h *CarHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewNotFoundError("car not found"))
	}

	var req request.UpdateCarReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if err := h.CarUC.Update(ctx, int64(id), &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "car updated",
	})
}

func (h *CarHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewNotFoundError("car not found"))
	}

	if err := h.CarUC.Delete(ctx, int64(id)); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "car deleted",
	})
}
