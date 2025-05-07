package service

import (
	"context"

	"github.com/alextavella/fc-observability-weather/internal/domain/model"
)

type IAppService interface {
	GetWeather(ctx context.Context, zipcode string) (*model.AppResult, error)
}

type IAddressService interface {
	GetAddressByZipcode(ctx context.Context, zipcode string) (*model.ViaCepResult, error)
}

type IWeatherService interface {
	GetWeatherByZipCode(ctx context.Context, zipcode string) (*model.WeatherResult, error)
}
