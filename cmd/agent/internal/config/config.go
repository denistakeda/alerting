package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval uint   `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func GetConfig() (Config, error) {
	config := Config{
		Address:        "localhost:8080",
		ReportInterval: 10,
		PollInterval:   2,
	}
	if err := env.Parse(&config); err != nil {
		return Config{}, errors.Wrap(err, "failed to parse agent configuration from the environment variables")
	}

	return config, nil
}
