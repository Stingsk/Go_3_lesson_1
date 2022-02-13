package httputil

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
	"reflect"
	"sync"
	"testing"
)

func TestRunRecipient(t *testing.T) {
	type args struct {
		wg      *sync.WaitGroup
		sigChan chan os.Signal
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunRecipient(tt.args.wg, tt.args.sigChan)
		})
	}
}

func Test_recipient(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recipient(tt.args.w, tt.args.r)
		})
	}
}

func Test_service(t *testing.T) {
	tests := []struct {
		name string
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := service(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setMiddlewares(t *testing.T) {
	type args struct {
		router *chi.Mux
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setMiddlewares(tt.args.router)
		})
	}
}
