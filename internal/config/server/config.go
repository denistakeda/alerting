package servercfg

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

// Config is a configuration for alerting server
type Config struct {
	Config        string        `env:"CONFIG"`
	Address       string        `env:"ADDRESS" json:"address"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" json:"store_interval"`
	StoreFile     string        `env:"STORE_FILE" json:"store_file"`
	Restore       bool          `env:"RESTORE" json:"restore"`
	Key           string        `env:"KEY" json:"key"`
	DatabaseDSN   string        `env:"DATABASE_DSN" json:"database_dsn"`
	Certificate   string        `env:"CERTIFICATE" json:"certificate"`
	CryptoKey     string        `env:"CRYPTO_KEY" json:"crypto_key"`
	TrustedSubnet string        `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
}

// GetConfig extracts the configuration from environment variables and flags
func GetConfig() (Config, error) {
	// Set default values
	config := Config{
		Address:       "localhost:8080",
		StoreInterval: 300 * time.Second,
		StoreFile:     "/tmp/devops-metrics-db.json",
		Restore:       true,
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
	flag.StringVar(&config.Address, "a", config.Address, "Where to start server")
	flag.BoolVar(&config.Restore, "r", config.Restore, "Restore from the file")
	flag.DurationVar(&config.StoreInterval, "i", config.StoreInterval, "Interval to dump state")
	flag.StringVar(&config.StoreFile, "f", config.StoreFile, "Database file")
	flag.StringVar(&config.Key, "k", config.Key, "Hash key")
	flag.StringVar(&config.DatabaseDSN, "d", config.DatabaseDSN, "Database DSN")
	flag.StringVar(&config.Certificate, "certificate", config.Certificate, "Path to a file with a certificate")
	flag.StringVar(&config.CryptoKey, "crypto-key", config.CryptoKey, "Path to a file with a private key")
	flag.StringVar(&config.TrustedSubnet, "t", config.TrustedSubnet, "Trusted subnet")
	flag.Parse()

	// Populate data from the env variables
	if err := env.Parse(&config); err != nil {
		return Config{}, errors.Wrap(err, "failed to parse server configuration from the environment variables")
	}

	return config, nil
}
