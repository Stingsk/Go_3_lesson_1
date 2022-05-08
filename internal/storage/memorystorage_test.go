package storage

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage_GetMetric(t *testing.T) {
	type fields struct {
		Metric map[string]*Metric
		Mutex  sync.Mutex
	}
	type args struct {
		in0  context.Context
		name string
		in2  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Metric
		wantErr assert.ErrorAssertionFunc
	}{
		/*{
			name: "",
			fields: fields{
				Metric: nil,
				Mutex:  sync.Mutex{},
			},
			args: args{
				in0:  nil,
				name: "",
				in2:  "",
			},
			want: &Metric{
				ID:    "",
				MType: "",
				Delta: nil,
				Value: nil,
				Hash:  "",
			},
			wantErr: nil,
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryStorage{
				Metric: tt.fields.Metric,
				Mutex:  tt.fields.Mutex,
			}
			got, err := m.GetMetric(tt.args.in0, tt.args.name, tt.args.in2)
			if !tt.wantErr(t, err, fmt.Sprintf("GetMetric(%v, %v, %v)", tt.args.in0, tt.args.name, tt.args.in2)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetMetric(%v, %v, %v)", tt.args.in0, tt.args.name, tt.args.in2)
		})
	}
}

func TestMemoryStorage_GetMetrics(t *testing.T) {
	type fields struct {
		Metric map[string]*Metric
		Mutex  sync.Mutex
	}
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]*Metric
		wantErr assert.ErrorAssertionFunc
	}{
		/*{
			name: "",
			fields: fields{
				Metric: nil,
				Mutex:  sync.Mutex{},
			},
			args: args{
				in0: nil,
			},
			want:    nil,
			wantErr: nil,
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryStorage{
				Metric: tt.fields.Metric,
				Mutex:  tt.fields.Mutex,
			}
			got, err := m.GetMetrics(tt.args.in0)
			if !tt.wantErr(t, err, fmt.Sprintf("GetMetrics(%v)", tt.args.in0)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetMetrics(%v)", tt.args.in0)
		})
	}
}

func TestMemoryStorage_NewMetric(t *testing.T) {
	type fields struct {
		Metric map[string]*Metric
		Mutex  sync.Mutex
	}
	type args struct {
		value      string
		metricType string
		name       string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Metric
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "positive test 1#",
			fields: fields{
				Metric: nil,
				Mutex:  sync.Mutex{},
			},
			args: args{
				value:      "10",
				metricType: "counter",
				name:       "counter",
			},
			want: Metric{
				ID:    "counter",
				MType: "counter",
				Delta: sumInt(9, 1),
				Value: nil,
				Hash:  "",
			},
			wantErr: assert.NoError,
		},
		{
			name: "positive test 2#",
			fields: fields{
				Metric: nil,
				Mutex:  sync.Mutex{},
			},
			args: args{
				value:      "10",
				metricType: "gauge",
				name:       "counter",
			},
			want: Metric{
				ID:    "counter",
				MType: "gauge",
				Delta: nil,
				Value: sumFloat(9, 1),
				Hash:  "",
			},
			wantErr: assert.NoError,
		},
		{
			name: "negative test 1#",
			fields: fields{
				Metric: nil,
				Mutex:  sync.Mutex{},
			},
			args: args{
				value:      "ert",
				metricType: "gauge",
				name:       "counter",
			},
			want: Metric{
				ID:    "",
				MType: "",
				Delta: nil,
				Value: nil,
				Hash:  "",
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryStorage{
				Metric: tt.fields.Metric,
				Mutex:  tt.fields.Mutex,
			}
			got, err := m.NewMetric(tt.args.value, tt.args.metricType, tt.args.name)
			if !tt.wantErr(t, err, fmt.Sprintf("NewMetric(%v, %v, %v)", tt.args.value, tt.args.metricType, tt.args.name)) {
				return
			}
			assert.Equalf(t, tt.want, got, "NewMetric(%v, %v, %v)", tt.args.value, tt.args.metricType, tt.args.name)
		})
	}
}

func TestMemoryStorage_Ping(t *testing.T) {
	type fields struct {
		Metric map[string]*Metric
		Mutex  sync.Mutex
	}
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		/*{
			name: "",
			fields: fields{
				Metric: nil,
				Mutex:  sync.Mutex{},
			},
			args: args{
				in0: nil,
			},
			wantErr: nil,
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryStorage{
				Metric: tt.fields.Metric,
				Mutex:  tt.fields.Mutex,
			}
			tt.wantErr(t, m.Ping(tt.args.in0), fmt.Sprintf("Ping(%v)", tt.args.in0))
		})
	}
}

func TestMemoryStorage_UpdateMetric(t *testing.T) {
	type fields struct {
		Metric map[string]*Metric
		Mutex  sync.Mutex
	}
	type args struct {
		in0    context.Context
		metric Metric
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    assert.ErrorAssertionFunc
		wantMetric Metric
	}{
		{
			name: "positive test 1#",
			fields: fields{
				Metric: map[string]*Metric{
					"gauge": {
						ID:    "gauge",
						MType: "gauge",
						Delta: nil,
						Value: sumFloat(9, 1),
					},
				},
				Mutex: sync.Mutex{},
			},
			args: args{
				in0: nil,
				metric: Metric{
					ID:    "gauge",
					MType: "gauge",
					Delta: nil,
					Value: sumFloat(9, 1),
					Hash:  "",
				},
			},
			wantErr: assert.NoError,
			wantMetric: Metric{
				ID:    "gauge",
				MType: "gauge",
				Delta: nil,
				Value: sumFloat(9, 1),
				Hash:  "",
			},
		},
		{
			name: "positive test 2#",
			fields: fields{
				Metric: map[string]*Metric{
					"counter": {
						ID:    "counter",
						MType: "counter",
						Delta: sumInt(9, 1),
						Value: nil,
					},
				},
				Mutex: sync.Mutex{},
			},
			args: args{
				in0: nil,
				metric: Metric{
					ID:    "counter",
					MType: "counter",
					Delta: sumInt(9, 1),
					Value: nil,
					Hash:  "",
				},
			},
			wantErr: assert.NoError,
			wantMetric: Metric{
				ID:    "counter",
				MType: "counter",
				Delta: sumInt(19, 1),
				Value: nil,
				Hash:  "",
			},
		},
		{
			name: "negative test 1#",
			fields: fields{
				Metric: map[string]*Metric{
					"counter": {
						ID:    "counter",
						MType: "counter",
						Delta: sumInt(9, 1),
						Value: nil,
					},
				},
				Mutex: sync.Mutex{},
			},
			args: args{
				in0: nil,
				metric: Metric{
					ID:    "3sda",
					MType: "345",
					Delta: sumInt(9, 1),
					Value: nil,
					Hash:  "",
				},
			},
			wantErr: assert.Error,
			wantMetric: Metric{
				ID:    "counter",
				MType: "counter",
				Delta: sumInt(19, 1),
				Value: nil,
				Hash:  "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryStorage{
				Metric: tt.fields.Metric,
				Mutex:  tt.fields.Mutex,
			}
			ctx, _ := context.WithCancel(context.Background())
			tt.wantErr(t, m.UpdateMetric(ctx, tt.args.metric), fmt.Sprintf("UpdateMetric(%v, %v)", ctx, tt.wantMetric))
		})
	}
}

func TestMemoryStorage_UpdateMetricByParameters(t *testing.T) {
	type fields struct {
		Metric map[string]*Metric
		Mutex  sync.Mutex
	}
	type args struct {
		in0        context.Context
		metricName string
		metricType string
		value      string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "positive test 1#",
			fields: fields{
				Metric: map[string]*Metric{
					"gauge": {
						ID:    "gauge",
						MType: "gauge",
						Delta: nil,
						Value: sumFloat(9, 1),
					},
				},
				Mutex: sync.Mutex{},
			},
			args: args{
				in0:        nil,
				metricName: "gauge",
				metricType: "gauge",
				value:      "9",
			},
			wantErr: assert.NoError,
		},
		{
			name: "positive test 2#",
			fields: fields{
				Metric: map[string]*Metric{
					"gauge": {
						ID:    "counter",
						MType: "counter",
						Delta: sumInt(9, 1),
						Value: nil,
					},
				},
				Mutex: sync.Mutex{},
			},
			args: args{
				in0:        nil,
				metricName: "counter",
				metricType: "counter",
				value:      "9",
			},
			wantErr: assert.NoError,
		},
		{
			name: "negative test 1#",
			fields: fields{
				Metric: map[string]*Metric{
					"gauge": {
						ID:    "counter",
						MType: "counter",
						Delta: sumInt(9, 1),
						Value: nil,
					},
				},
				Mutex: sync.Mutex{},
			},
			args: args{
				in0:        nil,
				metricName: "counter",
				metricType: "counter",
				value:      "fdgf",
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryStorage{
				Metric: tt.fields.Metric,
				Mutex:  tt.fields.Mutex,
			}

			ctx, _ := context.WithCancel(context.Background())
			tt.wantErr(t, m.UpdateMetricByParameters(ctx, tt.args.metricName, tt.args.metricType, tt.args.value), fmt.Sprintf("UpdateMetricByParameters(%v, %v, %v, %v)", ctx, tt.args.metricName, tt.args.metricType, tt.args.value))
		})
	}
}

func TestMemoryStorage_UpdateMetrics(t *testing.T) {
	type fields struct {
		Metric map[string]*Metric
		Mutex  sync.Mutex
	}
	type args struct {
		in0          context.Context
		metricsBatch []*Metric
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "positive test 1#",
			fields: fields{
				Metric: map[string]*Metric{
					"gauge": {
						ID:    "gauge",
						MType: "gauge",
						Delta: nil,
						Value: sumFloat(9, 8),
					},
				},
				Mutex: sync.Mutex{},
			},
			args: args{
				in0: context.Background(),
				metricsBatch: []*Metric{
					{
						ID:    "gauge",
						MType: "gauge",
						Delta: nil,
						Value: sumFloat(9, 1),
					}, {
						ID:    "gauge1",
						MType: "gauge",
						Delta: nil,
						Value: sumFloat(9, 1),
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryStorage{
				Metric: tt.fields.Metric,
				Mutex:  tt.fields.Mutex,
			}

			ctx, _ := context.WithCancel(context.Background())
			tt.wantErr(t, m.UpdateMetrics(ctx, tt.args.metricsBatch), fmt.Sprintf("UpdateMetrics(%v, %v)", ctx, tt.args.metricsBatch))
		})
	}
}

func TestNewMemoryStorage(t *testing.T) {
	tests := []struct {
		name string
		want *MemoryStorage
	}{
		{
			name: "positive test 1#",
			want: &MemoryStorage{
				Metric: map[string]*Metric{},
				Mutex:  sync.Mutex{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewMemoryStorage(), "NewMemoryStorage()")
		})
	}
}
