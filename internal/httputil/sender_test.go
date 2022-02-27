package httputil

import (
	"context"
	"sync"
	"testing"

	"github.com/Stingsk/Go_3_lesson_1/internal/metrics"
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
			if err := RunSender(tt.args.ctx, tt.args.duration, tt.args.messages, tt.args.wg, "localhost:8080"); (err != nil) != tt.wantErr {
				t.Errorf("RunSender() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSend(t *testing.T) {
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
			send(tt.args.send, "localhost:8080")
		})
	}
}
