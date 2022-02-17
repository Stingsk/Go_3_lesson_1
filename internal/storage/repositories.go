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
	NewMetric(metricName string, metricType string, value string) error
	UpdateMetric(value string) (Metric, error)
}
