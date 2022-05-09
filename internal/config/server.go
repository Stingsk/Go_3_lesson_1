package config

import (
	"fmt"
	"regexp"
	"time"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "server",
		Short: "Metrics for Server",
		Long:  "Metrics for Server",
		RunE: func(cmd *cobra.Command, args []string) error {
			re := regexp.MustCompile(`(DEBUG|INFO|WARNING|ERROR)`)

			if !re.MatchString(LogLevel) {
				return fmt.Errorf("%w: --log-level", "invalid param specified")
			}

			return nil
		},
	}
	Address            string
	Restore            bool
	StoreInterval      time.Duration
	StoreFile          string
	SignKey            string
	DataBaseConnection string
	LogLevel           string
)

const (
	defaultServerAddress = "localhost:8080"
	defaultStoreFile     = "/tmp/devops-metrics-db.json"
	defaultStoreInterval = 300 * time.Second
)

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
	rootCmd.Flags().StringVarP(&LogLevel, "log-level", "l", "ERROR",
		"Set log level: DEBUG|INFO|WARNING|ERROR")
}

func GetServerConfig() error {
	return rootCmd.Execute()
}
