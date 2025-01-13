package memory_test

import (
	"context"
	"github.com/bjlag/go-metrics/internal/storage"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bjlag/go-metrics/internal/storage/memory"
)

func TestStorage_Counter(t *testing.T) {
	type args struct {
		name string
	}

	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				name: "name",
			},
			want:    15,
			wantErr: false,
		},
		{
			name: "not found",
			args: args{
				name: "unknown",
			},
			want:    0,
			wantErr: true,
		},
	}

	s := memory.NewStorage()
	for _, value := range []int64{1, 2, 3, 4, 5} {
		s.AddCounter(context.Background(), "name", value)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curValue, err := s.GetCounter(context.Background(), tt.args.name)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.want, curValue)
			assert.Nil(t, err)
		})
	}
}

func TestStorage_GetAllCounters(t *testing.T) {
	s := memory.NewStorage()
	for _, value := range []int64{1, 2, 3, 4, 5} {
		s.AddCounter(context.Background(), "counter1", value)
		s.AddCounter(context.Background(), "counter2", value+1)
	}

	counters := s.GetAllCounters(context.Background())

	assert.Equal(t, 2, len(counters))
	assert.Equal(t, int64(15), counters["counter1"])
	assert.Equal(t, int64(20), counters["counter2"])
}

func TestStorage_SetGauge(t *testing.T) {
	type args struct {
		name string
	}

	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				name: "name",
			},
			want:    5,
			wantErr: false,
		},
		{
			name: "not found",
			args: args{
				name: "unknown",
			},
			want:    0,
			wantErr: true,
		},
	}

	s := memory.NewStorage()
	for _, value := range []float64{1, 2.2, 3.3, 0, 5} {
		s.SetGauge(context.Background(), "name", value)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curValue, err := s.GetGauge(context.Background(), tt.args.name)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.want, curValue)
			assert.Nil(t, err)
		})
	}
}

func TestStorage_GetAllGauges(t *testing.T) {
	s := memory.NewStorage()
	for _, value := range []float64{1, 2.2, 3.3, 0, 5} {
		s.SetGauge(context.Background(), "gauge1", value)
		s.SetGauge(context.Background(), "gauge2", value+2)
	}

	gauges := s.GetAllGauges(context.Background())

	assert.Equal(t, 2, len(gauges))
	assert.Equal(t, float64(5), gauges["gauge1"])
	assert.Equal(t, float64(7), gauges["gauge2"])
}

func TestStorage_SetGauges(t *testing.T) {
	s := memory.NewStorage()
	_ = s.SetGauges(context.Background(), []storage.Gauge{
		{
			ID:    "gauge1",
			Value: 1,
		},
		{
			ID:    "gauge2",
			Value: 5,
		},
		{
			ID:    "gauge1",
			Value: 3,
		},
	})

	g1, err := s.GetGauge(context.Background(), "gauge1")
	assert.Nil(t, err)
	g2, err := s.GetGauge(context.Background(), "gauge2")
	assert.Nil(t, err)

	assert.Equal(t, float64(3), g1)
	assert.Equal(t, float64(5), g2)
}

func TestStorage_Counters(t *testing.T) {
	s := memory.NewStorage()
	_ = s.AddCounters(context.Background(), []storage.Counter{
		{
			ID:    "counter1",
			Value: 1,
		},
		{
			ID:    "counter2",
			Value: 5,
		},
		{
			ID:    "counter1",
			Value: 3,
		},
	})

	c1, err := s.GetCounter(context.Background(), "counter1")
	assert.Nil(t, err)
	c2, err := s.GetCounter(context.Background(), "counter2")
	assert.Nil(t, err)

	assert.Equal(t, int64(4), c1)
	assert.Equal(t, int64(5), c2)
}
