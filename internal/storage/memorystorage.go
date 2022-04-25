package storage

import (
	"errors"
	"strconv"
	"strings"
)

type MemoryStorage struct {
}

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

func (m *MemoryStorage) Ping() error {
	return nil
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
