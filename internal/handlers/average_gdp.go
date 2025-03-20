package handlers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
	"math"
	"sync"

	"github.com/gin-gonic/gin"
	"go_microservice/internal/models"
	"go_microservice/internal/repository"
)

func AverageGdpHandler(c *gin.Context) {
	// Получаем список стран из БД через репозиторий
	var countryRecords []models.CountryDB
	if err := repository.GetCountries(&countryRecords); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка запроса к БД"})
		return
	}

	// Извлекаем названия стран
	var countryNames []string
	for _, rec := range countryRecords {
		countryNames = append(countryNames, rec.Name)
	}

	// Создаем HTTP-клиент с отключенной проверкой сертификата (для теста)
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	var results []models.CountryGDPResult

	// Выполняем запросы синхронно с задержкой 2 секунды между запросами
	for _, countryName := range countryNames {
		// Экранируем название страны для корректного формирования URL
		escapedCountry := url.QueryEscape(countryName)
		url := fmt.Sprintf("https://api.api-ninjas.com/v1/gdp?country=%s", escapedCountry)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println("Ошибка создания запроса для", countryName, ":", err)
			time.Sleep(2 * time.Second)
			continue
		}

		// Добавляем требуемый заголовок
		req.Header.Add("Origin", "https://www.api-ninjas.com")
		// Если нужен API-ключ, можно добавить его здесь:
		// req.Header.Add("X-Api-Key", "ваш_api_key")

		resp, err := client.Do(req)
		if err != nil {
			log.Println("Ошибка выполнения запроса для", countryName, ":", err)
			time.Sleep(2 * time.Second)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Println("Некорректный статус для", countryName, ":", resp.Status)
			time.Sleep(2 * time.Second)
			continue
		}

		// Декодирование полученного JSON-массива
		var data []models.GDPData
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Println("Ошибка декодирования JSON для", countryName, ":", err)
			time.Sleep(2 * time.Second)
			continue
		}

		// Если данных нет, пропускаем страну
		if len(data) == 0 {
			log.Println("Нет данных для", countryName)
			time.Sleep(2 * time.Second)
			continue
		}

		// Агрегация: считаем среднее значение gdp_nominal
		numWorkers := 4
		n := len(data)
		chunkSize := (n + numWorkers - 1) / numWorkers

		partialSums := make(chan float64, numWorkers)
		var wg sync.WaitGroup

		for i := 0; i < n; i += chunkSize {
			wg.Add(1)
			go func(start int) {
				defer wg.Done()
				end := start + chunkSize
				if end > n {
					end = n
				}
				var partial float64
				for j := start; j < end; j++ {
					partial += data[j].GDPNominal
				}
				partialSums <- partial
			}(i)
		}

		wg.Wait()
		close(partialSums)

		var total float64
		for psum := range partialSums {
			total += psum
		}

		avg := total / float64(n)
		avg = math.Round(avg*100) / 100

		results = append(results, models.CountryGDPResult{
			Country:           countryName,
			AverageGDPNominal: avg,
		})

		// Задержка 2 секунды перед следующим запросом
		time.Sleep(2 * time.Second)
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}
