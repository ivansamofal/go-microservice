package main

//import (
//	"github.com/gin-gonic/gin"
//	"net/http"
//)
//
//func main() {
//	r := gin.Default()
//	InitDB()
//
//	// Автоматически создаем таблицы в БД
//	DB.AutoMigrate(&User{})
//
//	// Эндпоинты
//	r.GET("/users", GetUsers)
//	r.POST("/users", CreateUser)
//
//	r.Run(":8080") // Запуск сервера
//}
//
//// Получение списка пользователей
//func GetUsers(c *gin.Context) {
//	var users []User
//	DB.Find(&users)
//	c.JSON(http.StatusOK, users)
//}
//
//// Создание пользователя
//func CreateUser(c *gin.Context) {
//	var user User
//	if err := c.ShouldBindJSON(&user); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	DB.Create(&user)
//	c.JSON(http.StatusCreated, user)
//}
