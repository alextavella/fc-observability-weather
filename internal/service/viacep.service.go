package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alextavella/fc-observability-weather/internal/domain/model"
	"github.com/alextavella/fc-observability-weather/pkg/otel"
	"go.opentelemetry.io/otel/codes"
)

type viaCepService struct {
}

func NewViaCepService() IAddressService {
	return &viaCepService{}
}

func (s *viaCepService) GetAddressByZipcode(ctx context.Context, zipcode string) (*model.ViaCepResult, error) {
	tracer := otel.GetTracer()
	ctx, span := tracer.Start(ctx, "get-address")
	defer span.End()

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", zipcode)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		span.SetStatus(codes.Error, "Failure")
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		span.SetStatus(codes.Error, "Failure")
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		span.SetStatus(codes.Error, "Failure")
		return nil, fmt.Errorf("error response from server: %s", res.Status)
	}

	viacepResult := new(model.ViaCepResult)
	if err := json.NewDecoder(res.Body).Decode(viacepResult); err != nil {
		span.SetStatus(codes.Error, "Failure")
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	span.SetStatus(codes.Ok, "Success")
	return viacepResult, nil
}
