package main

import (
	"net/http"
	"time"

	_ "carApi/docs"
	"carApi/utils"

	"carApi/config"
	httpDelivery "carApi/delivery/http"
	appMiddleware "carApi/delivery/middleware"
	"carApi/infrastructure/datastore"
	pgsqlRepository "carApi/repository/pgsql"
	redisRepository "carApi/repository/redis"
	"carApi/usecase"
	"carApi/utils/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	// Load config
	configApp := config.LoadConfig()

	// Setup logger
	appLogger := logger.NewApiLogger(configApp)
	appLogger.InitLogger()

	// Setup infra
	dbInstance, err := datastore.NewDatabase(configApp.DatabaseURL)
	utils.PanicIfNeeded(err)

	cacheInstance, err := datastore.NewCache(configApp.CacheURL)
	utils.PanicIfNeeded(err)

	// Setup repository
	redisRepo := redisRepository.NewRedisRepository(cacheInstance)
	carRepo := pgsqlRepository.NewPgsqlCarRepository(dbInstance)

	// Setup usecase
	ctxTimeout := time.Duration(configApp.ContextTimeout) * time.Second
	carUC := usecase.NewCarUsecase(carRepo, redisRepo, ctxTimeout)

	// Setup app middleware
	appMiddleware := appMiddleware.NewMiddleware(appLogger)

	// Setup route engine & middleware
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(appMiddleware.RequestID())
	e.Use(appMiddleware.Logger())
	e.Use(middleware.Recover())

	// Setup handler
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "i am alive")
	})

	httpDelivery.NewCarHandler(e, appMiddleware, carUC)

	e.Logger.Fatal(e.Start(":" + configApp.ServerPORT))
}
