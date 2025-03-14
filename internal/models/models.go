package models

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
