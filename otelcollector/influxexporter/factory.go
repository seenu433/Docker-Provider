// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package influxexporter

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	// The value of "type" key in configuration.
	typeStr                   = "influx"
	defaultSamplingInitial    = 1 //= 2
	defaultSamplingThereafter = 1 //500
	measurementName			  = "prom2mdm1"
)

// NewFactory creates a factory for Logging exporter
func NewFactory() component.ExporterFactory {
	return exporterhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		exporterhelper.WithTraces(createTraceExporter),
		exporterhelper.WithMetrics(createMetricsExporter),
		exporterhelper.WithLogs(createLogsExporter))
}

func createDefaultConfig() configmodels.Exporter {
	return &Config{
		ExporterSettings: configmodels.ExporterSettings{
			TypeVal: typeStr,
			NameVal: typeStr,
		},
		LogLevel:           "info",
		SamplingInitial:    defaultSamplingInitial,
		SamplingThereafter: defaultSamplingThereafter,
	}
}

func createTraceExporter(_ context.Context, _ component.ExporterCreateParams, config configmodels.Exporter) (component.TracesExporter, error) {
	cfg := config.(*Config)

	exporterLogger, err := createLogger(cfg)
	if err != nil {
		return nil, err
	}

	return newTraceExporter(config, cfg.LogLevel, exporterLogger)
}

func createMetricsExporter(_ context.Context, _ component.ExporterCreateParams, config configmodels.Exporter) (component.MetricsExporter, error) {
	cfg := config.(*Config)

	exporterLogger, err := createLogger(cfg)
	if err != nil {
		return nil, err
	}

	return newMetricsExporter(config, cfg.LogLevel, exporterLogger)
}

func createLogsExporter(_ context.Context, _ component.ExporterCreateParams, config configmodels.Exporter) (component.LogsExporter, error) {
	cfg := config.(*Config)

	exporterLogger, err := createLogger(cfg)
	if err != nil {
		return nil, err
	}

	return newLogsExporter(config, cfg.LogLevel, exporterLogger)
}

func createLogger(cfg *Config) (*zap.Logger, error) {
	var level zapcore.Level
	err := (&level).UnmarshalText([]byte(cfg.LogLevel))
	if err != nil {
		return nil, err
	}

	// We take development config as the base since it matches the purpose
	// of logging exporter being used for debugging reasons (so e.g. console encoder)
	conf := zap.NewDevelopmentConfig()
	conf.Level = zap.NewAtomicLevelAt(level)
	conf.Sampling = &zap.SamplingConfig{
		Initial:    cfg.SamplingInitial,
		Thereafter: cfg.SamplingThereafter,
	}

	logginglogger, err := conf.Build()
	if err != nil {
		return nil, err
	}
	return logginglogger, nil
}