package main

import (
	"context"
	"github.com/Stingsk/Go_3_lesson_1/internal/config"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Stingsk/Go_3_lesson_1/internal/httputil"
	"github.com/Stingsk/Go_3_lesson_1/internal/logs"
	"github.com/Stingsk/Go_3_lesson_1/internal/metrics"
	"github.com/sirupsen/logrus"
)

func main() {
	logs.Init()
	cfg := config.GetAgentConfig()
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
