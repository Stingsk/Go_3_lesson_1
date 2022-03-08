package httputil

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Stingsk/Go_3_lesson_1/internal/metrics"
)

func TestRunSender(t *testing.T) {
	type args struct {
		ctx      context.Context
		duration time.Duration
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
			RunSender(tt.args.ctx, tt.args.duration, tt.args.messages, tt.args.wg, "localhost:8080")
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
