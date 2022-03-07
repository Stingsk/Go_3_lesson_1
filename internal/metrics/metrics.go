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
	last []storage.MetricResource
}

func (d *SensorData) Store(data []storage.MetricResource) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.last = data
}

func (d *SensorData) Get() []storage.MetricResource {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.last
}

func RunGetMetrics(ctx context.Context, duration time.Duration, messages *SensorData, wg *sync.WaitGroup) error {
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
			return errors.New("crash shutdown")
		}
	}
}

func getMetrics(count int) []storage.MetricResource {

	var metricData []storage.MetricResource
	var rtm runtime.MemStats
	// Read full mem stats
	runtime.ReadMemStats(&rtm)

	if val, err := storage.NewMetricResourceFromParams(strconv.Itoa(runtime.NumGoroutine()), storage.MetricTypeCounter, "numgoroutine"); err == nil {
		metricData = append(metricData, val)
	}

	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.Alloc, 10), storage.MetricTypeGauge, "alloc"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.BuckHashSys, 10), storage.MetricTypeGauge, "buckhashsys"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.Frees, 10), storage.MetricTypeGauge, "frees"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatFloat(rtm.GCCPUFraction, 'f', 2, 64), storage.MetricTypeGauge, "gccpufraction"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.GCSys, 10), storage.MetricTypeGauge, "gcsys"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.HeapAlloc, 10), storage.MetricTypeGauge, "heapalloc"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.HeapIdle, 10), storage.MetricTypeGauge, "heapidle"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.HeapInuse, 10), storage.MetricTypeGauge, "heapinuse"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.HeapObjects, 10), storage.MetricTypeGauge, "heapobjects"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.HeapReleased, 10), storage.MetricTypeGauge, "heapreleased"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.HeapSys, 10), storage.MetricTypeGauge, "heapsys"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.LastGC, 10), storage.MetricTypeGauge, "lastgc"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.Lookups, 10), storage.MetricTypeGauge, "lookups"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.MCacheInuse, 10), storage.MetricTypeGauge, "mcacheinuse"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.MCacheSys, 10), storage.MetricTypeGauge, "mcachesys"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.MSpanInuse, 10), storage.MetricTypeGauge, "mspaninuse"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.MSpanSys, 10), storage.MetricTypeGauge, "mspansys"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.Mallocs, 10), storage.MetricTypeGauge, "mallocs"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.NextGC, 10), storage.MetricTypeGauge, "nextgc"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(uint64(rtm.NumForcedGC), 10), storage.MetricTypeGauge, "numforcedgc"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(uint64(rtm.NumGC), 10), storage.MetricTypeGauge, "numgc"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.OtherSys, 10), storage.MetricTypeGauge, "othersys"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.PauseTotalNs, 10), storage.MetricTypeGauge, "pausetotalns"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.StackInuse, 10), storage.MetricTypeGauge, "stackinuse"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.StackSys, 10), storage.MetricTypeGauge, "stacksys"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.Sys, 10), storage.MetricTypeGauge, "sys"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.Mallocs, 10), storage.MetricTypeGauge, "mallocs"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.Frees, 10), storage.MetricTypeGauge, "frees"); err == nil {
		metricData = append(metricData, val)
	}
	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rtm.TotalAlloc, 10), storage.MetricTypeGauge, "totalalloc"); err == nil {
		metricData = append(metricData, val)
	}

	if val, err := storage.NewMetricResourceFromParams(strconv.FormatUint(rand.Uint64(), 10), storage.MetricTypeGauge, "randomvalue"); err == nil {
		metricData = append(metricData, val)
	}

	if val, err := storage.NewMetricResourceFromParams(strconv.Itoa(count), storage.MetricTypeCounter, "pollcount"); err == nil {
		metricData = append(metricData, val)
	}

	return metricData
}
