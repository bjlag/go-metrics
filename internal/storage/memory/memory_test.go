package memory_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bjlag/go-metrics/internal/storage/memory"
)

func TestMemStorage_Counter(t *testing.T) {
	type args struct {
		name   string
		values []int64
	}

	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "success",
			args: args{
				name:   "test",
				values: []int64{1, 2, 3, 4, 5},
			},
			want: 15,
		},
		{
			name: "empty",
			args: args{
				name:   "test",
				values: []int64{},
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := memory.NewMemStorage()

			for _, value := range tt.args.values {
				s.AddCounter(tt.args.name, value)
			}

			assert.Equal(t, tt.want, s.GetCounter(tt.args.name))
		})
	}
}

func TestMemStorage_Gauge(t *testing.T) {
	type args struct {
		name   string
		values []float64
	}

	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "success",
			args: args{
				name:   "test",
				values: []float64{1, 2.2, 3.3, 0, 5},
			},
			want: 5,
		},
		{
			name: "empty",
			args: args{
				name:   "test",
				values: []float64{},
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := memory.NewMemStorage()

			for _, value := range tt.args.values {
				s.SetGauge(tt.args.name, value)
			}

			assert.Equal(t, tt.want, s.GetGauge(tt.args.name))
		})
	}
}
