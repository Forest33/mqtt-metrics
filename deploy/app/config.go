package main

import (
	"encoding/json"
	"os"

	"mqtt-metrics/pkg/logger"
)

const (
	defaultConfigPath = "mqtt-metrics.json"
)

type Config struct {
	Broker     *BrokerConfig     `json:"broker"`
	Logger     *LoggerConfig     `json:"logger"`
	Prometheus *PrometheusConfig `json:"prometheus"`
	Metrics    []*MetricConfig   `json:"metrics"`
}

type LoggerConfig struct {
	Level             string `json:"level"`
	TimeFieldFormat   string `json:"time_field_format"`
	PrettyPrint       bool   `json:"pretty_print"`
	DisableSampling   bool   `json:"disable_sampling"`
	RedirectStdLogger bool   `json:"redirect_std_logger"`
	ErrorStack        bool   `json:"error_stack"`
	ShowCaller        bool   `json:"show_caller"`
}

type BrokerConfig struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	StateTopic string `json:"state_topic"`
	ClientID   string `json:"client_id"`
	UserName   string `json:"user_name"`
	Password   string `json:"password"`
}

type PrometheusConfig struct {
	Port      int    `json:"port"`
	Path      string `json:"path"`
	Namespace string `json:"namespace"`
}

type MetricConfig struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Topic       string `json:"topic"`
}

var (
	cfg = &Config{}
)

func init() {
	log := logger.NewDefaultZerolog()

	path, ok := os.LookupEnv("MQTT_METRICS_CONFIG")
	if !ok {
		path = defaultConfigPath
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		log.Fatal().Msg(err.Error())
	}
}
