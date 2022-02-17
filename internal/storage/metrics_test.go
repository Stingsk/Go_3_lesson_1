package storage

import (
	"reflect"
	"testing"
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
		want   Metric
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
			want: Metric{
				metricType:   "counter",
				valueCounter: 25,
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
			want: Metric{
				metricType: "gauge",
				valueGauge: 13,
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
				value:      "",
				metricType: "failType",
			},
			want: Metric{
				metricType: "",
				valueGauge: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Metric{
				metricType:   tt.fields.metricType,
				valueGauge:   tt.fields.valueGauge,
				valueCounter: tt.fields.valueCounter,
			}
			if got, _ := UpdateMetric(tt.args.value, *u); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateMetric() = %v, want %v", got, tt.want)
			}
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
		want    Metric
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
			want: Metric{
				metricType:   "",
				valueGauge:   0,
				valueCounter: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Metric{
				metricType:   tt.fields.metricType,
				valueGauge:   tt.fields.valueGauge,
				valueCounter: tt.fields.valueCounter,
			}
			got, err := UpdateMetric(tt.args.value, *u)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateMetric() got = %v, want %v", got, tt.want)
			}
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
			u := &Metric{
				metricType:   tt.fields.metricType,
				valueGauge:   tt.fields.valueGauge,
				valueCounter: tt.fields.valueCounter,
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
			u := &Metric{
				metricType:   tt.fields.metricType,
				valueGauge:   tt.fields.valueGauge,
				valueCounter: tt.fields.valueCounter,
			}
			if got := u.GetMetricType(); got != tt.want {
				t.Errorf("GetMetricType() = %v, want %v", got, tt.want)
			}
		})
	}
}
