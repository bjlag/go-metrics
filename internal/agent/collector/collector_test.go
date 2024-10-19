package collector_test

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bjlag/go-metrics/internal/agent/collector"
)

func TestMetricCollector_Collect(t *testing.T) {
	rtm := &runtime.MemStats{
		Alloc:         1,
		TotalAlloc:    2,
		BuckHashSys:   3,
		Frees:         4,
		GCCPUFraction: 5,
		GCSys:         6,
		HeapAlloc:     7,
		HeapIdle:      8,
		HeapInuse:     9,
		HeapObjects:   10,
		HeapReleased:  11,
		HeapSys:       12,
		LastGC:        13,
		Lookups:       14,
		MCacheInuse:   15,
		MCacheSys:     16,
		MSpanInuse:    17,
		MSpanSys:      18,
		Mallocs:       19,
		NextGC:        20,
		NumForcedGC:   21,
		NumGC:         22,
		OtherSys:      23,
		PauseTotalNs:  24,
		StackInuse:    25,
		StackSys:      26,
		Sys:           27,
	}

	c := collector.NewMetricCollector(rtm)
	metrics, err := c.Collect()
	assert.NoError(t, err)

	assert.Equal(t, collector.NewMetric("gauge", "Alloc", uint64(1)), metrics[0])
	assert.Equal(t, collector.NewMetric("gauge", "TotalAlloc", uint64(2)), metrics[1])
	assert.Equal(t, collector.NewMetric("gauge", "BuckHashSys", uint64(3)), metrics[2])
	assert.Equal(t, collector.NewMetric("gauge", "Frees", uint64(4)), metrics[3])
	assert.Equal(t, collector.NewMetric("gauge", "GCCPUFraction", float64(5)), metrics[4])
	assert.Equal(t, collector.NewMetric("gauge", "GCSys", uint64(6)), metrics[5])
	assert.Equal(t, collector.NewMetric("gauge", "HeapAlloc", uint64(7)), metrics[6])
	assert.Equal(t, collector.NewMetric("gauge", "HeapIdle", uint64(8)), metrics[7])
	assert.Equal(t, collector.NewMetric("gauge", "HeapInuse", uint64(9)), metrics[8])
	assert.Equal(t, collector.NewMetric("gauge", "HeapObjects", uint64(10)), metrics[9])
	assert.Equal(t, collector.NewMetric("gauge", "HeapReleased", uint64(11)), metrics[10])
	assert.Equal(t, collector.NewMetric("gauge", "HeapSys", uint64(12)), metrics[11])
	assert.Equal(t, collector.NewMetric("gauge", "LastGC", uint64(13)), metrics[12])
	assert.Equal(t, collector.NewMetric("gauge", "Lookups", uint64(14)), metrics[13])
	assert.Equal(t, collector.NewMetric("gauge", "MCacheInuse", uint64(15)), metrics[14])
	assert.Equal(t, collector.NewMetric("gauge", "MCacheSys", uint64(16)), metrics[15])
	assert.Equal(t, collector.NewMetric("gauge", "MSpanInuse", uint64(17)), metrics[16])
	assert.Equal(t, collector.NewMetric("gauge", "MSpanSys", uint64(18)), metrics[17])
	assert.Equal(t, collector.NewMetric("gauge", "Mallocs", uint64(19)), metrics[18])
	assert.Equal(t, collector.NewMetric("gauge", "NextGC", uint64(20)), metrics[19])
	assert.Equal(t, collector.NewMetric("gauge", "NumForcedGC", uint32(21)), metrics[20])
	assert.Equal(t, collector.NewMetric("gauge", "NumGC", uint32(22)), metrics[21])
	assert.Equal(t, collector.NewMetric("gauge", "OtherSys", uint64(23)), metrics[22])
	assert.Equal(t, collector.NewMetric("gauge", "PauseTotalNs", uint64(24)), metrics[23])
	assert.Equal(t, collector.NewMetric("gauge", "StackInuse", uint64(25)), metrics[24])
	assert.Equal(t, collector.NewMetric("gauge", "StackSys", uint64(26)), metrics[25])
	assert.Equal(t, collector.NewMetric("gauge", "Sys", uint64(27)), metrics[26])
}
