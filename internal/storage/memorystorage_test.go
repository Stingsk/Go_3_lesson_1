package storage

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoryStorage_NewMetric(t *testing.T) {
	type args struct {
		value      string
		metricType string
		name       string
	}
	tests := []struct {
		name    string
		args    args
		want    Metric
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryStorage{}
			got, err := m.NewMetric(tt.args.value, tt.args.metricType, tt.args.name)
			if !tt.wantErr(t, err, fmt.Sprintf("NewMetric(%v, %v, %v)", tt.args.value, tt.args.metricType, tt.args.name)) {
				return
			}
			assert.Equalf(t, tt.want, got, "NewMetric(%v, %v, %v)", tt.args.value, tt.args.metricType, tt.args.name)
		})
	}
}

func TestMemoryStorage_UpdateMetric(t *testing.T) {
	type args struct {
		metricResourceMap *MetricResourceMap
		metric            Metric
	}
	tests := []struct {
		name    string
		args    args
		want    Metric
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryStorage{}
			got, err := m.UpdateMetric(tt.args.metricResourceMap, tt.args.metric)
			if !tt.wantErr(t, err, fmt.Sprintf("UpdateMetric(%v, %v)", tt.args.metricResourceMap, tt.args.metric)) {
				return
			}
			assert.Equalf(t, tt.want, got, "UpdateMetric(%v, %v)", tt.args.metricResourceMap, tt.args.metric)
		})
	}
}

func TestMemoryStorage_UpdateMetricByParameters(t *testing.T) {
	type args struct {
		metricResourceMap *MetricResourceMap
		metricName        string
		metricType        string
		value             string
	}
	tests := []struct {
		name    string
		args    args
		want    Metric
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryStorage{}
			got, err := m.UpdateMetricByParameters(tt.args.metricResourceMap, tt.args.metricName, tt.args.metricType, tt.args.value)
			if !tt.wantErr(t, err, fmt.Sprintf("UpdateMetricByParameters(%v, %v, %v, %v)", tt.args.metricResourceMap, tt.args.metricName, tt.args.metricType, tt.args.value)) {
				return
			}
			assert.Equalf(t, tt.want, got, "UpdateMetricByParameters(%v, %v, %v, %v)", tt.args.metricResourceMap, tt.args.metricName, tt.args.metricType, tt.args.value)
		})
	}
}
