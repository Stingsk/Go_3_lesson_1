package storage

const (
	MetricTypeGauge   string = "gauge"
	MetricTypeCounter string = "counter"
)

type Metric struct {
	metricType   string
	valueGauge   float64
	valueCounter int64
}

type Repository interface {
	NewMetricString(metricName string, metricType string, value string) error
	UpdateMetric(value string, metricType string) Metric
}
