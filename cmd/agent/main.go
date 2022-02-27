package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"

	"github.com/Stingsk/Go_3_lesson_1/internal/httputil"
	"github.com/Stingsk/Go_3_lesson_1/internal/logs"
	"github.com/Stingsk/Go_3_lesson_1/internal/metrics"
)

type config struct {
	Address        string `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval int    `env:"REPORT_INTERVAL" envDefault:"10"`
	PollInterval   int    `env:"POLL_INTERVAL" envDefault:"2"`
}

func main() {
	logs.Init()

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		logrus.Error("%+v\n", err)
	}
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
