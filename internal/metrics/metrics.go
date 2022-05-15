package metrics

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/Stingsk/Go_3_lesson_1/internal/storage"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/sirupsen/logrus"
)

type SensorData struct {
	mu   sync.RWMutex
	last []storage.Metric
}

func (d *SensorData) Store(data []storage.Metric) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.last = data
}

func (d *SensorData) Get() []storage.Metric {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.last
}

func RunGetMetrics(ctx context.Context, duration time.Duration, messages *SensorData, wg *sync.WaitGroup) {
	ticker := time.NewTicker(duration)
	count := 0
	for {
		select {
		case <-ticker.C:
			metrics := getMetrics(count)
			messages.Store(metrics)
			count++
		case <-ctx.Done():
			wg.Done()
			logrus.Error("crash shutdown")
			return
		}
	}
}

func RunGetMemoryAndCpuMetrics(ctx context.Context, duration time.Duration, messages *SensorData, wg *sync.WaitGroup) {
	ticker := time.NewTicker(duration)
	count := 0
	for {
		select {
		case <-ticker.C:
			metrics := getMemoryAndCpuMetrics(count)
			messages.Store(metrics)
			count++
		case <-ctx.Done():
			wg.Done()
			logrus.Error("crash shutdown")
			return
		}
	}
}
func getMetrics(count int) []storage.Metric {

	var metricData []storage.Metric
	var rtm runtime.MemStats
	// Read full mem stats
	runtime.ReadMemStats(&rtm)
	var ms = storage.NewMemoryStorage()

	if val, err := ms.NewMetric(strconv.Itoa(runtime.NumGoroutine()), storage.MetricTypeCounter, "NumGoroutine"); err == nil {
		metricData = append(metricData, val)
	}

	if val, err := ms.NewMetric(strconv.FormatUint(rtm.Alloc, 10), storage.MetricTypeGauge, "Alloc"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.BuckHashSys, 10), storage.MetricTypeGauge, "BuckHashSys"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.Frees, 10), storage.MetricTypeGauge, "Frees"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatFloat(rtm.GCCPUFraction, 'f', 2, 64), storage.MetricTypeGauge, "GCCPUFraction"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.GCSys, 10), storage.MetricTypeGauge, "GCSys"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.HeapAlloc, 10), storage.MetricTypeGauge, "HeapAlloc"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.HeapIdle, 10), storage.MetricTypeGauge, "HeapIdle"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.HeapInuse, 10), storage.MetricTypeGauge, "HeapInuse"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.HeapObjects, 10), storage.MetricTypeGauge, "HeapObjects"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.HeapReleased, 10), storage.MetricTypeGauge, "HeapReleased"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.HeapSys, 10), storage.MetricTypeGauge, "HeapSys"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.LastGC, 10), storage.MetricTypeGauge, "LastGC"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.Lookups, 10), storage.MetricTypeGauge, "Lookups"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.MCacheInuse, 10), storage.MetricTypeGauge, "MCacheInuse"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.MCacheSys, 10), storage.MetricTypeGauge, "MCacheSys"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.MSpanInuse, 10), storage.MetricTypeGauge, "MSpanInuse"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.MSpanSys, 10), storage.MetricTypeGauge, "MSpanSys"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.Mallocs, 10), storage.MetricTypeGauge, "Mallocs"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.NextGC, 10), storage.MetricTypeGauge, "NextGC"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(uint64(rtm.NumForcedGC), 10), storage.MetricTypeGauge, "NumForcedGC"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(uint64(rtm.NumGC), 10), storage.MetricTypeGauge, "NumGC"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.OtherSys, 10), storage.MetricTypeGauge, "OtherSys"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.PauseTotalNs, 10), storage.MetricTypeGauge, "PauseTotalNs"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.StackInuse, 10), storage.MetricTypeGauge, "StackInuse"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.StackSys, 10), storage.MetricTypeGauge, "StackSys"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.Sys, 10), storage.MetricTypeGauge, "Sys"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.Mallocs, 10), storage.MetricTypeGauge, "Mallocs"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.Frees, 10), storage.MetricTypeGauge, "Frees"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := ms.NewMetric(strconv.FormatUint(rtm.TotalAlloc, 10), storage.MetricTypeGauge, "TotalAlloc"); err == nil {
		metricData = append(metricData, val)
	}

	if val, err := ms.NewMetric(strconv.FormatUint(rand.Uint64(), 10), storage.MetricTypeGauge, "RandomValue"); err == nil {
		metricData = append(metricData, val)
	}

	if val, err := ms.NewMetric(strconv.Itoa(count), storage.MetricTypeCounter, "PollCount"); err == nil {
		metricData = append(metricData, val)
	}

	return metricData
}

func getMemoryAndCpuMetrics(count int) []storage.Metric {

	var ms = storage.NewMemoryStorage()
	var metricData []storage.Metric

	memory, err := mem.VirtualMemory()
	if err == nil {
		if val, err := ms.NewMetric(strconv.FormatUint(memory.Total, 10), storage.MetricTypeGauge, "TotalMemory"); err == nil {
			metricData = append(metricData, val)
		}
		if val, err := ms.NewMetric(strconv.FormatUint(memory.Free, 10), storage.MetricTypeGauge, "FreeMemory"); err == nil {
			metricData = append(metricData, val)
		}
	}

	processor, err := cpu.Percent(1*time.Second, true)
	if err == nil {
		for index, p := range processor {
			if val, err := ms.NewMetric(strconv.FormatFloat(p, 'f', 2, 64), storage.MetricTypeGauge, fmt.Sprintf("CPUutilization%d", index+1)); err == nil {
				metricData = append(metricData, val)
			}
		}
	}

	return metricData
}
