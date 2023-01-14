package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

type Config struct {
	Address string `env:"ADDRESS"`
}

func GetConfig() (Config, error) {
	config := Config{
		Address: "localhost:8080",
	}
	if err := env.Parse(&config); err != nil {
		return Config{}, errors.Wrap(err, "failed to parse server configuration from the environment variables")
	}

	return config, nil
}
