package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go_microservice/internal/handlers"
	"go_microservice/internal/models"
)

func TestGeoHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Сохраняем оригинальную функцию и заменяем её на тестовую.
	originalFunc := handlers.GetCountriesWithCitiesFunc
	handlers.GetCountriesWithCitiesFunc = func(countries *[]models.CountryDB) error {
		*countries = []models.CountryDB{
			{ID: 1, Name: "Test Country"},
		}
		return nil
	}
	defer func() {
		handlers.GetCountriesWithCitiesFunc = originalFunc
	}()

	handlers.GeoHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []models.CountryDB
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, "Test Country", response[0].Name)
}

func TestGeoHandler_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Подменяем функцию для эмуляции ошибки.
	originalFunc := handlers.GetCountriesWithCitiesFunc
	handlers.GetCountriesWithCitiesFunc = func(countries *[]models.CountryDB) error {
		return errors.New("database error")
	}
	defer func() {
		handlers.GetCountriesWithCitiesFunc = originalFunc
	}()

	handlers.GeoHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Database query error", response["error"])
}
