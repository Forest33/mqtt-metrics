package metrics

import (
	"time"
)

type Counter interface {
	Inc()
	Add(float64)
}

type Gauge interface {
	Set(float64)
	Inc()
	Dec()
	Add(float64)
	Sub(float64)
	SetToCurrentTime()
}

type Histogram interface {
	Observe(float64)
}

type Summary interface {
	Observe(float64)
}

type CounterMetric struct {
	Namespace   string
	Subsystem   string
	Name        string
	Help        string
	ConstLabels map[string]string
}

type GaugeMetric struct {
	Namespace   string
	Subsystem   string
	Name        string
	Help        string
	ConstLabels map[string]string
}

type HistogramMetric struct {
	Namespace   string
	Subsystem   string
	Name        string
	Help        string
	ConstLabels map[string]string
	Buckets     []float64
}

type SummaryMetric struct {
	Namespace   string
	Subsystem   string
	Name        string
	Help        string
	ConstLabels map[string]string
	Objectives  map[float64]float64
	MaxAge      time.Duration
	AgeBuckets  uint32
	BufCap      uint32
}
