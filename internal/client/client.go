package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CurrencyResponse struct {
	Date      time.Time           `json:"Date"`
	Timestamp time.Time           `json:"Timestamp"`
	Valute    map[string]Currency `json:"Valute"`
}

type Currency struct {
	CharCode string  `json:"CharCode"`
	Nominal  int     `json:"Nominal"`
	Name     string  `json:"Name"`
	Value    float64 `json:"Value"`
}

type Client struct {
	httpClient *http.Client
	url        string
}

func NewClient(url string, timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		url: url,
	}
}

func (c *Client) GetRates(ctx context.Context) (*CurrencyResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var currencyResponse CurrencyResponse

	if err := json.Unmarshal(body, &currencyResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &currencyResponse, nil
}
