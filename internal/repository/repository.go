package repository

import (
	"go_microservice/internal/db"
	"go_microservice/internal/models"
	"go_microservice/internal/migrations"
)

// GetServiceByID находит запись модели Service по её ID.
// Результат записывается в переданный указатель service.
func GetServiceByID(id int, service *models.Service) error {
	return db.DB.First(service, id).Error
}

// GetCountriesWithCities получает список стран с подгруженными городами.
// Результат записывается в переданный срез countries.
func GetCountriesWithCities(countries *[]models.CountryDB) error {
	return db.DB.Preload("Cities").Order("id").Find(countries).Error
}

func GetCountries(countries *[]models.CountryDB) error {
	return db.DB.Limit(5).Find(countries).Error
}

func InsertCountry(name, code2, code3 string) (int, error) {
	country := models.CountryDB{
		Name:  name,
		Code2: code2,
		Code3: code3,
	}
	// Создаем запись. Если нужна проверка на существование, можно использовать FirstOrCreate.
	if err := db.DB.Create(&country).Error; err != nil {
		return 0, err
	}
	return int(country.ID), nil
}

func InsertCity(name string, countryID int, population int, active bool) error {
	city := models.CityDB{
		Name:       name,
		CountryID:  uint(countryID),
		Population: population,
		Active:     active,
	}
	return db.DB.Create(&city).Error
}

func CreateTradeRow(ticker *migrations.BinanceTicker) error {
	return db.DB.Create(ticker).Error
}

func GetTradeRows(tickers *[]migrations.BinanceTicker, limit int) error {
	subQuery := db.DB.Model(&migrations.BinanceTicker{}).
		Select("*").
		Order("created_at desc").
		Limit(limit)
	return db.DB.Table("(?) as sub", subQuery).
		Order("created_at asc").
		Find(tickers).Error
}
