package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go_microservice/internal/cache"
	"go_microservice/internal/logger"
	"go_microservice/internal/models"
	"go_microservice/internal/repository"

	"github.com/gin-gonic/gin"
)

// GetCountriesWithCitiesFunc – экспортированная переменная для получения стран и городов.
// В продакшене она ссылается на repository.GetCountriesWithCities,
// а в тестах её можно переопределить.
var GetCountriesWithCitiesFunc = repository.GetCountriesWithCities

// GeoHandler godoc
// @Summary Географическая информация
// @Description Возвращает список стран и городов (требуется авторизация)
// @Tags geo
// @Produce json
// @Success 200 {array} models.CountryDB "Список стран с городами"
// @Failure 401 {object} map[string]interface{} "Ошибка авторизации"
// @Failure 500 {object} map[string]interface{} "Ошибка БД"
// @Router /api/geo [get]
func GeoHandler(c *gin.Context) {
	ctx := context.Background()
	cacheKey := "countries"

	// Пытаемся получить данные из Redis
	cached, err := cache.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var countries []models.CountryDB
		if err := json.Unmarshal([]byte(cached), &countries); err == nil {
			c.JSON(http.StatusOK, countries)
			return
		}
		// Если ошибка при парсинге, логируем её и продолжаем – данные из кеша могут быть повреждены.
		logger.Log.WithError(err).Error("Ошибка декодирования кешированных данных")
	}

	// Если данные не найдены в кеше или произошла ошибка, получаем данные из БД.
	var countries []models.CountryDB
	if err := GetCountriesWithCitiesFunc(&countries); err != nil {
		logger.Log.WithError(err).Error("Ошибка получения стран и городов из базы данных")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
		return
	}

	// Сохраняем результат в Redis с TTL, например, 10 минут
	if bytes, err := json.Marshal(countries); err == nil {
		cache.RedisClient.Set(ctx, cacheKey, bytes, 10*time.Minute)
	} else {
		logger.Log.WithError(err).Error("Ошибка сериализации данных для кеша")
	}

	c.JSON(http.StatusOK, countries)
}
