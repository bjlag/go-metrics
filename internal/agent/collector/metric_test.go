package collector_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bjlag/go-metrics/internal/agent/collector"
)

func TestMetric_Kind(t *testing.T) {
	type fields struct {
		kind  string
		name  string
		value interface{}
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "integer",
			fields: fields{
				kind:  "mType",
				name:  "name",
				value: 1,
			},
		},
		{
			name: "float",
			fields: fields{
				kind:  "mType",
				name:  "name",
				value: 1.1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := collector.NewMetric(tt.fields.kind, tt.fields.name, tt.fields.value)

			assert.Equal(t, m.Kind(), tt.fields.kind)
			assert.Equal(t, m.Name(), tt.fields.name)
			assert.Equal(t, m.Value(), tt.fields.value)
		})
	}
}
