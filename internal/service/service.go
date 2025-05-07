package service

import (
	"context"

	"github.com/alextavella/fc-observability-weather/internal/domain/model"
)

type IAddressService interface {
	GetAddressByZipcode(ctx context.Context, zipcode string) (*model.ViaCepResult, error)
}

type IWeatherService interface {
	GetWeatherByZipCode(ctx context.Context, zipcode string) (any, error)
}
