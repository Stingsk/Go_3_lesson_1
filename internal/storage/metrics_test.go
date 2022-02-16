package storage

import (
	"reflect"
	"testing"
)

func TestMetricUpdateMetric(t *testing.T) {
	type fields struct {
		metricType string
		metricName string
		value      string
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
				metricType: "counter",
				metricName: "testCounter",
				value:      "12",
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
				value:      "12",
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
				value:      "12",
			},
			args: args{
				value:      "",
				metricType: "failType",
			},
			want: Metric{
				metricType: "gauge",
				valueGauge: 12,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Metric{
				metricType: tt.fields.metricType,
			}
			if got, _ := u.UpdateMetric(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}
