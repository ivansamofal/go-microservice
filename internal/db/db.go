package db

import (
	"log"
	"os"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"go_microservice/internal/logger"
	"go_microservice/internal/migrations"
)

var DB *gorm.DB

func InitDB() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pswd := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, pswd, dbName, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Log.WithError(err).Error("Failed to connect to database");
		log.Fatal("Failed to connect to database: ", err)
	}

	if err := db.AutoMigrate(&migrations.TradeRow{}); err != nil {
		log.Fatalf("Ошибка миграции БД: %v", err)
	}

	DB = db
	log.Println("Connected to the database!")
}
