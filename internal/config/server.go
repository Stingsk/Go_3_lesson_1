package config

import (
	"strconv"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "server",
		Short: "Metrics for Server",
		Long:  "Metrics for Server",
	}
	Address       string
	Restore       bool
	StoreInterval time.Duration
	StoreFile     string
)

const (
	defaultServerAddress = "localhost:8080"
	defaultStoreFile     = "/tmp/devops-metrics-db.json"
	defaultStoreInterval = 300 * time.Second
)

type Config struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
}

type configFromEVN struct {
	Address       string `env:"ADDRESS" envDefault:"notset"`
	StoreInterval string `env:"STORE_INTERVAL" envDefault:"notset"`
	StoreFile     string `env:"STORE_FILE" envDefault:"notset"`
	Restore       string `env:"RESTORE" envDefault:"notset"`
}

func init() {
	rootCmd.Flags().StringVarP(&Address, "address", "a", defaultServerAddress,
		"Pair of host:port to listen on")

	rootCmd.Flags().BoolVarP(&Restore, "restore", "r", true,
		"Flag to load initial metrics from storage ")

	rootCmd.Flags().StringVarP(&StoreFile, "file", "f", defaultStoreFile,
		"Path to save metrics")

	rootCmd.Flags().DurationVarP(&StoreInterval, "interval", "i", defaultStoreInterval,
		"Seconds to periodically save metrics if 0 save immediately")
}

func GetServerConfig() Config {
	cfg := Config{}
	configEVN := configFromEVN{}
	if err := env.Parse(&configEVN); err != nil {
		logrus.Error(err)
	}
	err := rootCmd.Execute()
	if err != nil {
		logrus.Error(err)
	}

	if configEVN.Address == "notset" {
		cfg.Address = Address
	} else {
		cfg.Address = configEVN.Address
	}

	if configEVN.StoreFile == "notset" {
		cfg.StoreFile = StoreFile
	} else {
		cfg.StoreFile = configEVN.StoreFile
	}

	if configEVN.Restore == "notset" {
		cfg.Restore = Restore
	} else {
		var err error
		cfg.Restore, err = strconv.ParseBool(configEVN.Restore)
		if err != nil {
			cfg.Restore = true
		}
	}

	if configEVN.StoreInterval == "notset" {
		cfg.StoreInterval = StoreInterval
	} else {
		var err error
		cfg.StoreInterval, err = time.ParseDuration(configEVN.StoreInterval)
		if err != nil {
			cfg.StoreInterval = defaultStoreInterval
		}
	}

	return cfg
}
