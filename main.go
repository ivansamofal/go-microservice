package main

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
)

var db *sql.DB

type NativeName struct {
	Eng struct {
		Official string `json:"official"`
		Common   string `json:"common"`
	} `json:"eng"`
}

type Name struct {
	Common     string     `json:"common"`
	Official   string     `json:"official"`
	NativeName NativeName `json:"nativeName"`
}

type Country struct {
	Name       Name     `json:"name"`
	Alpha2Code string   `json:"cca2"`
	Alpha3Code string   `json:"cca3"`
	Capital    []string `json:"capital"`
	Population int      `json:"population"`
}

func initDB() {
	var err error
	// Замените на ваш URL подключения
	connStr := "user=app_user password=mypswd host=postgres port=5432 dbname=app_db sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Проверяем соединение с БД
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the database!")
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, Go Microservice!")
}

// JSON-ответ с информацией о статусе сервиса
func statusHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "running",
		"version": "1.0.0",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// JSON-ответ с информацией о сервисе
func infoHandler(w http.ResponseWriter, r *http.Request) {
	// Пример запроса к базе данных
	var name string
	err := db.QueryRow("SELECT name FROM services WHERE id = 1").Scan(&name)
	if err != nil {
		http.Error(w, "Failed to fetch service info", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"name":        name,
		"description": "Simple API with static data",
		"author":      "Your Name",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func geoHandler(w http.ResponseWriter, r *http.Request) {
	// Структуры для формирования ответа
	type CityResponse struct {
		Name       string `json:"name"`
		Population int    `json:"population"`
		Active     bool   `json:"active"`
	}
	type CountryResponse struct {
		ID     int            `json:"id"`
		Name   string         `json:"name"`
		Code2  string         `json:"code2"`
		Code3  string         `json:"code3"`
		Cities []CityResponse `json:"cities"`
	}

	// Запрос с объединением таблиц. Используем LEFT JOIN, чтобы включить страны без городов.
	query := `
		SELECT c.id, c.name, c.code2, c.code3, ci.name, ci.population, ci.active
		FROM countries c
		LEFT JOIN cities ci ON c.id = ci.country_id
		ORDER BY c.id
	`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Используем map для агрегации городов по стране
	countriesMap := make(map[int]*CountryResponse)
	for rows.Next() {
		var countryID int
		var countryName, code2, code3 string
		var cityName sql.NullString
		var cityPopulation sql.NullInt64
		var cityActive sql.NullBool

		if err := rows.Scan(&countryID, &countryName, &code2, &code3, &cityName, &cityPopulation, &cityActive); err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return
		}

		// Если страны еще нет в мапе, создаём новую запись
		country, exists := countriesMap[countryID]
		if !exists {
			country = &CountryResponse{
				ID:     countryID,
				Name:   countryName,
				Code2:  code2,
				Code3:  code3,
				Cities: []CityResponse{},
			}
			countriesMap[countryID] = country
		}

		// Если найдены данные о городе, добавляем город в срез
		if cityName.Valid {
			city := CityResponse{
				Name:       cityName.String,
				Population: int(cityPopulation.Int64),
				Active:     cityActive.Bool,
			}
			country.Cities = append(country.Cities, city)
		}
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Rows error", http.StatusInternalServerError)
		return
	}

	// Преобразуем map в срез для отправки JSON-ответа
	countries := make([]CountryResponse, 0, len(countriesMap))
	for _, country := range countriesMap {
		countries = append(countries, *country)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(countries)
}

func main() {
	initDB()
	http.HandleFunc("/", handler)
	http.HandleFunc("/api/status", statusHandler)
	http.HandleFunc("/api/info", infoHandler)
	http.HandleFunc("/api/geo", geoHandler)
	http.HandleFunc("/api/save", fetchAndSaveData)
	fmt.Println("Starting server on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func insertCountry(name, code2, code3 string) (int, error) {
	query := `INSERT INTO countries (name, code2, code3) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := db.QueryRow(query, name, code2, code3).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Функция для вставки столицы в таблицу
func insertCity(name string, countryID int, population int, active bool) error {
	query := `INSERT INTO cities (name, country_id, population, active) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(query, name, countryID, population, active)
	return err
}

// Функция для парсинга и сохранения данных
func fetchAndSaveData(w http.ResponseWriter, r *http.Request) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	// Отправляем GET-запрос на API
	resp, err := http.Get("https://restcountries.com/v3.1/all")
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to fetch data from API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to read API response", http.StatusInternalServerError)
		return
	}

	// Декодируем JSON в срез стран
	var countries []Country
	err = json.Unmarshal(body, &countries)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to parse API response", http.StatusInternalServerError)
		return
	}

	// Обрабатываем каждую страну
	for _, country := range countries {
		// Вставляем страну в таблицу
		countryName := country.Name.Common
		countryID, err := insertCountry(countryName, country.Alpha2Code, country.Alpha3Code)
		if err != nil {
			log.Printf("Failed to insert country %s: %v", countryName, err)
			continue
		}

		// Вставляем столицу в таблицу
		if len(country.Capital) > 0 {
			capitalName := country.Capital[0]
			population := country.Population
			log.Printf("Preparing to insert city: Name=%s, CountryID=%d, Population=%d", capitalName, countryID, population)
			err := insertCity(capitalName, countryID, population, true)
			if err != nil {
				log.Printf("Failed to insert city %s: %v", capitalName, err)
			} else {
				fmt.Printf("Inserted country %s and capital %s\n", countryName, capitalName)
			}
		}
	}

	// Ответ после успешной обработки
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Data saved successfully")
}
