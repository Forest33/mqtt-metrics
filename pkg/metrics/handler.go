package metrics

import (
	"fmt"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/prometheus/client_golang/prometheus"
	prom "github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"

	"mqtt-metrics/pkg/logger"
)

const (
	defaultMetricsPath        = "/metrics"
	defaultMetricsNamespace   = "home"
	defaultServerHostname     = ""
	defaultLogLevel           = "debug"
	defaultLogTimeFieldFormat = time.RFC3339
	defaultServerPort         = 9701
)

type Metrics struct {
	cfg *Config
	log *logger.Zerolog
}

type Config struct {
	Path                 string
	Namespace            string
	ServerHostname       string
	ServerPort           int
	LogLevel             string
	LogTimeFieldFormat   string
	LogPrettyPrint       bool
	LogDisableSampling   bool
	LogRedirectStdLogger bool
	LogErrorStack        bool
	LogShowCaller        bool
}

func New(cfg *Config) (*Metrics, error) {
	cfg.normalize()

	useCase := &Metrics{
		cfg: cfg,
		log: logger.NewZerolog(logger.ZeroConfig{
			Level:             cfg.LogLevel,
			TimeFieldFormat:   cfg.LogTimeFieldFormat,
			PrettyPrint:       cfg.LogPrettyPrint,
			DisableSampling:   cfg.LogDisableSampling,
			RedirectStdLogger: cfg.LogRedirectStdLogger,
			ErrorStack:        cfg.LogErrorStack,
			ShowCaller:        cfg.LogShowCaller,
		}),
	}

	cb := make(chan error)
	defer close(cb)

	useCase.listenServer(cb)
	select {
	case err := <-cb:
		return nil, err
	case <-time.After(time.Second):
	}

	return useCase, nil
}

func (useCase *Metrics) listenServer(cb chan error) {
	go func() {
		useCase.log.Info().Msgf("metrics server at %s:%d", useCase.cfg.ServerHostname, useCase.cfg.ServerPort)
		router := fasthttprouter.New()
		router.GET(useCase.cfg.Path, fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler()))
		err := fasthttp.ListenAndServe(fmt.Sprintf("%s:%d", useCase.cfg.ServerHostname, useCase.cfg.ServerPort), router.Handler)
		if err != nil {
			cb <- err
		}
	}()
}

func (useCase *Metrics) NewCounter(m *CounterMetric) Counter {
	return prom.NewCounter(prometheus.CounterOpts{
		Namespace:   useCase.getNamespace(m.Namespace),
		Subsystem:   m.Subsystem,
		Name:        m.Name,
		Help:        m.Help,
		ConstLabels: m.ConstLabels,
	})
}

func (useCase *Metrics) NewGauge(m *GaugeMetric) Gauge {
	return prom.NewGauge(prometheus.GaugeOpts{
		Namespace:   useCase.getNamespace(m.Namespace),
		Subsystem:   m.Subsystem,
		Name:        m.Name,
		Help:        m.Help,
		ConstLabels: m.ConstLabels,
	})
}

func (useCase *Metrics) NewHistogram(m *HistogramMetric) Histogram {
	return prom.NewHistogram(prometheus.HistogramOpts{
		Namespace:   useCase.getNamespace(m.Namespace),
		Subsystem:   m.Subsystem,
		Name:        m.Name,
		Help:        m.Help,
		ConstLabels: m.ConstLabels,
		Buckets:     m.Buckets,
	})
}

func (useCase *Metrics) NewSummary(m *SummaryMetric) Summary {
	return prom.NewSummary(prometheus.SummaryOpts{
		Namespace:   useCase.getNamespace(m.Namespace),
		Subsystem:   m.Subsystem,
		Name:        m.Name,
		Help:        m.Help,
		ConstLabels: m.ConstLabels,
		Objectives:  m.Objectives,
		MaxAge:      m.MaxAge,
		AgeBuckets:  m.AgeBuckets,
		BufCap:      m.BufCap,
	})
}

func (useCase Metrics) getNamespace(namespace string) string {
	if len(namespace) > 0 {
		return namespace
	}
	return useCase.cfg.Namespace
}

func (cfg *Config) normalize() {
	if len(cfg.Path) == 0 {
		cfg.Path = defaultMetricsPath
	}
	if len(cfg.Namespace) == 0 {
		cfg.Namespace = defaultMetricsNamespace
	}
	if len(cfg.ServerHostname) == 0 {
		cfg.ServerHostname = defaultServerHostname
	}
	if cfg.ServerPort == 0 {
		cfg.ServerPort = defaultServerPort
	}
	if len(cfg.LogLevel) == 0 {
		cfg.LogLevel = defaultLogLevel
	}
	if len(cfg.LogTimeFieldFormat) == 0 {
		cfg.LogTimeFieldFormat = defaultLogTimeFieldFormat
	}
}
