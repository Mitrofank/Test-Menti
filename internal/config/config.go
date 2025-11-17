package config

import (
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Postgres PostgresConfig `koanf:"postgres"`
	Redis    RedisConfig    `koanf:"redis"`
	HTTP     HTTPConfig     `koanf:"http"`
	JWT      JWTConfig      `koanf:"jwt"`
	Currency CurrencyConfig `koanf:"currency"`
}

type PostgresConfig struct {
	URL string `koanf:"url"`
}

type RedisConfig struct {
	Host     string `koanf:"host"`
	Port     string `koanf:"port"`
	Password string `koanf:"password"`
}

type HTTPConfig struct {
	Port string `koanf:"port"`
}

type CurrencyConfig struct {
	URL     string        `koanf:"url"`
	Timeout time.Duration `koanf:"timeout"`
	TTL     time.Duration `koanf:"ttl"`
}

type JWTConfig struct {
	SigningKey string        `koanf:"signing_key"`
	TokenTTL   time.Duration `koanf:"ttl"`
}

func LoadConfig() (*Config, error) {
	k := koanf.New(".")

	if err := k.Load(file.Provider(".env"), dotenv.Parser()); err != nil {
		log.Error("No .env file found, relying on environment variables.")
	}

	callback := func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "APP_")), "_", ".", 1)
	}

	if err := k.Load(env.Provider("APP_", ".", callback), nil); err != nil {
		return nil, err
	}

	var cfg Config

	if err := k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{Tag: "koanf"}); err != nil {
		return nil, err
	}

	return &cfg, nil
}
