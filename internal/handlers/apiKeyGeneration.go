package handlers

import (
	"go_microservice/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go_microservice/internal/middleware"
)

// GetAPIKeyHandler godoc
// @Summary Получение API-ключа
// @Description Генерирует JWT токен для доступа к API по переданным учетным данным.
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.Credentials true "Учетные данные"
// @Success 200 {object} map[string]string "api_key"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Неверные учетные данные"
// @Router /api/login [post]
func GetAPIKeyHandler(c *gin.Context) {
	// Структура запроса с учетными данными.
	var creds models.Credentials

	// Попытка распарсить JSON из тела запроса.
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	// Проверка учетных данных (пример, замените на свою логику аутентификации).
	if creds.Username != "admin" || creds.Password != "password" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные учетные данные"})
		return
	}

	// Создание claims для JWT токена.
	claims := jwt.MapClaims{
		"username": creds.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // Токен действителен 24 часа.
	}

	// Создаем новый JWT токен.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.JWTKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать токен"})
		return
	}

	// Возвращаем сгенерированный токен клиенту.
	c.JSON(http.StatusOK, gin.H{"api_key": tokenString})
}
