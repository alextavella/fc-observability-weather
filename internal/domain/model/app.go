package model

type AppResult struct {
	Zipcode    string  `json:"zipcode"`
	Celsius    float64 `json:"celsius"`
	Fahrenheit float64 `json:"fahrenheit"`
	Kelvin     float64 `json:"kelvin"`
}
