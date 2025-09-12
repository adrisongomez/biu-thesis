package config

import (
	"github.com/spf13/viper"
)

type TelemetryConfig struct {
	ServiceName        string `mapstructure:"SERVICE_NAME"`
	ServiceVersion     string `mapstructure:"SERVICE_VERSION"`
	ServiceEnv         string `mapstructure:"ENVIRONMENT"`
	TracerCollectorURL string `mapstructure:"OTEL_TRACER_COLLECTOR_URL"`
	LoggerCollectorURL string `mapstructure:"OTEL_LOGGER_COLLECTOR_URL"`
	MetricCollectorURL string `mapstructure:"OTEL_METRIC_COLLECTOR_URL"`
}

func NewTelemetryConfig(pathToConfig string) (*TelemetryConfig, error) {
	viper.SetConfigFile(pathToConfig)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	cfg := &TelemetryConfig{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
