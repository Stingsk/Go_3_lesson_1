package main

import (
	"context"
	"github.com/Stingsk/Go_3_lesson_1/internal/httputil"
	"github.com/Stingsk/Go_3_lesson_1/internal/logs"
	"github.com/Stingsk/Go_3_lesson_1/internal/metrics"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	logs.Init()

	logrus.Debug("Запуск сервера")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	var wg sync.WaitGroup
	var sensorData metrics.SensorData
	sensorData.Store(metrics.GetNames())
	ctx, _ := context.WithCancel(context.Background())

	wg.Add(1)
	go httputil.RunRecipient(ctx, &wg, sigChan)

	wg.Wait()

	logrus.Debug("Сервер завершил работу")
}
