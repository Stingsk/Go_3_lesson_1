package metrics

import (
	"context"
	"errors"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/Stingsk/Go_3_lesson_1/internal/storage"
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

func RunGetMetrics(ctx context.Context, duration int, messages *SensorData, wg *sync.WaitGroup) error {
	ticker := time.NewTicker(time.Duration(duration) * time.Second)
	count := 0
	for {
		select {
		case <-ticker.C:
			metrics := getMetrics(count)
			messages.Store(metrics)
			count++
		case <-ctx.Done():
			wg.Done()
			return errors.New("crash shutdown")
		}
	}
}

func getMetrics(count int) []storage.Metric {

	var metricData []storage.Metric
	var rtm runtime.MemStats
	// Read full mem stats
	runtime.ReadMemStats(&rtm)

	if val, err := storage.NewMetric(strconv.Itoa(runtime.NumGoroutine()), storage.MetricTypeCounter, "numgoroutine"); err == nil {
		metricData = append(metricData, val)
	}

	if val, err := storage.NewMetric(strconv.FormatUint(rtm.Alloc, 10), storage.MetricTypeGauge, "alloc"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.BuckHashSys, 10), storage.MetricTypeGauge, "buckhashsys"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.Frees, 10), storage.MetricTypeGauge, "frees"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatFloat(rtm.GCCPUFraction, 'f', 2, 64), storage.MetricTypeGauge, "gccpufraction"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.GCSys, 10), storage.MetricTypeGauge, "gcsys"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.HeapAlloc, 10), storage.MetricTypeGauge, "heapalloc"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.HeapIdle, 10), storage.MetricTypeGauge, "heapidle"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.HeapInuse, 10), storage.MetricTypeGauge, "heapinuse"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.HeapObjects, 10), storage.MetricTypeGauge, "heapobjects"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.HeapReleased, 10), storage.MetricTypeGauge, "heapreleased"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.HeapSys, 10), storage.MetricTypeGauge, "heapsys"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.LastGC, 10), storage.MetricTypeGauge, "lastgc"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.Lookups, 10), storage.MetricTypeGauge, "lookups"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.MCacheInuse, 10), storage.MetricTypeGauge, "mcacheinuse"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.MCacheSys, 10), storage.MetricTypeGauge, "mcachesys"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.MSpanInuse, 10), storage.MetricTypeGauge, "mspaninuse"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.MSpanSys, 10), storage.MetricTypeGauge, "mspansys"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.Mallocs, 10), storage.MetricTypeGauge, "mallocs"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.NextGC, 10), storage.MetricTypeGauge, "nextgc"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(uint64(rtm.NumForcedGC), 10), storage.MetricTypeGauge, "numforcedgc"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(uint64(rtm.NumGC), 10), storage.MetricTypeGauge, "numgc"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.OtherSys, 10), storage.MetricTypeGauge, "othersys"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.PauseTotalNs, 10), storage.MetricTypeGauge, "pausetotalns"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.StackInuse, 10), storage.MetricTypeGauge, "stackinuse"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.StackSys, 10), storage.MetricTypeGauge, "stacksys"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.Sys, 10), storage.MetricTypeGauge, "sys"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.Mallocs, 10), storage.MetricTypeGauge, "mallocs"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.Frees, 10), storage.MetricTypeGauge, "frees"); err == nil {
		metricData = append(metricData, val)
	}

	if val, err := storage.NewMetric(strconv.FormatUint(rand.Uint64(), 10), storage.MetricTypeGauge, "randomvalue"); err == nil {
		metricData = append(metricData, val)
	}

	if val, err := storage.NewMetric(strconv.Itoa(count), storage.MetricTypeCounter, "pollcount"); err == nil {
		metricData = append(metricData, val)
	}

	return metricData
}
