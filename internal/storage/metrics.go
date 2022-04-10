package storage

import (
	"errors"
	"strconv"
	"strings"
)

func (m *MemoryStorage) NewMetric(value string, metricType string, name string) (Metric, error) {
	metric := Metric{
		ID:    strings.ToLower(name),
		MType: strings.ToLower(metricType),
		Delta: nil,
		Value: nil,
	}
	if strings.ToLower(metricType) == MetricTypeGauge {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return Metric{}, err
		}
		metric.Value = &v
	} else if strings.ToLower(metric.MType) == MetricTypeCounter {
		newValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return Metric{}, err
		}

		delta := int64(0)
		if metric.Delta != nil {
			delta = *metric.Delta
		}
		metric.Delta = sumInt(delta, newValue)
	}

	return metric, nil
}

func (m *Metric) UpdateMetricResource(value string) error {
	if strings.ToLower(m.MType) == MetricTypeGauge {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		m.Value = &v
	} else if strings.ToLower(m.MType) == MetricTypeCounter {
		newValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}

		delta := int64(0)
		if m.Delta != nil {
			delta = *m.Delta
		}
		m.Delta = sumInt(delta, newValue)
	}

	return nil
}

func (m *Metric) Update(newMetric Metric) {
	if strings.ToLower(newMetric.MType) == MetricTypeGauge {
		m.Value = newMetric.Value
	} else if strings.ToLower(newMetric.MType) == MetricTypeCounter {
		m.Delta = sumInt(*m.Delta, *newMetric.Delta)
	}
}

func (m *Metric) GetMetricType() string {
	return m.MType
}

func (m *Metric) GetValue() string {
	if strings.ToLower(m.MType) == MetricTypeGauge {
		if m.Value == nil {
			return ""
		}
		return strconv.FormatFloat(*m.Value, 'f', 3, 64)
	} else if strings.ToLower(m.MType) == MetricTypeCounter {
		if m.Delta == nil {
			return ""
		}
		return strconv.FormatInt(*m.Delta, 10)
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

func (m *MemoryStorage) UpdateMetric(metricResourceMap *MetricResourceMap, metric Metric) (Metric, error) {
	metricResourceMap.Mutex.Lock()
	defer metricResourceMap.Mutex.Unlock()

	var valueMetric = metricResourceMap.Metric[strings.ToLower(metric.ID)]
	if valueMetric.GetValue() != "" {
		if metric.Delta != nil || metric.Value != nil {
			valueMetric.Update(metric)
			metricResourceMap.Metric[strings.ToLower(metric.ID)] = valueMetric
			return valueMetric, nil
		} else {
			return Metric{}, errors.New("data is empty")
		}
	} else {
		if metric.Delta != nil || metric.Value != nil {
			metricResourceMap.Metric[strings.ToLower(metric.ID)] = metric
			return metric, nil
		} else {
			return Metric{}, errors.New("data is empty")
		}
	}
}

func (m *MemoryStorage) UpdateMetricByParameters(metricResourceMap *MetricResourceMap, metricName string, metricType string, value string) (Metric, error) {
	metricResourceMap.Mutex.Lock()
	defer metricResourceMap.Mutex.Unlock()
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return Metric{}, errors.New("only numbers  params in request are allowed")
	}

	var valueMetric, found = metricResourceMap.Metric[metricName]
	if found {
		err := valueMetric.UpdateMetricResource(value)
		metricResourceMap.Metric[strings.ToLower(metricName)] = valueMetric
		if err != nil {
			return Metric{}, err
		}
		return valueMetric, nil
	} else {
		metric, err := m.NewMetric(value, metricType, metricName)
		if err != nil {
			return Metric{}, err
		}
		metricResourceMap.Metric[metricName] = metric
		return metric, nil
	}
}
