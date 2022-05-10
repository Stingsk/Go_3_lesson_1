package storage

import (
	"context"
	"errors"
	"strconv"
	"sync"
)

type MemoryStorage struct {
	Metric map[string]*Metric
	Mutex  sync.Mutex
}

func NewMemoryStorage() *MemoryStorage {
	var m MemoryStorage

	m.Metric = make(map[string]*Metric)

	return &m
}

func (m *MemoryStorage) UpdateMetric(_ context.Context, metric Metric) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	var id = metric.ID
	var valueMetric = m.Metric[id]
	if valueMetric.GetValue() != "" {
		if metric.Delta != nil || metric.Value != nil {
			valueMetric.Update(metric)
			m.Metric[id] = valueMetric
			return nil
		} else {
			return errors.New("data is empty")
		}
	} else {
		if metric.Delta != nil || metric.Value != nil {
			m.Metric[id] = &metric
			return nil
		} else {
			return errors.New("data is empty")
		}
	}
}

func (m *MemoryStorage) UpdateMetricByParameters(_ context.Context, metricName string, metricType string, value string) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return errors.New("only numbers  params in request are allowed")
	}

	var valueMetric, found = m.Metric[metricName]
	if found {
		err := valueMetric.UpdateMetricResource(value)
		m.Metric[metricName] = valueMetric
		if err != nil {
			return err
		}
		return nil
	} else {
		_, err := m.NewMetric(value, metricType, metricName)
		if err != nil {
			return err
		}

		return nil
	}
}

func (m *MemoryStorage) GetMetric(_ context.Context, name string, _ string) (*Metric, error) {
	if name == "" {
		return &Metric{}, errors.New("empty name")
	}
	value, ok := m.Metric[name]

	if !ok {
		return &Metric{}, errors.New("name not found")
	}

	return value, nil
}
func (m *MemoryStorage) GetMetrics(_ context.Context) (map[string]*Metric, error) {
	return m.Metric, nil
}

func (m *MemoryStorage) UpdateMetrics(_ context.Context, metricsBatch []*Metric) error {
	for _, mb := range metricsBatch {
		m.Metric[mb.ID] = mb
	}

	return nil
}

func (m *MemoryStorage) Ping(_ context.Context) error {
	return nil
}

func (m *MemoryStorage) WriteMetrics() error {
	return nil
}

func (m *MemoryStorage) NewMetric(value string, metricType string, name string) (Metric, error) {
	metric, err := New(value, metricType, name)
	if err != nil {
		return Metric{}, err
	}
	m.Metric[name] = &metric
	return metric, nil
}
