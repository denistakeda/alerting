package servercfg

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

type Config struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	Key           string        `env:"KEY"`
	DatabaseDSN   string        `env:"DATABASE_DSN"`
}

func GetConfig() (Config, error) {
	config := Config{}

	// Get flags
	flag.StringVar(&config.Address, "a", "localhost:8080", "Where to start server")
	flag.BoolVar(&config.Restore, "r", true, "Restore from the file")
	flag.DurationVar(&config.StoreInterval, "i", 300*time.Second, "Interval to dump state")
	flag.StringVar(&config.StoreFile, "f", "/tmp/devops-metrics-db.json", "Database file")
	flag.StringVar(&config.Key, "k", "", "Hash key")
	flag.StringVar(&config.DatabaseDSN, "d", "", "Database DSN")
	flag.Parse()

	// Populate data from the env variables
	if err := env.Parse(&config); err != nil {
		return Config{}, errors.Wrap(err, "failed to parse server configuration from the environment variables")
	}

	return config, nil
}
