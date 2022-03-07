package storage

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricUpdateMetric(t *testing.T) {
	type fields struct {
		metricType   string
		metricName   string
		valueGauge   float64
		valueCounter int64
	}
	type args struct {
		value      string
		metricType string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   MetricResource
	}{
		{
			name: "Update counter",
			fields: fields{
				metricType:   "counter",
				metricName:   "testCounter",
				valueCounter: 12,
			},
			args: args{
				value:      "13",
				metricType: "counter",
			},
			want: MetricResource{
				Metric: &Metric{
					MType: "counter",
					Delta: sumInt(24, 1),
					Value: nil,
					ID:    "testCounter",
				},
				Updated: nil,
				Mutex:   sync.Mutex{},
			},
		},
		{
			name: "Update gauge",
			fields: fields{
				metricType: "gauge",
				metricName: "stackinuse",
				valueGauge: 12,
			},
			args: args{
				value:      "13",
				metricType: "gauge",
			},
			want: MetricResource{
				Metric: &Metric{
					MType: "gauge",
					Value: sumFloat(12, 1),
				},
				Updated: nil,
				Mutex:   sync.Mutex{},
			},
		},
		{
			name: "Update fail type",
			fields: fields{
				metricType: "gauge",
				metricName: "stackinuse",
				valueGauge: 12,
			},
			args: args{
				value:      "12",
				metricType: "failType",
			},
			want: MetricResource{
				Metric: &Metric{
					MType: "gauge",
					Value: sumFloat(11, 1),
				},
				Updated: nil,
				Mutex:   sync.Mutex{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &MetricResource{
				Metric: &Metric{
					MType: tt.fields.metricType,
					Value: &tt.fields.valueGauge,
					Delta: &tt.fields.valueCounter,
				},
				Updated: nil,
				Mutex:   sync.Mutex{},
			}

			got.UpdateMetricResource(tt.args.value)
			assert.Equal(t, got.GetValue(), tt.want.GetValue())
			assert.Equal(t, got.GetMetricType(), tt.want.GetMetricType())
		})
	}
}

func TestMetric_UpdateMetric(t *testing.T) {
	type fields struct {
		metricType   string
		valueGauge   float64
		valueCounter int64
	}
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    MetricResource
		wantErr bool
	}{
		{
			name: "positive test 1#",
			fields: fields{
				metricType:   "",
				valueGauge:   0,
				valueCounter: 0,
			},
			args: args{
				value: "",
			},
			want: MetricResource{
				Metric: &Metric{
					MType: "",
					Value: nil,
					Delta: nil,
				},
				Updated: nil,
				Mutex:   sync.Mutex{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MetricResource{
				Metric: &Metric{
					MType: tt.fields.metricType,
					Value: &tt.fields.valueGauge,
					Delta: &tt.fields.valueCounter,
				},
				Updated: nil,
				Mutex:   sync.Mutex{},
			}
			err := got.UpdateMetricResource(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateMetricResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, got.GetValue(), tt.want.GetValue())
			assert.Equal(t, got.GetMetricType(), tt.want.GetMetricType())
		})
	}
}

func TestMetric_GetValue(t *testing.T) {
	type fields struct {
		metricType   string
		valueGauge   float64
		valueCounter int64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "positive test #1",
			fields: fields{
				metricType:   "gauge",
				valueGauge:   1,
				valueCounter: 2,
			},
			want: "1.000",
		},
		{
			name: "positive test #2",
			fields: fields{
				metricType:   "counter",
				valueGauge:   1,
				valueCounter: 2,
			},
			want: "2",
		},
		{
			name: "positive test #3",
			fields: fields{
				metricType:   "ty",
				valueGauge:   1,
				valueCounter: 2,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := MetricResource{
				Metric: &Metric{
					MType: tt.fields.metricType,
					Value: &tt.fields.valueGauge,
					Delta: &tt.fields.valueCounter,
				},
				Updated: nil,
				Mutex:   sync.Mutex{},
			}
			if got := u.GetValue(); got != tt.want {
				t.Errorf("GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetric_GetMetricType(t *testing.T) {
	type fields struct {
		metricType   string
		valueGauge   float64
		valueCounter int64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "positive test #1",
			fields: fields{
				metricType:   "gauge",
				valueGauge:   1,
				valueCounter: 2,
			},
			want: "gauge",
		},
		{
			name: "positive test #2",
			fields: fields{
				metricType:   "counter",
				valueGauge:   1,
				valueCounter: 2,
			},
			want: "counter",
		},
		{
			name: "positive test #3",
			fields: fields{
				metricType:   "ty",
				valueGauge:   1,
				valueCounter: 2,
			},
			want: "ty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := MetricResource{
				Metric: &Metric{
					MType: tt.fields.metricType,
					Value: &tt.fields.valueGauge,
					Delta: &tt.fields.valueCounter,
				},
				Updated: nil,
				Mutex:   sync.Mutex{},
			}
			if got := u.GetMetricType(); got != tt.want {
				t.Errorf("GetMetricType() = %v, want %v", got, tt.want)
			}
		})
	}
}
