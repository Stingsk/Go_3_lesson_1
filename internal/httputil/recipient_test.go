package httputil

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/go-chi/chi/v5"
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
			RunServer(tt.args.wg, tt.args.sigChan)
		})
	}
}

func Test_getAllMetrics(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := concatenationMetrics(); got != tt.want {
				t.Errorf("concatenationMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_recipientGet(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		method string
		target string
		want   want
	}{
		{
			name:   "negative test #1",
			method: http.MethodGet,
			target: "/status",
			want: want{
				code:        404,
				response:    "Only 2 params in request are allowed!\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test #2",
			method: http.MethodGet,
			target: "/value/counter/testSetGet33",
			want: want{
				code:        404,
				response:    "Value NotFound!\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test #3",
			method: http.MethodGet,
			target: "/value/gauge/Alloc",
			want: want{
				code:        404,
				response:    "Value NotFound!\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.target, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := http.HandlerFunc(getMetric)
			// запускаем сервер
			h.ServeHTTP(w, request)
			res := w.Result()

			// проверяем код ответа
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(resBody) != tt.want.response {
				t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			}

			// заголовок ответа
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}

func Test_recipientGetAllMetrics(t *testing.T) {
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
			getAllMetrics(tt.args.w, tt.args.r)
		})
	}
}

func Test_recipientPost(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		method string
		target string
		want   want
	}{
		{
			name:   "negative test #1",
			method: http.MethodPost,
			target: "/status",
			want: want{
				code:        404,
				response:    "Only 3 params in request are allowed!\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test #2",
			method: http.MethodPost,
			target: "/update/asd/Alloc/345016",
			want: want{
				code:        501,
				response:    "MetricName NotImplemented!\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "negative test #3",
			method: http.MethodPost,
			target: "/update/gauge/Alloc/none",
			want: want{
				code:        400,
				response:    "Only Numbers  params in request are allowed!\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "positive test #1",
			method: http.MethodPost,
			target: "/update/gauge/Alloc/345016",
			want: want{
				code:        200,
				response:    "",
				contentType: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.target, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := http.HandlerFunc(postMetric)
			// запускаем сервер
			h.ServeHTTP(w, request)
			res := w.Result()

			// проверяем код ответа
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(resBody) != tt.want.response {
				t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			}

			// заголовок ответа
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
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
