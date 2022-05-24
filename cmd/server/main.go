package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Stingsk/Go_3_lesson_1/cmd/server/config"
	"github.com/Stingsk/Go_3_lesson_1/internal/httputil"
	"github.com/Stingsk/Go_3_lesson_1/internal/logs"
	"github.com/sirupsen/logrus"
)

func main() {
	logs.Init()
	serverConfig, err := config.GetServerConfig()
	if err != nil {
		logrus.Fatal("Fail to read config server: ", err)
	}
	logrus.Info("Config Server : ", serverConfig)
	level, err := logrus.ParseLevel(serverConfig.LogLevel)
	if err != nil {
		logrus.Error("Fail to read log level: ", err)
	}
	logrus.SetLevel(level)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	var wg sync.WaitGroup
	wg.Add(1)
	go httputil.RunServer(serverConfig, &wg, sigChan)

	wg.Wait()

	logrus.Debug("Shutdown server")
}
