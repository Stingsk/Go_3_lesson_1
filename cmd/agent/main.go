package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Stingsk/Go_3_lesson_1/internal/httputil"
	"github.com/Stingsk/Go_3_lesson_1/internal/logs"
	"github.com/Stingsk/Go_3_lesson_1/internal/metrics"
)

const (
	defaultServerAddress  = "localhost:8080"
	defaultReportInterval = 10 * time.Second
	defaultPollInterval   = 2 * time.Second
)

type config struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

type configFromEVN struct {
	Address        string `env:"ADDRESS" envDefault:"notset"`
	ReportInterval string `env:"REPORT_INTERVAL" envDefault:"notset"`
	PollInterval   string `env:"POLL_INTERVAL" envDefault:"notset"`
}

func main() {
	logs.Init()
	cfg := getConfig()
	logrus.Debug("Start agent")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	wg := &sync.WaitGroup{}
	var sensorData metrics.SensorData

	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go metrics.RunGetMetrics(ctx, cfg.PollInterval, &sensorData, wg)

	wg.Add(1)
	go httputil.RunSender(ctx, cfg.ReportInterval, &sensorData, wg, cfg.Address)

	go func() {
		<-sigChan
		cancel()
	}()

	wg.Wait()
	logrus.Debug("Shutdown agent")
}

func getConfig() config {
	cfg := config{}
	configEVN := configFromEVN{}
	if err := env.Parse(&configEVN); err != nil {
		logrus.Error(err)
	}
	rootCmd := &cobra.Command{
		Use:   "agent",
		Short: "Metrics for Agent",
		Long:  "Metrics for Agent",
	}

	if configEVN.Address == "notset" {
		rootCmd.Flags().StringVarP(&cfg.Address, "address", "a", defaultServerAddress,
			"Pair of host:port to send data")
	} else {
		cfg.Address = configEVN.Address
	}

	if configEVN.ReportInterval == "notset" {
		rootCmd.Flags().DurationVarP(&cfg.ReportInterval, "reportInterval", "r", defaultReportInterval,
			"Seconds to periodically save metrics")
	} else {
		var err error
		cfg.ReportInterval, err = time.ParseDuration(configEVN.ReportInterval)
		if err != nil {
			cfg.ReportInterval = defaultReportInterval
		}
	}

	if configEVN.PollInterval == "notset" {
		rootCmd.Flags().DurationVarP(&cfg.PollInterval, "pollInterval", "p", defaultPollInterval,
			"Seconds to periodically send metrics to server")
	} else {
		var err error
		cfg.PollInterval, err = time.ParseDuration(configEVN.PollInterval)
		if err != nil {
			cfg.PollInterval = defaultPollInterval
		}
	}

	return cfg
}
