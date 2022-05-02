package storage

import (
	"context"
)

const (
	MetricTypeGauge   string = "gauge"
	MetricTypeCounter string = "counter"
)

type Gauge float64
type Counter int64

type Repository interface {
	GetMetric(ctx context.Context, name string, metricType string) (*Metric, error)
	GetMetrics(ctx context.Context) (map[string]*Metric, error)

	UpdateMetrics(ctx context.Context, metricsBatch []*Metric) error
	UpdateMetricByParameters(ctx context.Context, metricName string, metricType string, value string) error
	UpdateMetric(ctx context.Context, metric Metric) error

	Ping(ctx context.Context) error
}

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *Counter `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *Gauge   `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // Значение хеш-функции
}
