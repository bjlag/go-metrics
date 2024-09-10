package collector

import (
	"math/rand"
	"runtime"
)

type MetricCollector struct {
	rtm *runtime.MemStats
}

func NewMetricCollector(rtm *runtime.MemStats) *MetricCollector {
	return &MetricCollector{
		rtm: rtm,
	}
}

func (c MetricCollector) ReadStats() {
	runtime.ReadMemStats(c.rtm)
}

func (c MetricCollector) Collect() []*Metric {
	return []*Metric{
		NewMetric(Gauge, "Alloc", c.rtm.Alloc),
		NewMetric(Gauge, "TotalAlloc", c.rtm.TotalAlloc),
		NewMetric(Gauge, "BuckHashSys", c.rtm.BuckHashSys),
		NewMetric(Gauge, "Frees", c.rtm.Frees),
		NewMetric(Gauge, "GCCPUFraction", c.rtm.GCCPUFraction),
		NewMetric(Gauge, "GCSys", c.rtm.GCSys),
		NewMetric(Gauge, "HeapAlloc", c.rtm.HeapAlloc),
		NewMetric(Gauge, "HeapIdle", c.rtm.HeapIdle),
		NewMetric(Gauge, "HeapInuse", c.rtm.HeapInuse),
		NewMetric(Gauge, "HeapObjects", c.rtm.HeapObjects),
		NewMetric(Gauge, "HeapReleased", c.rtm.HeapReleased),
		NewMetric(Gauge, "HeapSys", c.rtm.HeapSys),
		NewMetric(Gauge, "LastGC", c.rtm.LastGC),
		NewMetric(Gauge, "Lookups", c.rtm.Lookups),
		NewMetric(Gauge, "MCacheInuse", c.rtm.MCacheInuse),
		NewMetric(Gauge, "MCacheSys", c.rtm.MCacheSys),
		NewMetric(Gauge, "MSpanInuse", c.rtm.MSpanInuse),
		NewMetric(Gauge, "MSpanSys", c.rtm.MSpanSys),
		NewMetric(Gauge, "Mallocs", c.rtm.Mallocs),
		NewMetric(Gauge, "NextGC", c.rtm.NextGC),
		NewMetric(Gauge, "NumForcedGC", c.rtm.NumForcedGC),
		NewMetric(Gauge, "NumGC", c.rtm.NumGC),
		NewMetric(Gauge, "OtherSys", c.rtm.OtherSys),
		NewMetric(Gauge, "PauseTotalNs", c.rtm.PauseTotalNs),
		NewMetric(Gauge, "StackInuse", c.rtm.StackInuse),
		NewMetric(Gauge, "StackSys", c.rtm.StackSys),
		NewMetric(Gauge, "Sys", c.rtm.Sys),
		NewMetric(Gauge, "RandomValue", getRandomFloat(1, 100)),
	}
}

func getRandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
