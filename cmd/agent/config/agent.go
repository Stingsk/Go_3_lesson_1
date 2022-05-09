package config

import (
	"time"

	"github.com/spf13/cobra"
)

var (
	rootAgentCmd = &cobra.Command{
		Use:   "agent",
		Short: "Metrics for Agent",
		Long:  "Metrics for Agent",
	}
	Address        string
	PollInterval   time.Duration
	ReportInterval time.Duration
	ServerTimeout  time.Duration
	SignKey        string
	LogLevel       string
)

const (
	defaultServerAddress  = "localhost:8080"
	defaultPollInterval   = 2 * time.Second
	defaultReportInterval = 10 * time.Second
	defaultServerTimeout  = 1 * time.Second
)

func init() {
	rootAgentCmd.Flags().StringVarP(&Address, "address", "a", defaultServerAddress,
		"Pair of host:port to send data")

	rootAgentCmd.Flags().DurationVarP(&ServerTimeout, "timeout", "t", defaultServerTimeout,
		"Timeout for server connection")

	rootAgentCmd.Flags().StringVarP(&SignKey, "key", "k", "",
		"Key for generate hash")

	rootAgentCmd.Flags().DurationVarP(&ReportInterval, "reportInterval", "r", defaultReportInterval,
		"Seconds to periodically save metrics")

	rootAgentCmd.Flags().DurationVarP(&PollInterval, "pollInterval", "p", defaultPollInterval,
		"Seconds to periodically send metrics to server")

	rootAgentCmd.Flags().StringVarP(&LogLevel, "log-level", "l", "INFO",
		"Set log level: DEBUG|INFO|WARNING|ERROR")
}

func GetAgentConfig() error {
	return rootAgentCmd.Execute()
}
