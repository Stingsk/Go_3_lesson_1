package main

import (
	"context"
	"github.com/Stingsk/Go_3_lesson_1/internal/httputil"
	"github.com/Stingsk/Go_3_lesson_1/internal/metrics"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
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
	go metrics.RunGetMetrics(ctx, 2, &sensorData, &wg, sigChan)

	wg.Add(2)
	go httputil.RunSender(ctx, 10, &sensorData, &wg, sigChan)

	wg.Wait()
}
