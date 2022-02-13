package storage

import (
	"fmt"
	"strconv"
)

func NewMetricNameString(metricName string) (MetricName, error) {
	switch metricName {
	case "alloc":
		return Alloc, nil
	case "buckhashsys":
		return BuckHashSys, nil
	case "frees":
		return Frees, nil
	case "gcsys":
		return GCSys, nil
	case "heapalloc":
		return HeapAlloc, nil
	case "heapidle":
		return HeapIdle, nil
	case "heapinuse":
		return HeapInuse, nil
	case "heapobjects":
		return HeapObjects, nil
	case "heapreleased":
		return HeapReleased, nil
	case "heapsys":
		return HeapSys, nil
	case "lastgc":
		return LastGC, nil
	case "lookups":
		return Lookups, nil
	case "mcacheinuse":
		return MCacheInuse, nil
	case "mcachesys":
		return MCacheSys, nil
	case "mspaninuse":
		return MSpanInuse, nil
	case "mspansys":
		return MSpanSys, nil
	case "mallocs":
		return Mallocs, nil
	case "nextgc":
		return NextGC, nil
	case "othersys":
		return OtherSys, nil
	case "pausetotalns":
		return PauseTotalNs, nil
	case "stackinuse":
		return StackInuse, nil
	case "stacksys":
		return StackSys, nil
	case "randomvalue":
		return RandomValue, nil
	case "sys":
		return Sys, nil
	case "numgc":
		return NumGC, nil
	case "numforcedgc":
		return NumForcedGC, nil
	case "gccpufraction":
		return GCCPUFraction, nil
	}

	return MetricName{}, fmt.Errorf("invalid '%s' MetricName", metricName)
}

func (u *Metric) NewMetricString(metricName string, metricType string, value string) {
	u.metricType = metricType
	u.metricName = metricName
	u.value = value
}
func (u *MetricName) NewMetricNameString(metricName string) {
	u.s = metricName
}

func (u *Metric) UpdateMetric(value string, metricType string) Metric {
	if metricType == "gauge" {
		u.value = value
	} else if metricType == "counter" {
		oldValue, err := strconv.ParseInt(u.value, 10, 64)
		if err != nil {
			oldValue = 0
		}
		newValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			newValue = 0
		}
		u.value = strconv.FormatInt(oldValue+newValue, 10)
	}
	return *u
}
func (u *Metric) GetMetricName() string {
	return u.metricName
}

func (u *Metric) GetMetricType() string {
	return u.metricType
}
func (u *Metric) GetValue() string {
	return u.value
}
