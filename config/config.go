package config

//import (
//	"gorm.io/driver/postgres"
//	"gorm.io/gorm"
//	"log"
//)
//
//var DB *gorm.DB
//
//func InitDB() {
//	dsn := "host=localhost user=postgres password=secret dbname=go_microservice port=5432 sslmode=disable"
//	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
//	if err != nil {
//		log.Fatal("Ошибка подключения к БД:", err)
//	}
//	DB = db
//}
//
//import (
//"fmt"
//"gorm.io/driver/postgres"
//"gorm.io/gorm"
//"log"
//"os"
//)
//
//var DB *gorm.DB
//
//func InitDB() {
//	dsn := fmt.Sprintf(
//		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
//		os.Getenv("DB_HOST"),
//		os.Getenv("DB_USER"),
//		os.Getenv("DB_PASSWORD"),
//		os.Getenv("DB_NAME"),
//		os.Getenv("DB_PORT"),
//	)
//
//	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
//	if err != nil {
//		log.Fatal("Ошибка подключения к БД:", err)
//	}
//	DB = db
//}
