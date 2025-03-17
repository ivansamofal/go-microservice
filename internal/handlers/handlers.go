package handlers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"go_microservice/internal/logger"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go_microservice/internal/models"
	"go_microservice/internal/repository"
)

var getCountriesWithCities = repository.GetCountriesWithCities

// Простой обработчик, возвращающий строку
func Handler(c *gin.Context) {
	c.String(http.StatusOK, "Hello, Go Microservice!")
}

// StatusHandler godoc
// @Summary Статус сервиса
// @Description Возвращает информацию о статусе сервиса
// @Tags status
// @Produce json
// @Success 200 {object} map[string]interface{} "Статус сервиса"
// @Failure 401 {object} map[string]interface{} "Ошибка авторизации"
// @Router /api/status [get]
func StatusHandler(c *gin.Context) {
	response := map[string]string{
		"status":  "running",
		"version": "1.0.0",
	}
	c.JSON(http.StatusOK, response)
}

// InfoHandler godoc
// @Summary Информация о сервисе
// @Description Возвращает общую информацию о сервисе
// @Tags info
// @Produce json
// @Success 200 {object} map[string]interface{} "Информация о сервисе"
// @Failure 401 {object} map[string]interface{} "Ошибка авторизации"
// @Router /api/info [get]
func InfoHandler(c *gin.Context) {
	var service models.Service
	// Находим сервис с ID = 1
	if err := repository.GetServiceByID(1, &service); err != nil {
		logger.Log.WithError(err).Error("Failed to fetch service info")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch service info"})
		return
	}

	response := map[string]string{
		"name":        service.Name,
		"description": "Simple API with static data",
		"author":      "Your Name",
	}
	c.JSON(http.StatusOK, response)
}

// FetchAndSaveData godoc
// @Summary Сохранить данные
// @Description Получает данные из запроса и сохраняет их в БД (требуется авторизация)
// @Tags save
// @Accept json
// @Produce json
// @Param data body map[string]interface{} true "Данные для сохранения"
// @Success 200 {object} map[string]interface{} "Результат операции"
// @Failure 400 {object} map[string]interface{} "Неверный запрос"
// @Failure 401 {object} map[string]interface{} "Ошибка авторизации"
// @Failure 500 {object} map[string]interface{} "Ошибка БД"
// @Router /api/save [post]
func FetchAndSaveData(c *gin.Context) {
	// Настройка TLS для пропуска проверки сертификата
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get("https://restcountries.com/v3.1/all")
	if err != nil {
		log.Println(err)
		logger.Log.WithError(err).Error("Failed to fetch data from API save data")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from API"})
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		logger.Log.WithError(err).Error("Failed to read API response save data")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read API response"})
		return
	}

	// Предположим, что структура models.Country используется для парсинга JSON
	var externalCountries []models.Country
	if err = json.Unmarshal(body, &externalCountries); err != nil {
		log.Println(err)
		logger.Log.WithError(err).Error("Failed to parse API save data")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse API response"})
		return
	}

	// Перебираем полученные страны и сохраняем их через репозиторий
	for _, extCountry := range externalCountries {
		// Сохраняем страну (InsertCountry возвращает ID сохранённой страны)
		countryID, err := repository.InsertCountry(extCountry.Name.Common, extCountry.Alpha2Code, extCountry.Alpha3Code)
		if err != nil {
			logger.Log.WithError(err).Error("Failed to insert country " + extCountry.Name.Common)
			log.Printf("Failed to insert country %s: %v", extCountry.Name.Common, err)
			continue
		}

		// Если есть столица – сохраняем её как город
		if len(extCountry.Capital) > 0 {
			err := repository.InsertCity(extCountry.Capital[0], countryID, extCountry.Population, true)
			if err != nil {
				logger.Log.WithError(err).Error("Failed to insert city " + extCountry.Capital[0])
				log.Printf("Failed to insert city %s: %v", extCountry.Capital[0], err)
			} else {
				fmt.Printf("Inserted country %s and capital %s\n", extCountry.Name.Common, extCountry.Capital[0])
			}
		}
	}

	c.String(http.StatusOK, "Data saved successfully")
}
