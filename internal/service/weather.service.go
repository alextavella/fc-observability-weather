package service

import (
	"context"
	"fmt"
	"net/http"
)

type weatherService struct {
	host string
}

func NewWeatherService(host string) IWeatherService {
	return &weatherService{
		host: host,
	}
}

func (s *weatherService) GetWeatherByZipCode(ctx context.Context, zipcode string) (any, error) {
	url := fmt.Sprintf("%s/weather/%s/zipcode", s.host, zipcode)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from server: %s", res.Status)
	}

	return res, nil
}
