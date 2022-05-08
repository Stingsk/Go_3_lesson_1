package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Stingsk/Go_3_lesson_1/internal/config"

	"github.com/Stingsk/Go_3_lesson_1/internal/httputil"
	"github.com/Stingsk/Go_3_lesson_1/internal/logs"
	"github.com/sirupsen/logrus"
)

func main() {
	logs.Init()
	cfg := config.GetServerConfig()
	logrus.Debug("Start server")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	var wg sync.WaitGroup
	wg.Add(1)
	/*var serverConfig = &httputil.ServerConfig{
		WaitGroup:          &wg,
		SigChan:            sigChan,
		Host:               cfg.Address,
		Restore:            cfg.Restore,
		StoreFile:          cfg.StoreFile,
		StoreInterval:      cfg.StoreInterval,
		SignKey:            cfg.SignKey,
		DataBaseConnection: cfg.DataBaseConnection,
	}*/

	var serverConfig = &httputil.ServerConfig{
		WaitGroup:          &wg,
		SigChan:            sigChan,
		Host:               cfg.Address,
		Restore:            cfg.Restore,
		StoreFile:          "",
		StoreInterval:      cfg.StoreInterval,
		SignKey:            cfg.SignKey,
		DataBaseConnection: "",
	}
	go httputil.RunServer(*serverConfig)

	wg.Wait()

	logrus.Debug("Shutdown server")
}
