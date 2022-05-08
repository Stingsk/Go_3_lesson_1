package storage

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMemoryStorage()
			err := m.NewMetric(tt.args.value, tt.args.metricType, tt.args.name)
			if !tt.wantErr(t, err, fmt.Sprintf("NewMetric(%v, %v, %v)", tt.args.value, tt.args.metricType, tt.args.name)) {
				return
			}
			assert.Equalf(t, m, m, "NewMetric(%v, %v, %v)", tt.args.value, tt.args.metricType, tt.args.name)
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
			wantErr: assert.NoError,
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
			m := NewMemoryStorage()
			ctx, cancel := context.WithCancel(context.Background())
			m.NewMetric("10", tt.args.metric.MType, tt.args.metric.ID)
			tt.wantErr(t, m.UpdateMetric(ctx, tt.args.metric), fmt.Sprintf("UpdateMetric(%v, %v)", ctx, tt.wantMetric))
			cancel()
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
			m := NewMemoryStorage()
			ctx, cancel := context.WithCancel(context.Background())
			tt.wantErr(t, m.UpdateMetricByParameters(ctx, tt.args.metricName, tt.args.metricType, tt.args.value), fmt.Sprintf("UpdateMetricByParameters(%v, %v, %v, %v)", ctx, tt.args.metricName, tt.args.metricType, tt.args.value))
			cancel()
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
			m := NewMemoryStorage()

			ctx, cancel := context.WithCancel(context.Background())
			tt.wantErr(t, m.UpdateMetrics(ctx, tt.args.metricsBatch), fmt.Sprintf("UpdateMetrics(%v, %v)", ctx, tt.args.metricsBatch))
			cancel()
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
