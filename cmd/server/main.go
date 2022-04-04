package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Stingsk/Go_3_lesson_1/cmd/server/config"
	"github.com/Stingsk/Go_3_lesson_1/internal/file"
	"github.com/Stingsk/Go_3_lesson_1/internal/httputil"
	"github.com/Stingsk/Go_3_lesson_1/internal/logs"
	"github.com/Stingsk/Go_3_lesson_1/internal/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	logs.Init()
	var metricData storage.MetricResourceMap
	cfg := config.GetConfig()
	if cfg.Restore {
		logrus.Info("Load data from " + cfg.StoreFile)
		metricRead, _ := file.ReadMetrics(cfg.StoreFile)
		metricData.Metric = &metricRead
	} else {
		metricResource := make(map[string]storage.MetricResource)
		metricData.Metric = &metricResource
		logrus.Info("Load data is Off ")
	}
	err := os.Remove(cfg.StoreFile)
	if err != nil {
		logrus.Error(err)
	}

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
