package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Stingsk/Go_3_lesson_1/internal/config"
	"github.com/caarlos0/env/v6"

	"github.com/Stingsk/Go_3_lesson_1/internal/httputil"
	"github.com/Stingsk/Go_3_lesson_1/internal/logs"
	"github.com/Stingsk/Go_3_lesson_1/internal/metrics"
	"github.com/sirupsen/logrus"
)

func main() {
	logs.Init()
	if err := config.GetAgentConfig(); err != nil {
		logrus.Error("Failed to parse command line arguments:", err)
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

	var agentConfig = &httputil.AgentConfig{
		Context:        ctx,
		ReportInterval: config.ReportInterval,
		PollInterval:   config.PollInterval,
		Messages:       &sensorData,
		WaitGroup:      wg,
		Address:        config.Address,
		SignKey:        config.SignKey,
		LogLevel:       config.LogLevel,
	}

	if err := env.Parse(&agentConfig); err != nil {
		logrus.Error("Failed to parse environment variables", err)
	}
	level, err := logrus.ParseLevel(agentConfig.LogLevel)
	if err != nil {
		logrus.Error(err)
	}
	logrus.SetLevel(level)

	wg.Add(1)
	go metrics.RunGetMetrics(ctx, agentConfig.PollInterval, &sensorData, wg)

	wg.Add(1)

	go httputil.RunSender(*agentConfig)

	go func() {
		<-sigChan
		cancel()
	}()

	wg.Wait()
	logrus.Debug("Shutdown agent")
}
