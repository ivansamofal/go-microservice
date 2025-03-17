package handlers

import (
	"net/http"

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
	var countries []models.CountryDB
	if err := GetCountriesWithCitiesFunc(&countries); err != nil {
		logger.Log.WithError(err).Error("Ошибка получения стран и городов из базы данных")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
		return
	}
	c.JSON(http.StatusOK, countries)
}
