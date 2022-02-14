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

func TestRunGetMetrics(t *testing.T) {
	type args struct {
		ctx      context.Context
		duration int
		messages *SensorData
		wg       *sync.WaitGroup
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
