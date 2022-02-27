package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"

	"github.com/Stingsk/Go_3_lesson_1/internal/httputil"
	"github.com/Stingsk/Go_3_lesson_1/internal/logs"
)

type config struct {
	Address string `env:"ADDRESS" envDefault:"http://localhost:8080"`
}

func main() {
	logs.Init()
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
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
	go httputil.RunServer(&wg, sigChan, cfg.Address)

	wg.Wait()

	logrus.Debug("Shutdown server")
}
