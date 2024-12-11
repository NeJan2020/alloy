// Package rabbitmq_exporter embeds https://github.com/kbudde/rabbitmq_exporter
package rabbitmq_exporter

import (
	"fmt"

	"github.com/go-kit/log"

	"github.com/grafana/alloy/internal/runtime/logging/level"
	"github.com/grafana/alloy/internal/static/integrations"
	integrations_v2 "github.com/grafana/alloy/internal/static/integrations/v2"
	"github.com/grafana/alloy/internal/static/integrations/v2/metricsutils"

	re "github.com/NeJan2020/rabbitmq_exporter"
)

func init() {
	integrations.RegisterIntegration(&Config{})
	integrations_v2.RegisterLegacy(&Config{}, integrations_v2.TypeMultiplex, metricsutils.NewNamedShim("rabbitmq"))
}

func New(logger log.Logger, c *Config) (integrations.Integration, error) {
	level.Debug(logger).Log("msg", "initializing rabbitmq_exporter", "config", c)

	cfg, err := initConfig(c)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize exporter, config is invalid: %w", err)
	}

	// setup logger
	logrusLogger := integrations.NewLogger(logger)
	re.SetLogger(logrusLogger)

	exporter := re.NewExporter(cfg)
	return integrations.NewCollectorIntegration(
		c.Name(),
		integrations.WithCollectors(exporter),
	), nil
}
