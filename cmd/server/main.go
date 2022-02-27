package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Stingsk/Go_3_lesson_1/internal/file"
	"github.com/Stingsk/Go_3_lesson_1/internal/storage"
	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"

	"github.com/Stingsk/Go_3_lesson_1/internal/httputil"
	"github.com/Stingsk/Go_3_lesson_1/internal/logs"
)

type config struct {
	Address       string        `env:"ADDRESS" envDefault:"localhost:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

func main() {
	logs.Init()
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		logrus.Error(err)
	}
	var metricData = make(map[string]storage.Metric)
	if cfg.Restore {
		logrus.Info("Load data from " + cfg.StoreFile)
		metricData, _ = file.ReadMetrics(cfg.StoreFile)
	} else {
		logrus.Info("Load data is Off ")
	}
	logrus.Debug("Start server")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	var wg sync.WaitGroup
	wg.Add(1)
	go httputil.RunServer(&wg, sigChan, cfg.Address, metricData, cfg.StoreFile, cfg.StoreInterval)

	wg.Wait()

	logrus.Debug("Shutdown server")
}
