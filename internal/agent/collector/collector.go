package collector

import (
	"math/rand"
	"runtime"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
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

func (c MetricCollector) Collect() ([]*Metric, error) {
	memStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	cpuCount, err := cpu.Counts(true)
	if err != nil {
		return nil, err
	}

	return []*Metric{
		NewGaugeMetric("Alloc", c.rtm.Alloc),
		NewGaugeMetric("TotalAlloc", c.rtm.TotalAlloc),
		NewGaugeMetric("BuckHashSys", c.rtm.BuckHashSys),
		NewGaugeMetric("Frees", c.rtm.Frees),
		NewGaugeMetric("GCCPUFraction", c.rtm.GCCPUFraction),
		NewGaugeMetric("GCSys", c.rtm.GCSys),
		NewGaugeMetric("HeapAlloc", c.rtm.HeapAlloc),
		NewGaugeMetric("HeapIdle", c.rtm.HeapIdle),
		NewGaugeMetric("HeapInuse", c.rtm.HeapInuse),
		NewGaugeMetric("HeapObjects", c.rtm.HeapObjects),
		NewGaugeMetric("HeapReleased", c.rtm.HeapReleased),
		NewGaugeMetric("HeapSys", c.rtm.HeapSys),
		NewGaugeMetric("LastGC", c.rtm.LastGC),
		NewGaugeMetric("Lookups", c.rtm.Lookups),
		NewGaugeMetric("MCacheInuse", c.rtm.MCacheInuse),
		NewGaugeMetric("MCacheSys", c.rtm.MCacheSys),
		NewGaugeMetric("MSpanInuse", c.rtm.MSpanInuse),
		NewGaugeMetric("MSpanSys", c.rtm.MSpanSys),
		NewGaugeMetric("Mallocs", c.rtm.Mallocs),
		NewGaugeMetric("NextGC", c.rtm.NextGC),
		NewGaugeMetric("NumForcedGC", c.rtm.NumForcedGC),
		NewGaugeMetric("NumGC", c.rtm.NumGC),
		NewGaugeMetric("OtherSys", c.rtm.OtherSys),
		NewGaugeMetric("PauseTotalNs", c.rtm.PauseTotalNs),
		NewGaugeMetric("StackInuse", c.rtm.StackInuse),
		NewGaugeMetric("StackSys", c.rtm.StackSys),
		NewGaugeMetric("Sys", c.rtm.Sys),
		NewGaugeMetric("FreeMemory", memStat.Free),
		NewGaugeMetric("TotalMemory", memStat.Total),
		NewGaugeMetric("CPUutilization1", cpuCount),
		NewGaugeMetric("RandomValue", getRandomFloat(1, 100)),
	}, nil
}

func getRandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
