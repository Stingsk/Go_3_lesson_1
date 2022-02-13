package storage

import (
	"reflect"
	"testing"
)

func TestMetric_GetMetricName(t *testing.T) {
}

func TestMetric_GetMetricType(t *testing.T) {
	type fields struct {
		metricType string
		metricName string
		value      string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test",
			fields: fields{
				metricType: "Type",
				metricName: "Name",
				value:      "12",
			},
			want: "Type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Metric{
				metricType: tt.fields.metricType,
				metricName: tt.fields.metricName,
				value:      tt.fields.value,
			}
			if got := u.GetMetricType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMetricType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetric_NewMetricString(t *testing.T) {
}

func TestMetric_UpdateMetric(t *testing.T) {
}

func TestNewMetricNameString(t *testing.T) {
	type args struct {
		metricName string
	}
	tests := []struct {
		name    string
		args    args
		want    MetricName
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMetricNameString(tt.args.metricName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMetricNameString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMetricNameString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMetricTypeString(t *testing.T) {
}
