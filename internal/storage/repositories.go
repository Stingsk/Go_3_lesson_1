package storage

import (
	"context"
)

const (
	MetricTypeGauge   string = "gauge"
	MetricTypeCounter string = "counter"
)

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
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // Значение хеш-функции
}
