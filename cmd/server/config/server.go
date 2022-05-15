package config

import (
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "server",
		Short: "Metrics for Server",
		Long:  "Metrics for Server",
	}
	address            string
	restore            bool
	storeInterval      time.Duration
	storeFile          string
	signKey            string
	dataBaseConnection string
	logLevel           string
)

type Config struct {
	Address            string        `env:"ADDRESS"`
	StoreInterval      time.Duration `env:"STORE_INTERVAL"`
	StoreFile          string        `env:"STORE_FILE"`
	Restore            bool          `env:"RESTORE"`
	SignKey            string        `env:"KEY"`
	DataBaseConnection string        `env:"DATABASE_DSN"`
	LogLevel           string        `env:"LOG_LEVEL"`
}

const (
	defaultServerAddress = "localhost:8080"
	defaultStoreFile     = "/tmp/devops-metrics-db.json"
	defaultStoreInterval = 300 * time.Second
)

func setConfig() {
	rootCmd.Flags().StringVarP(&address, "address", "a", defaultServerAddress,
		"Pair of host:port to listen on")

	rootCmd.Flags().StringVarP(&signKey, "key", "k", "",
		"Key for generate hash")

	rootCmd.Flags().StringVarP(&dataBaseConnection, "database-connection", "d", "",
		"Key for generate hash")

	rootCmd.Flags().BoolVarP(&restore, "restore", "r", true,
		"Flag to load initial metrics from storage ")

	rootCmd.Flags().StringVarP(&storeFile, "file", "f", defaultStoreFile,
		"Path to save metrics")

	rootCmd.Flags().DurationVarP(&storeInterval, "interval", "i", defaultStoreInterval,
		"Seconds to periodically save metrics if 0 save immediately")
	rootCmd.Flags().StringVarP(&logLevel, "log-level", "l", "INFO",
		"Set log level: DEBUG|INFO|WARNING|ERROR")
}

func GetServerConfig() (Config, error) {
	setConfig()

	if err := rootCmd.Execute(); err != nil {
		return Config{}, err
	}

	var serverConfig = Config{
		Address:            address,
		Restore:            restore,
		StoreFile:          storeFile,
		StoreInterval:      storeInterval,
		SignKey:            signKey,
		DataBaseConnection: dataBaseConnection, // "postgresql://localhost:5432/metrics",
		LogLevel:           logLevel,
	}

	if err := env.Parse(&serverConfig); err != nil {
		return Config{}, err
	}
	if serverConfig.DataBaseConnection == "" {
		serverConfig.DataBaseConnection = serverConfig.DataBaseConnection
	}

	return serverConfig, nil
}
