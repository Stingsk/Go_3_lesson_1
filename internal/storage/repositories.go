package storage

import "sync"

const (
	MetricTypeGauge   string = "gauge"
	MetricTypeCounter string = "counter"
)

type Repository interface {
	NewMetric(value string, metricType string, name string) (Metric, error)
	UpdateMetric(metricResourceMap *MetricResourceMap, metric Metric, singKey string) (Metric, error)
	UpdateMetricByParameters(metricResourceMap *MetricResourceMap, metricName string, metricType string, value string) (Metric, error)
}

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // Значение хеш-функции
}

type MetricResourceMap struct {
	Metric     map[string]Metric
	Mutex      sync.Mutex
	Repository Repository
}
