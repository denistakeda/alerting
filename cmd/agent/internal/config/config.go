package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type Config struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}

func GetConfig() (Config, error) {
	config := Config{
		Address:        "localhost:8080",
		ReportInterval: 10 * time.Second,
		PollInterval:   2 * time.Second,
	}
	if err := env.Parse(&config); err != nil {
		return Config{}, errors.Wrap(err, "failed to parse agent configuration from the environment variables")
	}

	if !strings.HasPrefix(config.Address, "http") {
		config.Address = fmt.Sprintf("http://%s", config.Address)
	}

	return config, nil
}
