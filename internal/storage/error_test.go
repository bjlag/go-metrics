package storage_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bjlag/go-metrics/internal/storage"
)

func TestNotFoundError(t *testing.T) {
	wrapErr := errors.New("some error")
	err := storage.NewMetricNotFoundError("counter", "pool", wrapErr)

	assert.Equal(t, "counter metric 'pool' not found", err.Error())
	assert.Equal(t, wrapErr, err.Unwrap())
	assert.Equal(t, "counter", err.Kind())
	assert.Equal(t, "pool", err.Name())
}
