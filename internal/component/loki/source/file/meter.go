package file

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/go-kit/log"
	"github.com/grafana/alloy/internal/runtime/logging/level"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

type OTLPSDK struct {
	meter api.Meter

	levelCounter     api.Int64Counter
	exceptionCounter api.Int64Counter

	listenAddress string
	srvMux        *http.ServeMux
}

const meterName = "github.com/CloudDetail/apo-alloy"

func InitMeter(log log.Logger, address string) *OTLPSDK {
	sdk := &OTLPSDK{
		srvMux:        http.NewServeMux(),
		listenAddress: address,
	}

	exporter, err := prometheus.New()
	if err != nil {
		level.Error(log).Log("msg", "PROM_EXPORTER_ERROR", "err", err.Error())
	}

	hostname, find := os.LookupEnv("_node_name_")
	if !find {
		hostname, _ = os.Hostname()
	}
	hostIP := os.Getenv("_node_ip_")
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.HostIPKey.String(hostIP),
			semconv.HostName(hostname),
		),
	)
	if err != nil {
		level.Error(log).Log("msg", "PROM_EXPORTER_ERROR", "err", err.Error())
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter), metric.WithResource(res))
	sdk.meter = provider.Meter(meterName)

	// TODO deal error
	sdk.levelCounter, _ = sdk.meter.Int64Counter(
		"originx_logparser_level_count",
		api.WithDescription("log level counter"))
	sdk.exceptionCounter, _ = sdk.meter.Int64Counter(
		"originx_logparser_exception_count",
		api.WithDescription("log exception counter"))

	level.Info(log).Log("msg", "PROM_EXPORTER_INIT", "addr", "serving metrics at "+address)

	sdk.srvMux.Handle("/metrics", promhttp.Handler())
	setUpHTTPServer(log, sdk.srvMux, address)

	return sdk
}

func setUpHTTPServer(log log.Logger, handler http.Handler, listenAddr string) (*http.Server, error) {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}

	server := &http.Server{Handler: handler}
	go func() {
		err := server.Serve(listener)
		if err != nil && err != http.ErrServerClosed {
			level.Error(log).Log("msg", "PROM_EXPORTER_ERROR", "err", err.Error())
		}
		listener.Close()
	}()
	return server, nil
}
