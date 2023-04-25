package agentcfg

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

// Config is a configuration for agent
type Config struct {
	Config         string        `env:"CONFIG"`
	Address        string        `env:"ADDRESS" json:"address"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" json:"report_interval"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" json:"poll_interval"`
	Key            string        `env:"KEY" json:"key"`
	RateLimit      int           `env:"RATE_LIMIT" json:"rate_limit"`
	CryptoKey      string        `env:"CRYPTO_KEY" json:"crypto_key"`
}

// GetConfig extracts the configuration from environment variables and flags
func GetConfig() (Config, error) {
	// Set default values
	config := Config{
		Address:        "http://localhost:8080",
		ReportInterval: 10 * time.Second,
		PollInterval:   2 * time.Second,
		RateLimit:      1,
	}

	// Read from file
	flag.StringVar(&config.Config, "c", "", "Path to configuration file")
	flag.Parse()

	if config.Config != "" {
		content, err := os.ReadFile(config.Config)
		if err != nil {
			log.Fatal("Error when opening file: ", err)
		}

		// Now let's unmarshall the data into `payload`
		err = json.Unmarshal(content, &config)
		if err != nil {
			log.Fatal("Error during Unmarshal(): ", err)
		}
	}

	// Get flags
	flag.StringVar(&config.Address, "a", config.Address, "Server to send metrics to")
	flag.DurationVar(&config.ReportInterval, "r", config.ReportInterval, "Interval to send metrics to server")
	flag.DurationVar(&config.PollInterval, "p", config.PollInterval, "Interval to collect metrics")
	flag.StringVar(&config.Key, "k", config.Key, "Key to sign")
	flag.IntVar(&config.RateLimit, "l", config.RateLimit, "The maximum amount of active requests")
	flag.StringVar(&config.CryptoKey, "c", config.CryptoKey, "Path to the certificate")
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
