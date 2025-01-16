package file_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bjlag/go-metrics/internal/storage/file"
)

func TestStorage_Load(t *testing.T) {
	// Arrange
	content := `[
  {
    "id": "PollCount",
    "type": "counter",
    "delta": 53
  },
  {
    "id": "Sys",
    "type": "gauge",
    "value": 17320976
  },
  {
    "id": "GCCPUFraction",
    "type": "gauge",
    "value": 0.00022545379332210888
  }
]`
	tmpfile, err := os.CreateTemp("", "metrics")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Remove(tmpfile.Name())
	}()

	if _, err = tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}

	// Act
	store, err := file.NewStorage(tmpfile.Name())
	assert.Nil(t, err)

	metrics, err := store.Load()
	assert.Nil(t, err)

	// Assert
	assert.Len(t, metrics, 3)
	assert.Contains(t, metrics, newCounter("PollCount", 53))
	assert.Contains(t, metrics, newGauge("Sys", 17320976))
	assert.Contains(t, metrics, newGauge("GCCPUFraction", 0.00022545379332210888))
}

func TestStorage_Save(t *testing.T) {
	// Arrange
	tmpfile, err := os.CreateTemp("", "metrics")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Remove(tmpfile.Name())
	}()

	// Act
	store, err := file.NewStorage(tmpfile.Name())
	assert.Nil(t, err)

	metrics := []file.Metric{
		newCounter("PollCount", 53),
		newGauge("Sys", 17320976),
		newGauge("GCCPUFraction", 0.00022545379332210888),
	}
	err = store.Save(metrics)
	assert.Nil(t, err)

	// Assert
	loadedMetrics, err := store.Load()
	assert.Nil(t, err)

	assert.Equal(t, metrics, loadedMetrics)
}

func newCounter(id string, value int64) file.Metric {
	return file.Metric{
		ID:    id,
		MType: "counter",
		Delta: &value,
	}
}

func newGauge(id string, value float64) file.Metric {
	return file.Metric{
		ID:    id,
		MType: "gauge",
		Value: &value,
	}
}
