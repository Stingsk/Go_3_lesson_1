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
	last []string
}

func (d *SensorData) Store(data []string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.last = data
}

func (d *SensorData) Get() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.last
}

func RunGetMetrics(ctx context.Context, duration int, messages *SensorData, wg *sync.WaitGroup) error {
	ticker := time.NewTicker(time.Duration(duration) * time.Second)
	for {
		select {
		case <-ticker.C:
			metrics := getMetrics()
			messages.Store(metrics)
		case <-ctx.Done():
			wg.Done()
			return errors.New("crash shutdown")
		}
	}
}

func getMetrics() []string {
	monitor := newMonitor()

	result := make([]string, len(monitor))
	i := 0
	for name, metric := range monitor {
		result[i] = metric.GetMetricType() + "/" + name + "/" + metric.GetValue()
		i++
	}

	return result
}

func newMonitor() map[string]storage.Metric {

	var metricData = make(map[string]storage.Metric)
	var rtm runtime.MemStats
	// Read full mem stats
	runtime.ReadMemStats(&rtm)

	if val, err := storage.NewMetric(strconv.Itoa(runtime.NumGoroutine()), storage.MetricTypeCounter); err == nil {
		metricData["numgoroutine"] = val
	}

	if val, err := storage.NewMetric(strconv.FormatUint(rtm.Alloc, 10), storage.MetricTypeGauge); err == nil {
		metricData["alloc"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.BuckHashSys, 10), storage.MetricTypeGauge); err == nil {
		metricData["buckhashsys"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.Frees, 10), storage.MetricTypeGauge); err == nil {
		metricData["frees"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatFloat(rtm.GCCPUFraction, 'f', 2, 64), storage.MetricTypeGauge); err == nil {
		metricData["gccpufraction"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.GCSys, 10), storage.MetricTypeGauge); err == nil {
		metricData["gcsys"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.HeapAlloc, 10), storage.MetricTypeGauge); err == nil {
		metricData["heapalloc"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.HeapIdle, 10), storage.MetricTypeGauge); err == nil {
		metricData["heapidle"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.HeapInuse, 10), storage.MetricTypeGauge); err == nil {
		metricData["heapinuse"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.HeapObjects, 10), storage.MetricTypeGauge); err == nil {
		metricData["heapobjects"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.HeapReleased, 10), storage.MetricTypeGauge); err == nil {
		metricData["heapreleased"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.HeapSys, 10), storage.MetricTypeGauge); err == nil {
		metricData["heapsys"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.LastGC, 10), storage.MetricTypeGauge); err == nil {
		metricData["lastgc"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.Lookups, 10), storage.MetricTypeGauge); err == nil {
		metricData["lookups"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.MCacheInuse, 10), storage.MetricTypeGauge); err == nil {
		metricData["mcacheinuse"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.MCacheSys, 10), storage.MetricTypeGauge); err == nil {
		metricData["mcachesys"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.MSpanInuse, 10), storage.MetricTypeGauge); err == nil {
		metricData["mspaninuse"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.MSpanSys, 10), storage.MetricTypeGauge); err == nil {
		metricData["mspansys"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.Mallocs, 10), storage.MetricTypeGauge); err == nil {
		metricData["mallocs"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.NextGC, 10), storage.MetricTypeGauge); err == nil {
		metricData["nextgc"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(uint64(rtm.NumForcedGC), 10), storage.MetricTypeGauge); err == nil {
		metricData["numforcedgc"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(uint64(rtm.NumGC), 10), storage.MetricTypeGauge); err == nil {
		metricData["numgc"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.OtherSys, 10), storage.MetricTypeGauge); err == nil {
		metricData["othersys"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.PauseTotalNs, 10), storage.MetricTypeGauge); err == nil {
		metricData["pausetotalns"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.StackInuse, 10), storage.MetricTypeGauge); err == nil {
		metricData["stackinuse"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.StackSys, 10), storage.MetricTypeGauge); err == nil {
		metricData["stacksys"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.Sys, 10), storage.MetricTypeGauge); err == nil {
		metricData["sys"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.Mallocs, 10), storage.MetricTypeGauge); err == nil {
		metricData["mallocs"] = val
	}
	if val, err := storage.NewMetric(strconv.FormatUint(rtm.Frees, 10), storage.MetricTypeGauge); err == nil {
		metricData["frees"] = val
	}

	if val, err := storage.NewMetric(strconv.FormatUint(rand.Uint64(), 10), storage.MetricTypeGauge); err == nil {
		metricData["randomvalue"] = val
	}

	return metricData
}
