package storage

import (
	"strconv"
	"strings"
)

func NewMetric(value string, metricType string, name string) (Metric, error) {
	var u Metric
	u.MType = metricType
	u.ID = name
	if strings.ToLower(metricType) == MetricTypeGauge {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return Metric{}, err
		}
		u.Value = &v
	} else if strings.ToLower(u.MType) == MetricTypeCounter {
		newValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return Metric{}, err
		}

		delta := int64(0)
		if u.Delta != nil {
			delta = *u.Delta
		}
		u.Delta = getAdress(delta + newValue)
	}
	return u, nil
}

func UpdateMetric(value string, u Metric) (Metric, error) {
	if strings.ToLower(u.MType) == MetricTypeGauge {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return Metric{}, err
		}
		u.Value = &v
	} else if strings.ToLower(u.MType) == MetricTypeCounter {
		newValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return Metric{}, err
		}

		delta := int64(0)
		if u.Delta != nil {
			delta = *u.Delta
		}
		u.Delta = getAdress(delta + newValue)
	}
	return u, nil
}

func Update(newMetric Metric, u Metric) Metric {
	if strings.ToLower(newMetric.MType) == MetricTypeGauge {
		return newMetric
	} else if strings.ToLower(newMetric.MType) == MetricTypeCounter {
		newMetric.Delta = getAdress(*u.Delta + *newMetric.Delta)
	}
	return newMetric
}

func (u *Metric) GetMetricType() string {
	return u.MType
}
func (u *Metric) GetValue() string {
	if strings.ToLower(u.MType) == MetricTypeGauge {
		if u.Value == nil {
			return ""
		}
		return strconv.FormatFloat(*u.Value, 'f', 3, 64)
	} else if strings.ToLower(u.MType) == MetricTypeCounter {
		if u.Delta == nil {
			return ""
		}
		return strconv.FormatInt(*u.Delta, 10)
	}

	return ""
}

func getAdress[T any](t T) *T {
	return &t
}
