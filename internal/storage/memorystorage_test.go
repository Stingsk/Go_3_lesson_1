package storage

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
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
		{
			name: "positive test 1#",
			args: args{
				value:      "10",
				metricType: "counter",
				name:       "counter",
			},
			want: Metric{
				ID:    "counter",
				MType: "counter",
				Delta: sumInt(9, 1),
				Value: nil,
			},
			wantErr: assert.NoError,
		},
		{
			name: "positive test 2#",
			args: args{
				value:      "10",
				metricType: "gauge",
				name:       "counter",
			},
			want: Metric{
				ID:    "counter",
				MType: "gauge",
				Delta: nil,
				Value: sumFloat(9, 1),
			},
			wantErr: assert.NoError,
		},
		{
			name: "negative test 1#",
			args: args{
				value:      "ert",
				metricType: "gauge",
				name:       "counter",
			},
			want: Metric{
				ID:    "",
				MType: "",
				Delta: nil,
				Value: nil,
			},
			wantErr: assert.Error,
		},
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
		{
			name: "positive test 1#",
			args: args{
				metricResourceMap: &MetricResourceMap{
					Metric: map[string]Metric{
						"gauge": Metric{
							ID:    "gauge",
							MType: "gauge",
							Delta: nil,
							Value: sumFloat(9, 1),
						},
					},
					Mutex:      sync.Mutex{},
					Repository: &MemoryStorage{},
				},
				metric: Metric{
					ID:    "gauge",
					MType: "gauge",
					Delta: nil,
					Value: sumFloat(19, 1),
				},
			},
			want: Metric{
				ID:    "gauge",
				MType: "gauge",
				Delta: nil,
				Value: sumFloat(19, 1),
			},
			wantErr: assert.NoError,
		},
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
