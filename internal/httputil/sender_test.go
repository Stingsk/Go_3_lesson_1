package httputil

import (
	"context"
	"github.com/Stingsk/Go_3_lesson_1/internal/metrics"
	"sync"
	"testing"
)

func TestRunSender(t *testing.T) {
	type args struct {
		ctx      context.Context
		duration int
		messages *metrics.SensorData
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
			if err := RunSender(tt.args.ctx, tt.args.duration, tt.args.messages, tt.args.wg); (err != nil) != tt.wantErr {
				t.Errorf("RunSender() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_send(t *testing.T) {
	type args struct {
		send string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			send(tt.args.send)
		})
	}
}
