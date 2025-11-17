package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MitrofanK/Test-Menti/internal/client"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *Redis {
	return &Redis{
		client: client,
	}
}

func (r *Redis) SetRates(ctx context.Context, rates map[string]client.Currency, ttl time.Duration) error {
	rateBytes, err := json.Marshal(rates)
	if err != nil {
		return fmt.Errorf("failed to marshal rates for redis: %w", err)
	}

	key := "currency_rates"
	err = r.client.Set(ctx, key, rateBytes, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set rates in redis: %w", err)
	}
	return nil
}

func (r *Redis) GetRates(ctx context.Context) (map[string]client.Currency, error) {
	key := "currency_rates"

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, redis.Nil
		}
		return nil, fmt.Errorf("failed to get rates from redis: %w", err)
	}

	var rates map[string]client.Currency
	if err := json.Unmarshal([]byte(val), &rates); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rates from redis: %w", err)
	}

	return rates, nil
}
