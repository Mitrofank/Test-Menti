package currency

import (
	"context"
	"fmt"
	"time"

	client_model "github.com/MitrofanK/Test-Menti/internal/client"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

type CurrencyClient interface {
	GetRates(ctx context.Context) (*client_model.CurrencyResponse, error)
}

type CurrencyCache interface {
	GetRates(ctx context.Context) (map[string]client_model.Currency, error)
	SetRates(ctx context.Context, rates map[string]client_model.Currency, ttl time.Duration) error
}

type Service struct {
	client   CurrencyClient
	cache    CurrencyCache
	cacheTTL time.Duration
}

func NewService(client CurrencyClient, cache CurrencyCache, ttl time.Duration) *Service {
	return &Service{
		client:   client,
		cache:    cache,
		cacheTTL: ttl,
	}
}

func (s *Service) GetRates(ctx context.Context) (map[string]client_model.Currency, error) {
	cachedRates, err := s.cache.GetRates(ctx)
	if err == nil {
		log.Info("Cache hit!")
		return cachedRates, nil
	}

	if err != redis.Nil {
		return nil, fmt.Errorf("error getting rates from cache: %w", err)
	}

	log.Info("Cache miss! Fetching from external API.")

	response, err := s.client.GetRates(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get rates from client: %w", err)
	}

	processedRates := make(map[string]client_model.Currency)
	for code, currency := range response.Valute {
		if currency.Nominal > 0 {
			currency.Value = currency.Value / float64(currency.Nominal)
			currency.Nominal = 1
			processedRates[code] = currency
		}
	}

	go func() {
		if err := s.cache.SetRates(context.Background(), processedRates, s.cacheTTL); err != nil {
			log.Printf("failed to set rates in cache: %v", err)
		}
	}()

	return processedRates, nil

}

func (s *Service) Convert(ctx context.Context, from, to string, amount float64) (float64, error) {
	rates, err := s.GetRates(ctx)
	if err != nil {
		return 0.0, err
	}

	rates["RUB"] = client_model.Currency{Value: 1.0, Nominal: 1}

	fromRate, ok := rates[from]
	if !ok {
		return 0.0, fmt.Errorf("currency %s not found", from)
	}

	toRate, ok := rates[to]
	if !ok {
		return 0.0, fmt.Errorf("currency %s not found", to)
	}

	result := (amount * fromRate.Value) / toRate.Value

	return result, nil
}
