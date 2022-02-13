package storage

import (
	"reflect"
	"testing"
)

func TestMetric_GetMetricName(t *testing.T) {
	type fields struct {
		metricType MetricType
		metricName MetricName
		counter    int
		value      string
	}
	tests := []struct {
		name   string
		fields fields
		want   MetricName
	}{
		{
			name: "Test",
			fields: fields{
				metricType: MetricType{"Type"},
				metricName: MetricName{"Name"},
				counter:    0,
				value:      "12",
			},
			want: MetricName{"Name"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Metric{
				metricType: tt.fields.metricType,
				metricName: tt.fields.metricName,
				counter:    tt.fields.counter,
				value:      tt.fields.value,
			}
			if got := u.GetMetricName(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMetricName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetric_GetMetricType(t *testing.T) {
	type fields struct {
		metricType MetricType
		metricName MetricName
		counter    int
		value      string
	}
	tests := []struct {
		name   string
		fields fields
		want   MetricType
	}{
		{
			name: "Test",
			fields: fields{
				metricType: MetricType{"Type"},
				metricName: MetricName{"Name"},
				counter:    0,
				value:      "12",
			},
			want: MetricType{"Type"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Metric{
				metricType: tt.fields.metricType,
				metricName: tt.fields.metricName,
				counter:    tt.fields.counter,
				value:      tt.fields.value,
			}
			if got := u.GetMetricType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMetricType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetric_NewMetricString(t *testing.T) {
	type fields struct {
		metricType MetricType
		metricName MetricName
		counter    int
		value      string
	}
	type args struct {
		metricName string
		metricType string
		value      string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test",
			fields: fields{
				metricType: MetricType{"Type"},
				metricName: MetricName{"Name"},
				counter:    0,
				value:      "12",
			},
			args: args{
				metricName: "Name",
				metricType: "Type",
				value:      "12",
			},
			wantErr: true,
		},
		{
			name: "Test2",
			fields: fields{
				metricType: MetricType{"counter"},
				metricName: MetricName{"pullcounter"},
				counter:    0,
				value:      "12",
			},
			args: args{
				metricName: "pullcounter",
				metricType: "counter",
				value:      "12",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Metric{
				metricType: tt.fields.metricType,
				metricName: tt.fields.metricName,
				counter:    tt.fields.counter,
				value:      tt.fields.value,
			}
			if err := u.NewMetricString(tt.args.metricName, tt.args.metricType, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("NewMetricString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMetric_UpdateMetric(t *testing.T) {
	type fields struct {
		metricType MetricType
		metricName MetricName
		counter    int
		value      string
	}
	type args struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Metric
	}{
		{
			name: "",
			fields: fields{
				metricType: MetricType{"counter"},
				metricName: MetricName{"pullcounter"},
				counter:    0,
				value:      "",
			},
			args: args{value: "5"},
			want: Metric{
				metricType: MetricType{"counter"},
				metricName: MetricName{"pullcounter"},
				counter:    1,
				value:      "5",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Metric{
				metricType: tt.fields.metricType,
				metricName: tt.fields.metricName,
				counter:    tt.fields.counter,
				value:      tt.fields.value,
			}
			if got := u.UpdateMetric(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateMetric() = %v, want %v", got, tt.want)
			}
		})
	}
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
	type args struct {
		metricType string
	}
	tests := []struct {
		name    string
		args    args
		want    MetricType
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMetricTypeString(tt.args.metricType)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMetricTypeString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMetricTypeString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
