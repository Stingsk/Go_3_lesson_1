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

func init() {
	rootAgentCmd.Flags().StringVarP(&Address, "address", "a", defaultServerAddress,
		"Pair of host:port to send data")

	rootAgentCmd.Flags().StringVarP(&SignKey, "key", "k", "",
		"Key for generate hash")

	rootAgentCmd.Flags().DurationVarP(&ReportInterval, "reportInterval", "r", defaultReportInterval,
		"Seconds to periodically save metrics")

	rootAgentCmd.Flags().DurationVarP(&PollInterval, "pollInterval", "p", defaultPollInterval,
		"Seconds to periodically send metrics to server")

	rootAgentCmd.Flags().StringVarP(&LogLevel, "log-level", "l", "ERROR",
		"Set log level: DEBUG|INFO|WARNING|ERROR")
}

func GetAgentConfig() error {
	return rootAgentCmd.Execute()
}
