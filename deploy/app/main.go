package main

import (
	"os"
	"os/signal"
	"syscall"

	"mqtt-metrics/adapter/broker"
	"mqtt-metrics/business/usecase"
	"mqtt-metrics/pkg/metrics"

	"mqtt-metrics/pkg/logger"
)

var (
	log            *logger.Zerolog
	brokerClient   *broker.Client
	metricsUseCase *usecase.MetricsUseCase
	metricsHandler *metrics.Metrics
)

func main() {
	defer shutdown()

	log = logger.NewZerolog(logger.ZeroConfig{
		Level:             cfg.Logger.Level,
		TimeFieldFormat:   cfg.Logger.TimeFieldFormat,
		PrettyPrint:       cfg.Logger.PrettyPrint,
		DisableSampling:   cfg.Logger.DisableSampling,
		RedirectStdLogger: cfg.Logger.RedirectStdLogger,
		ErrorStack:        cfg.Logger.ErrorStack,
		ShowCaller:        cfg.Logger.ShowCaller,
	})

	initAdapters()
	initUseCases()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit
}

func initAdapters() {
	var err error
	brokerClient, err = broker.NewBrokerClient(&broker.Config{
		Host:       cfg.Broker.Host,
		Port:       cfg.Broker.Port,
		StateTopic: cfg.Broker.StateTopic,
		ClientID:   cfg.Broker.ClientID,
		UserName:   cfg.Broker.UserName,
		Password:   cfg.Broker.Password,
	}, log)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	metricsHandler, err = metrics.New(&metrics.Config{
		Path:                 cfg.Prometheus.Path,
		Namespace:            cfg.Prometheus.Namespace,
		ServerPort:           cfg.Prometheus.Port,
		LogLevel:             cfg.Logger.Level,
		LogTimeFieldFormat:   cfg.Logger.TimeFieldFormat,
		LogPrettyPrint:       cfg.Logger.PrettyPrint,
		LogDisableSampling:   cfg.Logger.DisableSampling,
		LogRedirectStdLogger: cfg.Logger.RedirectStdLogger,
		LogErrorStack:        cfg.Logger.ErrorStack,
		LogShowCaller:        cfg.Logger.ShowCaller,
	})
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func initUseCases() {
	topicMetrics := make([]*usecase.Metric, 0, len(cfg.Metrics))
	for _, m := range cfg.Metrics {
		topicMetrics = append(topicMetrics, &usecase.Metric{
			Name:        m.Name,
			Type:        m.Type,
			Description: m.Description,
			Topic:       m.Topic,
		})
	}

	metricsCfg := &usecase.MetricsConfig{
		Metrics: topicMetrics,
	}

	var err error
	metricsUseCase, err = usecase.NewMetricsUseCase(metricsCfg, brokerClient, metricsHandler, log)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func shutdown() {
	if r := recover(); r != nil {
		log.Error().Msgf("panic: %v", r)
	}
	brokerClient.Close()
}
