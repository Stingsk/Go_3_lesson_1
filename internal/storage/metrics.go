package storage

import (
	"errors"
	"fmt"
)

// MetricType is enum-like type.
// We are using struct instead of string, to ensure about immutability.
type MetricType struct {
	s string
}

func (u MetricType) IsZero() bool {
	return u == MetricType{}
}

func (u MetricType) String() string {
	return u.s
}

var (
	Gauge   = MetricType{"gauge"}
	Counter = MetricType{"counter"}
)

func NewMetricTypeString(metricType string) (MetricType, error) {
	switch metricType {
	case "gauge":
		return Gauge, nil
	case "counter":
		return Counter, nil
	}

	return MetricType{}, errors.New(fmt.Sprintf("invalid '%s' MetricType", metricType))
}

// MetricType is enum-like type.
// We are using struct instead of string, to ensure about immutability.
type MetricName struct {
	s string
}

func (u MetricName) IsZero() bool {
	return u == MetricName{}
}

func (u MetricName) String() string {
	return u.s
}

var (
	Alloc         = MetricName{"alloc"}
	BuckHashSys   = MetricName{"buckhashsys"}
	Frees         = MetricName{"frees"}
	GCSys         = MetricName{"gcsys"}
	HeapAlloc     = MetricName{"heapalloc"}
	HeapIdle      = MetricName{"heapidle"}
	HeapInuse     = MetricName{"heapinuse"}
	HeapObjects   = MetricName{"heapobjects"}
	HeapReleased  = MetricName{"heapreleased"}
	HeapSys       = MetricName{"heapsys"}
	LastGC        = MetricName{"lastgc"}
	Lookups       = MetricName{"lookups"}
	MCacheInuse   = MetricName{"mcacheinuse"}
	MCacheSys     = MetricName{"mcachesys"}
	MSpanInuse    = MetricName{"mspaninuse"}
	MSpanSys      = MetricName{"mspansys"}
	Mallocs       = MetricName{"mallocs"}
	NextGC        = MetricName{"nextgc"}
	OtherSys      = MetricName{"othersys"}
	PauseTotalNs  = MetricName{"pausetotalns"}
	StackInuse    = MetricName{"stackinuse"}
	StackSys      = MetricName{"stacksys"}
	RandomValue   = MetricName{"randomvalue"}
	Sys           = MetricName{"sys"}
	NumGC         = MetricName{"numgc"}
	NumForcedGC   = MetricName{"numforcedgc"}
	GCCPUFraction = MetricName{"gccpufraction"}
	NumGoroutine  = MetricName{"numgoroutine"}
	PollCount     = MetricName{"pollcount"}
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
		return MCacheInuse, nil
	case "mspansys":
		return MSpanInuse, nil
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
		return NextGC, nil
	case "numforcedgc":
		return NumForcedGC, nil
	case "gccpufraction":
		return GCCPUFraction, nil
	case "numgoroutine":
		return NumGoroutine, nil
	case "pollcount":
		return PollCount, nil
	}

	return MetricName{}, errors.New(fmt.Sprintf("invalid '%s' MetricName", metricName))
}

type Metric struct {
	metricType MetricType
	metricName MetricName
	counter    int
	value      string
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
