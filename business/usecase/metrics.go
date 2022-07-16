package usecase

import (
	"mqtt-metrics/business/entity"
	"mqtt-metrics/pkg/logger"
	"mqtt-metrics/pkg/metrics"
)

type MetricsUseCase struct {
	cfg            *MetricsConfig
	broker         Broker
	metricsHandler MetricsHandler
	log            *logger.Zerolog
	metrics        map[string]*Metric
}

type MetricsConfig struct {
	Metrics []*Metric
}

type Metric struct {
	Name        string
	Type        string
	Description string
	Topic       string
	handler     interface{}
}

func NewMetricsUseCase(cfg *MetricsConfig, broker Broker, metricsHandler MetricsHandler, log *logger.Zerolog) (*MetricsUseCase, error) {
	uc := &MetricsUseCase{
		cfg:            cfg,
		broker:         broker,
		metricsHandler: metricsHandler,
		log:            log,
	}

	uc.initMetrics()
	uc.broker.SetConnectHandler(uc.OnConnect)

	return uc, uc.broker.Start()
}

func (uc *MetricsUseCase) initMetrics() {
	uc.metrics = make(map[string]*Metric, len(uc.cfg.Metrics))
	for _, m := range uc.cfg.Metrics {
		switch m.Type {
		case "counter":
			m.handler = uc.metricsHandler.NewCounter(&metrics.CounterMetric{Name: m.Name, Help: m.Description})
		case "gauge":
			m.handler = uc.metricsHandler.NewGauge(&metrics.GaugeMetric{Name: m.Name, Help: m.Description})
		}

		uc.metrics[m.Topic] = m
	}
}

func (uc *MetricsUseCase) MessageHandler(topic string, payload []byte) {
	if m, ok := uc.metrics[topic]; ok {
		uc.log.Debug().Msgf("%s: %s", topic, payload)
		uc.set(m, string(payload))
	}
}

func (uc *MetricsUseCase) OnConnect() {
	for _, m := range uc.cfg.Metrics {
		uc.broker.Subscribe(m.Topic, uc.MessageHandler)
	}
}

func (uc *MetricsUseCase) set(m *Metric, v string) {
	value, err := entity.ConvertMetricValue(v)
	if err != nil {
		uc.log.Error().Msgf("failed to convert metric value \"%s\": %v", v, err)
		return
	}

	switch m.Type {
	case "gauge":
		m.handler.(metrics.Gauge).Set(value)
	case "counter":
		m.handler.(metrics.Counter).Add(value)
	}
}
