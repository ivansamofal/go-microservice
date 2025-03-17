package models

type CountryDB struct {
	ID         uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string   `gorm:"size:255" json:"name"`
	Code2      string   `gorm:"size:10" json:"code2"`
	Code3      string   `gorm:"size:10" json:"code3"`
	Population int      `json:"population"`
	Cities     []CityDB `gorm:"foreignKey:CountryID" json:"cities"`
}

func (CountryDB) TableName() string {
	return "countries"
}

type CityDB struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string `gorm:"size:255" json:"name"`
	Population int    `json:"population"`
	Active     bool   `json:"active"`
	CountryID  uint   `json:"country_id"`
}

func (CityDB) TableName() string {
	return "cities"
}

type Service struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"size:255" json:"name"`
}

func (Service) TableName() string {
	return "services"
}
