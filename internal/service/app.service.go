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

type appService struct {
	httpClient *http.Client
	host       string
}

func NewAppService(httpClient *http.Client, host string) IAppService {
	return &appService{
		httpClient: httpClient,
		host:       host,
	}
}

func (s *appService) GetWeather(ctx context.Context, zipcode string) (*model.AppResult, error) {
	tracer := otel.GetTracer()
	ctx, span := tracer.Start(ctx, "get-weather")
	defer span.End()

	url := fmt.Sprintf("%s/weather/%s/zipcode", s.host, zipcode)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		span.SetStatus(codes.Error, "Failure")
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		span.SetStatus(codes.Error, "Failure")
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		span.SetStatus(codes.Error, "Failure")
		return nil, fmt.Errorf("error response from server: %s", res.Status)
	}

	result := &model.AppResult{}
	if err := json.NewDecoder(res.Body).Decode(result); err != nil {
		span.SetStatus(codes.Error, "Failure")
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	span.SetStatus(codes.Ok, "Success")
	return result, nil
}
