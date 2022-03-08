package config

import (
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "agent",
		Short: "Metrics for Agent",
		Long:  "Metrics for Agent",
	}

	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
)

const (
	defaultServerAddress  = "localhost:8080"
	defaultReportInterval = 10 * time.Second
	defaultPollInterval   = 2 * time.Second
)

type Config struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

type configFromEVN struct {
	Address        string `env:"ADDRESS" envDefault:"notset"`
	ReportInterval string `env:"REPORT_INTERVAL" envDefault:"notset"`
	PollInterval   string `env:"POLL_INTERVAL" envDefault:"notset"`
}

func init() {
	rootCmd.Flags().StringVarP(&Address, "address", "a", defaultServerAddress,
		"Pair of host:port to send data")
	rootCmd.Flags().DurationVarP(&ReportInterval, "reportInterval", "r", defaultReportInterval,
		"Seconds to periodically save metrics")
	rootCmd.Flags().DurationVarP(&PollInterval, "pollInterval", "p", defaultPollInterval,
		"Seconds to periodically send metrics to server")
}

func GetConfig() Config {
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

	if configEVN.ReportInterval == "notset" {
		cfg.ReportInterval = ReportInterval
	} else {
		var err error
		cfg.ReportInterval, err = time.ParseDuration(configEVN.ReportInterval)
		if err != nil {
			cfg.ReportInterval = defaultReportInterval
		}
	}

	if configEVN.PollInterval == "notset" {
		cfg.PollInterval = PollInterval
	} else {
		var err error
		cfg.PollInterval, err = time.ParseDuration(configEVN.PollInterval)
		if err != nil {
			cfg.PollInterval = defaultPollInterval
		}
	}

	return cfg
}
