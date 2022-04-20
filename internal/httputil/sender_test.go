package httputil

import (
	"testing"
)

func TestSend(t *testing.T) {
	type args struct {
		send string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive test 1#",
			args: args{
				send: "dasdasd",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			send(tt.args.send, "localhost:8080")
		})
	}
}
