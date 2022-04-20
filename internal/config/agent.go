package config

import (
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rootAgentCmd = &cobra.Command{
		Use:   "agent",
		Short: "Metrics for Agent",
		Long:  "Metrics for Agent",
	}
	ReportInterval time.Duration
	PollInterval   time.Duration
)

const (
	defaultReportInterval = 10 * time.Second
	defaultPollInterval   = 2 * time.Second
)

type AgentConfig struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
	SignKey        string
}

type configAgentFromEVN struct {
	Address        string `env:"ADDRESS"`
	ReportInterval string `env:"REPORT_INTERVAL"`
	PollInterval   string `env:"POLL_INTERVAL"`
	SignKey        string `env:"KEY"`
}

func init() {
	rootAgentCmd.Flags().StringVarP(&Address, "address", "a", defaultServerAddress,
		"Pair of host:port to send data")

	rootAgentCmd.Flags().StringVarP(&SignKey, "key", "k", "",
		"Key for generate hash")

	rootAgentCmd.Flags().DurationVarP(&ReportInterval, "reportInterval", "r", defaultReportInterval,
		"Seconds to periodically save metrics")

	rootAgentCmd.Flags().DurationVarP(&PollInterval, "pollInterval", "p", defaultPollInterval,
		"Seconds to periodically send metrics to server")
}

func GetAgentConfig() AgentConfig {
	cfg := AgentConfig{}
	configEVN := configAgentFromEVN{}
	if err := env.Parse(&configEVN); err != nil {
		logrus.Error(err)
	}

	err := rootAgentCmd.Execute()
	if err != nil {
		logrus.Error(err)
	}

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

	if configEVN.ReportInterval == "" {
		cfg.ReportInterval = ReportInterval
	} else {
		var err error
		cfg.ReportInterval, err = time.ParseDuration(configEVN.ReportInterval)
		if err != nil {
			cfg.ReportInterval = defaultReportInterval
		}
	}

	if configEVN.PollInterval == "" {
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
