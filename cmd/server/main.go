package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Stingsk/Go_3_lesson_1/internal/config"
	"github.com/Stingsk/Go_3_lesson_1/internal/httputil"
	"github.com/Stingsk/Go_3_lesson_1/internal/logs"
	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
)

func main() {
	logs.Init()
	if err := config.GetServerConfig(); err != nil {
		logrus.Info("Failed to parse command line arguments", err)
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	var wg sync.WaitGroup
	wg.Add(1)
	var serverConfig = &httputil.Config{
		WaitGroup:          &wg,
		SigChan:            sigChan,
		Address:            config.Address,
		Restore:            config.Restore,
		StoreFile:          config.StoreFile,
		StoreInterval:      config.StoreInterval,
		SignKey:            config.SignKey,
		DataBaseConnection: config.DataBaseConnection, // "postgresql://localhost:5432/metrics",
		LogLevel:           config.LogLevel,
	}

	level, err := logrus.ParseLevel(serverConfig.LogLevel)
	if err != nil {
		logrus.Error(err)
	}
	logrus.SetLevel(level)
	if err := env.Parse(&serverConfig); err != nil {
		logrus.Info("Failed to parse environment variables", err)
	}
	logrus.Info("Config Server : ", serverConfig)
	go httputil.RunServer(*serverConfig)

	wg.Wait()

	logrus.Debug("Shutdown server")
}
