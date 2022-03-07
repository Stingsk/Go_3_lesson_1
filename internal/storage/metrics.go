package storage

import (
	"strconv"
	"strings"
)

func NewMetricResourceFromParams(value string, metricType string, name string) (MetricResource, error) {
	var u MetricResource
	u.Mutex.TryLock()
	defer u.Mutex.Unlock()

	u.Metric = &Metric{
		ID:    strings.ToLower(metricType),
		MType: strings.ToLower(name),
		Delta: nil,
		Value: nil,
	}
	if strings.ToLower(metricType) == MetricTypeGauge {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return MetricResource{}, err
		}
		u.Metric.Value = &v
	} else if strings.ToLower(u.Metric.MType) == MetricTypeCounter {
		newValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return MetricResource{}, err
		}

		delta := int64(0)
		if u.Metric.Delta != nil {
			delta = *u.Metric.Delta
		}
		u.Metric.Delta = sumInt(delta, newValue)
	}
	updated := true
	u.Updated = &updated

	return u, nil
}

func (u *MetricResource) UpdateMetricResource(value string) error {
	u.Mutex.TryLock()
	defer u.Mutex.Unlock()

	if strings.ToLower(u.Metric.MType) == MetricTypeGauge {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		u.Metric.Value = &v
	} else if strings.ToLower(u.Metric.MType) == MetricTypeCounter {
		newValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}

		delta := int64(0)
		if u.Metric.Delta != nil {
			delta = *u.Metric.Delta
		}
		u.Metric.Delta = sumInt(delta, newValue)
	}
	updated := true
	u.Updated = &updated

	return nil
}

func (u *MetricResource) Update(newMetric Metric) {
	u.Mutex.TryLock()
	if strings.ToLower(newMetric.MType) == MetricTypeGauge {
	} else if strings.ToLower(newMetric.MType) == MetricTypeCounter {
		newMetric.Delta = sumInt(*u.Metric.Delta, *newMetric.Delta)
	}

	u.Metric = &newMetric
	updated := true
	u.Updated = &updated
	u.Mutex.Unlock()
}

func (u *MetricResource) GetMetricType() string {
	return u.Metric.MType
}
func (u *MetricResource) GetValue() string {
	u.Mutex.TryLock()
	defer u.Mutex.Unlock()
	if strings.ToLower(u.Metric.MType) == MetricTypeGauge {
		if u.Metric.Value == nil {
			return ""
		}
		return strconv.FormatFloat(*u.Metric.Value, 'f', 3, 64)
	} else if strings.ToLower(u.Metric.MType) == MetricTypeCounter {
		if u.Metric.Delta == nil {
			return ""
		}
		return strconv.FormatInt(*u.Metric.Delta, 10)
	}

	return ""
}

func sumInt(first int64, second int64) *int64 {
	helper := first + second
	return &helper
}
func sumFloat(first float64, second float64) *float64 {
	helper := first + second
	return &helper
}

func NewMetricResource(metric Metric) MetricResource {
	falseValue := true
	return MetricResource{
		Metric:  &metric,
		Updated: &falseValue,
	}
}
