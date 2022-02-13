package storage

import (
	"errors"
	"fmt"
)

func NewMetricTypeString(metricType string) (MetricType, error) {
	switch metricType {
	case "gauge":
		return Gauge, nil
	case "counter":
		return Counter, nil
	}

	return MetricType{}, fmt.Errorf("invalid '%s' MetricType", metricType)
}

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
	case "numgoroutine":
		return NumGoroutine, nil
	case "pullcounter":
		return PullCounter, nil
	case "testgauge":
		return TestGauge, nil
	case "testcounter":
		return TestCounter, nil
	case "testsetget134":
		return TestSetGet134, nil
	case "testsetget33":
		return TestSetGet33, nil
	}

	return MetricName{}, fmt.Errorf("invalid '%s' MetricName", metricName)
}

func (u *Metric) NewMetricString(metricName string, metricType string, value string) error {
	mt, err := NewMetricTypeString(metricType)
	if err != nil {
		return errors.New("error Match MetricType")
	}
	mn, err := NewMetricNameString(metricName)
	if err != nil {
		return errors.New("error Match MetricName")
	}

	u.metricType = mt
	u.metricName = mn
	u.counter = 0
	u.value = value

	return nil
}

func (u *Metric) UpdateMetric(value string) Metric {
	u.value = value
	u.counter++
	return *u
}
func (u *Metric) GetMetricName() MetricName {
	return u.metricName
}

func (u *Metric) GetMetricType() MetricType {
	return u.metricType
}
func (u *Metric) GetValue() string {
	return u.value
}
