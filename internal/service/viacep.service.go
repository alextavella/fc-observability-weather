package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alextavella/fc-observability-weather/internal/domain/model"
)

type ViaCepService struct {
}

func NewViaCepService() IAddressService {
	return &ViaCepService{}
}

func (s *ViaCepService) GetAddressByZipcode(ctx context.Context, zipcode string) (*model.ViaCepResult, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", zipcode)

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

	viacepResult := new(model.ViaCepResult)
	if err := json.NewDecoder(res.Body).Decode(viacepResult); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return viacepResult, nil
}
