package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Stingsk/Go_3_lesson_1/cmd/agent/config"
	"github.com/Stingsk/Go_3_lesson_1/internal/httputil"
	"github.com/Stingsk/Go_3_lesson_1/internal/logs"
	"github.com/Stingsk/Go_3_lesson_1/internal/metrics"
	"github.com/sirupsen/logrus"
)

func main() {
	logs.Init()
	agentConfig, err := config.GetAgentConfig()
	if err != nil {
		logrus.Fatal("Failed to parse command line arguments:", err)
	}

	level, err := logrus.ParseLevel(agentConfig.LogLevel)
	if err != nil {
		logrus.Error(err)
	}
	logrus.SetLevel(level)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	wg := &sync.WaitGroup{}
	var sensorData metrics.SensorData

	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go metrics.RunGetMetrics(ctx, agentConfig.PollInterval, &sensorData, wg)
	wg.Add(1)
	go metrics.RunGetMemoryAndCPUMetrics(ctx, agentConfig.PollInterval, &sensorData, wg)

	wg.Add(1)

	go httputil.RunSender(agentConfig, &sensorData, wg, ctx)

	go func() {
		<-sigChan
		cancel()
	}()

	wg.Wait()
	logrus.Debug("Shutdown agent")
}
