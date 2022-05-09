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
	Address            string
	Restore            bool
	StoreInterval      time.Duration
	StoreFile          string
	SignKey            string
	DataBaseConnection string
)

const (
	defaultServerAddress = "localhost:8080"
	defaultStoreFile     = "/tmp/devops-metrics-db.json"
	defaultStoreInterval = 300 * time.Second
)

type Config struct {
	Address            string
	StoreInterval      time.Duration
	StoreFile          string
	Restore            bool
	SignKey            string
	DataBaseConnection string
}

type configFromEVN struct {
	Address            string `env:"ADDRESS"`
	StoreInterval      string `env:"STORE_INTERVAL"`
	StoreFile          string `env:"STORE_FILE"`
	Restore            string `env:"RESTORE"`
	SignKey            string `env:"KEY"`
	DataBaseConnection string `env:"DATABASE_DSN"`
}

func init() {
	rootCmd.Flags().StringVarP(&Address, "address", "a", defaultServerAddress,
		"Pair of host:port to listen on")

	rootCmd.Flags().StringVarP(&SignKey, "key", "k", "",
		"Key for generate hash")

	rootCmd.Flags().StringVarP(&DataBaseConnection, "database-connection", "d", "",
		"Key for generate hash")

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
	logrus.Info("Config EVN : ", configEVN)
	err := rootCmd.Execute()
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info("Config Cmd: ", rootCmd)

	if configEVN.Address == "" {
		cfg.Address = Address
	} else {
		cfg.Address = configEVN.Address
	}

	if configEVN.SignKey == "" {
		cfg.SignKey = SignKey
	} else {
		cfg.SignKey = configEVN.SignKey
	}

	if configEVN.DataBaseConnection == "" {
		cfg.DataBaseConnection = DataBaseConnection
	} else {
		cfg.DataBaseConnection = configEVN.DataBaseConnection
	}

	if configEVN.StoreFile == "" {
		cfg.StoreFile = StoreFile
	} else {
		cfg.StoreFile = configEVN.StoreFile
	}

	if configEVN.Restore == "" {
		cfg.Restore = Restore
	} else {
		var err error
		cfg.Restore, err = strconv.ParseBool(configEVN.Restore)
		if err != nil {
			cfg.Restore = true
		}
	}

	if configEVN.StoreInterval == "" {
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
