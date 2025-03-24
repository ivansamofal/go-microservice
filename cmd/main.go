package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"

	"go_microservice/internal/cache"
	"go_microservice/internal/controllers"
	"go_microservice/internal/db"
	"go_microservice/internal/handlers"
	"go_microservice/internal/logger"
	"go_microservice/internal/middleware"

	_ "go_microservice/docs" // пакет сгенерированных документов
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

	// Применяем CORS middleware ко всем маршрутам
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("HOST_FRONTEND_URL")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Документация Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Публичные маршруты
	router.GET("/", handlers.Handler)
	router.POST("/api/login", handlers.GetAPIKeyHandler)
	router.GET("/api/calculations", controllers.Calculations)
	router.POST("/trade", controllers.SaveTradeData)

	// Защищённая группа маршрутов
	api := router.Group("/api")
	api.Use(middleware.JWTAuthMiddleware())
	{
		api.GET("/status", handlers.StatusHandler)
		api.GET("/info", handlers.InfoHandler)
		api.GET("/geo", handlers.GeoHandler)
		api.POST("/save", handlers.FetchAndSaveData)

		// Все маршруты начинающиеся с /api/ будут применяться с JWT middleware
		api.GET("/gdp", handlers.AverageGdpHandler)
	}

	log.Println("Starting server on port 8080...")
	log.Fatal(router.Run(":8080"))
}
