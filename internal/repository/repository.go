package repository

import (
	"go_microservice/internal/db"
)

func InsertCountry(name, code2, code3 string) (int, error) {
	query := `INSERT INTO countries (name, code2, code3) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := db.DB.QueryRow(query, name, code2, code3).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func InsertCity(name string, countryID int, population int, active bool) error {
	query := `INSERT INTO cities (name, country_id, population, active) VALUES ($1, $2, $3, $4)`
	_, err := db.DB.Exec(query, name, countryID, population, active)
	return err
}
