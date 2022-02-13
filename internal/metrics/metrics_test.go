package metrics

import (
	"context"
	"reflect"
	"sync"
	"testing"
)

func TestGetNames(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunGetMetrics(t *testing.T) {
	type args struct {
		ctx      context.Context
		duration int
		messages *SensorData
		wg       *sync.WaitGroup
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RunGetMetrics(tt.args.ctx, tt.args.duration, tt.args.messages, tt.args.wg); (err != nil) != tt.wantErr {
				t.Errorf("RunGetMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSensorData_Get(t *testing.T) {
}

func TestSensorData_Store(t *testing.T) {
}

func Test_getMetrics(t *testing.T) {
	type args struct {
		count int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMetrics(tt.args.count); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newMonitor(t *testing.T) {
	type args struct {
		count int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "test", args: args{count: 2}, want: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newMonitor(tt.args.count)
			if !reflect.DeepEqual(got.PollCount, tt.want) {
				t.Errorf("newMonitor() got = %v, want %v", got, tt.want)
			}
		})
	}
}
