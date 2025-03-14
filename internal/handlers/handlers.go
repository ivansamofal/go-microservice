package handlers

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"go_microservice/internal/db"
	"go_microservice/internal/models"
	"go_microservice/internal/repository"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, Go Microservice!")
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "running",
		"version": "1.0.0",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	var name string
	err := db.DB.QueryRow("SELECT name FROM services WHERE id = 1").Scan(&name)
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

func GeoHandler(w http.ResponseWriter, r *http.Request) {
	// Используем структуры из пакета models
	type CityResponse = models.CityResponse
	type CountryResponse = models.CountryResponse

	query := `
        SELECT c.id, c.name, c.code2, c.code3, ci.name, ci.population, ci.active
        FROM countries c
        LEFT JOIN cities ci ON c.id = ci.country_id
        ORDER BY c.id
    `
	rows, err := db.DB.Query(query)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

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

	countries := make([]CountryResponse, 0, len(countriesMap))
	for _, country := range countriesMap {
		countries = append(countries, *country)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(countries)
}

func FetchAndSaveData(w http.ResponseWriter, r *http.Request) {
	// Настройка TLS для пропуска проверки сертификата
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get("https://restcountries.com/v3.1/all")
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to fetch data from API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to read API response", http.StatusInternalServerError)
		return
	}

	var countries []models.Country
	err = json.Unmarshal(body, &countries)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to parse API response", http.StatusInternalServerError)
		return
	}

	for _, country := range countries {
		countryName := country.Name.Common
		countryID, err := repository.InsertCountry(countryName, country.Alpha2Code, country.Alpha3Code)
		if err != nil {
			log.Printf("Failed to insert country %s: %v", countryName, err)
			continue
		}

		if len(country.Capital) > 0 {
			capitalName := country.Capital[0]
			population := country.Population
			log.Printf("Preparing to insert city: Name=%s, CountryID=%d, Population=%d", capitalName, countryID, population)
			err := repository.InsertCity(capitalName, countryID, population, true)
			if err != nil {
				log.Printf("Failed to insert city %s: %v", capitalName, err)
			} else {
				fmt.Printf("Inserted country %s and capital %s\n", countryName, capitalName)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Data saved successfully")
}
