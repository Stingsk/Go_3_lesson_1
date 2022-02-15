package storage

import (
	"reflect"
	"testing"
)

func TestMetricGetMetricType(t *testing.T) {
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

func TestNewMetricName(t *testing.T) {
	type args struct {
		metricName string
	}
	tests := []struct {
		name    string
		args    args
		want    MetricName
		wantErr bool
	}{
		{
			name:    "Add Real MetricName",
			args:    args{metricName: "stackinuse"},
			want:    MetricName{s: "stackinuse"},
			wantErr: false,
		},
		{
			name:    "Add Fail MetricName",
			args:    args{metricName: "MetricName"},
			want:    MetricName{s: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMetricName(tt.args.metricName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMetricName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMetricName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

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
				metricType: "counter",
				metricName: "testCounter",
				value:      "25",
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
				metricName: "stackinuse",
				value:      "13",
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
				value:      "13",
				metricType: "failType",
			},
			want: Metric{
				metricType: "gauge",
				metricName: "stackinuse",
				value:      "12",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Metric{
				metricType: tt.fields.metricType,
				metricName: tt.fields.metricName,
				value:      tt.fields.value,
			}
			if got := u.UpdateMetric(tt.args.value, tt.args.metricType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}
