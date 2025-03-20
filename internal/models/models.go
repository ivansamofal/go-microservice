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

type Credentials struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"password"`
}

type GDPData struct {
	Country             string  `json:"country"`
	Year                int     `json:"year"`
	GDPGrowth           float64 `json:"gdp_growth"`
	GDPNominal          float64 `json:"gdp_nominal"`
	GDPPerCapitaNominal float64 `json:"gdp_per_capita_nominal"`
	GDPPPP              float64 `json:"gdp_ppp"`
	GDPPerCapitaPPP     float64 `json:"gdp_per_capita_ppp"`
	GDPPPPShare         float64 `json:"gdp_ppp_share"`
}

type CountryGDPResult struct {
	Country           string  `json:"country"`
	AverageGDPNominal float64 `json:"average_gdp_nominal"`
}

type TickerData struct {
	Timestamp       int64   `json:"timestamp,string"`
	Open            float64 `json:"open,string"`
	High            float64 `json:"high,string"`
	Low             float64 `json:"low,string"`
	Last            float64 `json:"last,string"`
	Volume          float64 `json:"volume,string"`
	VWAP            float64 `json:"vwap,string"`
	Bid             float64 `json:"bid,string"`
	Ask             float64 `json:"ask,string"`
	Side            int     `json:"side,string"`
	Open24          float64 `json:"open_24,string"`
	PercentChange24 float64 `json:"percent_change_24,string"`
}
