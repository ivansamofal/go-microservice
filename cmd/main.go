package main

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go_microservice/internal/db"
	"go_microservice/internal/handlers"
	"go_microservice/internal/controllers"
	"go_microservice/internal/logger"
	"go_microservice/internal/middleware"
	"go_microservice/internal/cache"
	"github.com/joho/godotenv"
	"log"

	"github.com/gin-gonic/gin"
	_ "go_microservice/docs" // пакет сгенерированных документов (смотрите шаг 4)
)

// @title Go Microservice API
// @version 1.0
// @description API Documentation for Go Microservice.
// @host localhost:8080
// @BasePath /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Нет файла .env или ошибка его загрузки")
	}

	db.InitDB()
	logger.Init()
	cache.InitRedis()

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Публичные маршруты
	router.GET("/", handlers.Handler)
	router.POST("/api/login", handlers.GetAPIKeyHandler)

	// Защищённая группа маршрутов
	api := router.Group("/api")
	api.Use(middleware.JWTAuthMiddleware())
	{
		api.GET("/status", handlers.StatusHandler)
		api.GET("/info", handlers.InfoHandler)
		api.GET("/geo", handlers.GeoHandler)
		api.POST("/save", handlers.FetchAndSaveData)

		router.GET("/api/gdp", handlers.AverageGdpHandler)

		router.POST("/api/trade", controllers.SaveTradeData)
		router.GET("/api/calculations", controllers.Calculations)
	}

	log.Println("Starting server on port 8080...")
	log.Fatal(router.Run(":8080"))
}
