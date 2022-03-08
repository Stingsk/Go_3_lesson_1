package main

import (
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/Stingsk/Go_3_lesson_1/internal/file"
	"github.com/Stingsk/Go_3_lesson_1/internal/storage"
	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Stingsk/Go_3_lesson_1/internal/httputil"
	"github.com/Stingsk/Go_3_lesson_1/internal/logs"
)

const (
	defaultServerAddress = "localhost:8080"
	defaultStoreFile     = "/tmp/devops-metrics-db.json"
	defaultStoreInterval = 300 * time.Second
)

type config struct {
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

func main() {
	logs.Init()
	var metricData storage.MetricResourceMap
	cfg := getConfig()
	if cfg.Restore {
		logrus.Info("Load data from " + cfg.StoreFile)
		metricRead, _ := file.ReadMetrics(cfg.StoreFile)
		metricData.Metric = &metricRead
	} else {
		metricResource := make(map[string]storage.MetricResource)
		metricData.Metric = &metricResource
		logrus.Info("Load data is Off ")
	}
	os.Remove(cfg.StoreFile)
	logrus.Debug("Start server")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	var wg sync.WaitGroup
	wg.Add(1)
	go httputil.RunServer(&wg, sigChan, cfg.Address, &metricData, cfg.StoreFile, cfg.StoreInterval)

	wg.Wait()

	logrus.Debug("Shutdown server")
}

func getConfig() config {
	cfg := config{}
	configEVN := configFromEVN{}
	if err := env.Parse(&configEVN); err != nil {
		logrus.Error(err)
	}
	rootCmd := &cobra.Command{
		Use:   "server",
		Short: "Metrics for Server",
		Long:  "Metrics for Server",
	}

	if configEVN.Address == "notset" {
		rootCmd.Flags().StringVarP(&cfg.Address, "address", "a", defaultServerAddress,
			"Pair of host:port to listen on")
	} else {
		cfg.Address = configEVN.Address
	}

	if configEVN.StoreFile == "notset" {
		rootCmd.Flags().StringVarP(&cfg.StoreFile, "file", "f", defaultStoreFile,
			"Path to save metrics")
	} else {
		cfg.StoreFile = configEVN.StoreFile
	}

	if configEVN.Restore == "notset" {
		rootCmd.Flags().BoolVarP(&cfg.Restore, "restore", "r", true,
			"Flag to load initial metrics from storage ")
	} else {
		var err error
		cfg.Restore, err = strconv.ParseBool(configEVN.Restore)
		if err != nil {
			cfg.Restore = true
		}
	}

	if configEVN.StoreInterval == "notset" {
		rootCmd.Flags().DurationVarP(&cfg.StoreInterval, "interval", "i", defaultStoreInterval,
			"Seconds to periodically save metrics if 0 save immediately")
	} else {
		var err error
		cfg.StoreInterval, err = time.ParseDuration(configEVN.StoreInterval)
		if err != nil {
			cfg.StoreInterval = defaultStoreInterval
		}
	}

	return cfg
}
