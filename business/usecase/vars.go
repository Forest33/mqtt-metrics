package usecase

import (
	"mqtt-metrics/adapter/broker"
	"mqtt-metrics/pkg/metrics"
)

type Broker interface {
	Start() error
	Publish(topic string, data []byte)
	Subscribe(topic string, handler broker.MessageHandler)
	SetConnectHandler(h broker.ConnectHandler)
	SetDisconnectHandler(h broker.DisconnectHandler)
}

type MetricsHandler interface {
	NewCounter(m *metrics.CounterMetric) metrics.Counter
	NewGauge(m *metrics.GaugeMetric) metrics.Gauge
	NewHistogram(m *metrics.HistogramMetric) metrics.Histogram
	NewSummary(m *metrics.SummaryMetric) metrics.Summary
}
