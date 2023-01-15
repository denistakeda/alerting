package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
	"time"
)

type Config struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
}

func GetConfig() (Config, error) {
	config := Config{
		Address:       "localhost:8080",
		StoreInterval: 300 * time.Second,
		StoreFile:     "/tmp/devops-metrics-db.json",
		Restore:       true,
	}
	if err := env.Parse(&config); err != nil {
		return Config{}, errors.Wrap(err, "failed to parse server configuration from the environment variables")
	}

	return config, nil
}
