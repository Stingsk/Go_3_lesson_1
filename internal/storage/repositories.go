package storage

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

// MetricType is enum-like type.
// We are using struct instead of string, to ensure about immutability.
type MetricName struct {
	s string
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

func (u MetricName) IsZero() bool {
	return u == MetricName{}
}

func (u MetricName) String() string {
	return u.s
}

type Metric struct {
	metricType MetricType
	metricName MetricName
	counter    int
	value      string
}

type Repository interface {
	NewMetricTypeString(metricType string) (MetricType, error)
	NewMetricNameString(metricName string) (MetricName, error)
	NewMetricString(metricName string, metricType string, value string) error
	UpdateMetric(value string) Metric
	GetMetricName() MetricName
}
