package memory_test

import (
	"context"
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
