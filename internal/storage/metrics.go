package storage

import (
	"strconv"
)

func NewMetric(value string, metricType string) (Metric, error) {
	var u Metric
	u.metricType = metricType
	if metricType == MetricTypeGauge {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return Metric{}, err
		}
		u.valueGauge = v
	} else if u.metricType == MetricTypeCounter {
		newValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return Metric{}, err
		}

		u.valueCounter += newValue
	}
	return u, nil
}

func (u *Metric) UpdateMetric(value string) (Metric, error) {
	if u.metricType == MetricTypeGauge {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return Metric{}, err
		}
		u.valueGauge = v
	} else if u.metricType == MetricTypeCounter {
		newValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return Metric{}, err
		}

		u.valueCounter += newValue
	}
	return *u, nil
}

func (u *Metric) GetMetricType() string {
	return u.metricType
}
func (u *Metric) GetValue() string {
	if u.metricType == MetricTypeGauge {
		return strconv.FormatFloat(u.valueGauge, 'f', 10, 64)
	} else if u.metricType == MetricTypeCounter {
		return strconv.FormatInt(u.valueCounter, 10)
	}

	return ""
}
