package metrics

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNewMonitor(t *testing.T) {

}

func TestRunGetMetrics(t *testing.T) {
	type args struct {
		ctx      context.Context
		duration time.Duration
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
			RunGetMetrics(tt.args.ctx, tt.args.duration, tt.args.messages, tt.args.wg)
		})
	}
}
