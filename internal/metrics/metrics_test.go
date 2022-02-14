package metrics

import (
	"context"
	"reflect"
	"sync"
	"testing"
)

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

func TestSensorData_Store(t *testing.T) {
	type fields struct {
		mu   sync.RWMutex
		last []string
	}
	type args struct {
		data []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "",
			fields: fields{
				mu:   sync.RWMutex{},
				last: []string{"1", "2"},
			},
			args: args{
				data: []string{"1", "3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &SensorData{
				mu:   tt.fields.mu,
				last: tt.fields.last,
			}
			d.Store(tt.args.data)
		})
	}
}

func TestSensorData_Get(t *testing.T) {
	type fields struct {
		mu   sync.RWMutex
		last []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "Test",
			fields: fields{
				mu:   sync.RWMutex{},
				last: []string{"1", "2"},
			},
			want: []string{"1", "2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &SensorData{
				mu:   tt.fields.mu,
				last: tt.fields.last,
			}
			if got := d.Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
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

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test Fail",
			args: args{
				ctx:      ctx,
				duration: 30,
				messages: &SensorData{
					mu:   sync.RWMutex{},
					last: nil,
				},
				wg: wg,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cancel()
			if err := RunGetMetrics(tt.args.ctx, tt.args.duration, tt.args.messages, tt.args.wg); (err != nil) != tt.wantErr {
				t.Errorf("RunGetMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
