package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"go-clean-template/config"
	"go-clean-template/internal/domain"

	"net/http"
)

type dataAPI struct {
	client *http.Client
	url    string
	path   string
}

func NewDataAPI(client *http.Client, cfg config.API) *dataAPI {
	return &dataAPI{
		client,
		cfg.URL,
		cfg.Path,
	}
}

func (a *dataAPI) GetData(ctx context.Context,
	req domain.ServiceRequest) ([]struct{}, error) {
	const op string = "dataAPI.GetData"
	data, err := a.getDataFromAPI(ctx, a.path, req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return data, nil
}

func (a *dataAPI) getDataFromAPI(ctx context.Context, path string, req domain.ServiceRequest) ([]struct{}, error) {
	reqURL := fmt.Sprintf("%s%s", a.url, path)
	_ = req

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http.Do: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("code = %d; status = %s", resp.StatusCode, resp.Status)
	}

	var data []struct{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("json.NewDecoder: %w", err)
	}

	return data, nil
}
