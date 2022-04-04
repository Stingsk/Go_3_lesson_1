package storage

import "sync"

const (
	MetricTypeGauge   string = "gauge"
	MetricTypeCounter string = "counter"
)

type Repository interface {
	NewMetric(metricName string, metricType string, value string) error
	UpdateMetric(value string) (Metric, error)
}

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type MetricResource struct {
	Metric  *Metric
	Updated *bool
	Mutex   *sync.Mutex
}

type MetricResourceMap struct {
	Metric *map[string]MetricResource
}
