package storage

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

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
	if m == nil {
		return ""
	}
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
func (m *Metric) GetHash(key string) string {
	var metricString string
	switch {
	case m.MType == MetricTypeGauge:
		metricString = fmt.Sprintf("%s:%s:%f", m.ID, MetricTypeGauge, *(m.Value))
	case m.MType == MetricTypeCounter:
		metricString = fmt.Sprintf("%s:%s:%d", m.ID, MetricTypeCounter, *(m.Delta))
	}

	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(metricString))

	return hex.EncodeToString(h.Sum(nil))
}

func (m *Metric) IsHashValid(key string) bool {
	if key == "" {
		return true
	}

	return m.Hash == m.GetHash(key)
}

func (m *Metric) SetHash(key string) {
	if key != "" {
		m.Hash = m.GetHash(key)
	}
}

func sumInt(first int64, second int64) *int64 {
	helper := first + second
	return &helper
}

func sumFloat(first float64, second float64) *float64 {
	helper := first + second
	return &helper
}
