package util

func ConvertToFahrenheit(celsius float64) float64 {
	return (celsius * 9 / 5) + 32
}

func ConvertToKelvin(celsius float64) float64 {
	return celsius + 273.15
}
