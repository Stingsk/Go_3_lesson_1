package config

import (
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/spf13/cobra"
)

var (
	rootAgentCmd = &cobra.Command{
		Use:   "agent",
		Short: "Metrics for Agent",
		Long:  "Metrics for Agent",
	}
	address        string
	pollInterval   time.Duration
	reportInterval time.Duration
	serverTimeout  time.Duration
	signKey        string
	logLevel       string
)

const (
	defaultServerAddress  = "localhost:8080"
	defaultPollInterval   = 2 * time.Second
	defaultReportInterval = 10 * time.Second
	defaultServerTimeout  = 1 * time.Second
)

type Config struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	SignKey        string        `env:"KEY"`
	LogLevel       string        `env:"LOG_LEVEL"`
}

func setConfig() {
	rootAgentCmd.Flags().StringVarP(&address, "address", "a", defaultServerAddress,
		"Pair of host:port to send data")

	rootAgentCmd.Flags().DurationVarP(&serverTimeout, "timeout", "t", defaultServerTimeout,
		"Timeout for server connection")

	rootAgentCmd.Flags().StringVarP(&signKey, "key", "k", "",
		"Key for generate hash")

	rootAgentCmd.Flags().DurationVarP(&reportInterval, "reportInterval", "r", defaultReportInterval,
		"Seconds to periodically save metrics")

	rootAgentCmd.Flags().DurationVarP(&pollInterval, "pollInterval", "p", defaultPollInterval,
		"Seconds to periodically send metrics to server")

	rootAgentCmd.Flags().StringVarP(&logLevel, "log-level", "l", "INFO",
		"Set log level: DEBUG|INFO|WARNING|ERROR")
}

func GetAgentConfig() (Config, error) {
	setConfig()
	rootAgentCmd.Execute()

	var agentConfig = Config{
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
		Address:        address,
		SignKey:        signKey,
		LogLevel:       logLevel,
	}

	if err := env.Parse(&agentConfig); err != nil {
		return Config{}, err
	}

	return agentConfig, nil
}
