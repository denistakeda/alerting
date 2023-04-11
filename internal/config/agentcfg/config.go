package agentcfg

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

// Config is a configuration for agent
type Config struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	Key            string        `env:"KEY"`
	RateLimit      int           `env:"RATE_LIMIT"`
}

// GetConfig extracts the configuration from environment variables and flags
func GetConfig() (Config, error) {
	config := Config{}

	// Get flags
	flag.StringVar(&config.Address, "a", "http://localhost:8080", "Server to send metrics to")
	flag.DurationVar(&config.ReportInterval, "r", 10*time.Second, "Interval to send metrics to server")
	flag.DurationVar(&config.PollInterval, "p", 2*time.Second, "Interval to collect metrics")
	flag.StringVar(&config.Key, "k", "", "Key to sign")
	flag.IntVar(&config.RateLimit, "l", 1, "The maximum amount of active requests")
	flag.Parse()

	// Populate data from the env variables
	if err := env.Parse(&config); err != nil {
		return Config{}, errors.Wrap(err, "failed to parse agent configuration from the environment variables")
	}

	// Validation and post-processing
	if !strings.HasPrefix(config.Address, "http") {
		config.Address = fmt.Sprintf("http://%s", config.Address)
	}

	return config, nil
}
