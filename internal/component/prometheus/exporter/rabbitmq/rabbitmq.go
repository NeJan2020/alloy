package rabbitmq

import (
	"github.com/grafana/alloy/internal/component"
	"github.com/grafana/alloy/internal/component/prometheus/exporter"
	"github.com/grafana/alloy/internal/featuregate"
	"github.com/grafana/alloy/internal/static/integrations"
	"github.com/grafana/alloy/internal/static/integrations/rabbitmq_exporter"
	"github.com/grafana/alloy/syntax/alloytypes"
	config_util "github.com/prometheus/common/config"
)

func init() {
	component.Register(component.Registration{
		Name:      "prometheus.exporter.rabbitmq",
		Stability: featuregate.StabilityGenerallyAvailable,
		Args:      Arguments{},
		Exports:   exporter.Exports{},

		Build: exporter.New(createExporter, "rabbitmq"),
	})
}

func createExporter(opts component.Options, args component.Arguments, defaultInstanceKey string) (integrations.Integration, string, error) {
	a := args.(Arguments)
	return integrations.NewIntegrationWithInstanceKey(opts.Logger, a.Convert(), defaultInstanceKey)
}

// DefaultArguments holds non-zero default options for Arguments when it is
// unmarshaled from Alloy.
var DefaultArguments = Arguments{
	IncludeExporterMetrics:   false,
	RabbitURL:                "http://127.0.0.1:15672",
	RabbitUsername:           "guest",
	RabbitPassword:           "guest",
	RabbitConnection:         "direct",
	OutputFormat:             "TTY", //JSON
	CAFile:                   "ca.pem",
	CertFile:                 "client-cert.pem",
	KeyFile:                  "client-key.pem",
	InsecureSkipVerify:       false,
	ExcludeMetrics:           []string{},
	SkipExchangesString:      "^$",
	IncludeExchangesString:   ".*",
	SkipQueuesString:         "^$",
	IncludeQueuesString:      ".*",
	SkipVHostString:          "^$",
	IncludeVHostString:       ".*",
	RabbitCapabilitiesString: "no_sort,bert",
	AlivenessVhost:           "/",
	EnabledExporters:         []string{"exchange", "node", "overview", "queue"},
	Timeout:                  30,
	MaxQueues:                0,
}

type Arguments struct {
	IncludeExporterMetrics bool `alloy:"include_exporter_metrics,attr,optional"`

	// exporter-specific config.
	//
	// The exporter binary config differs to this, but these
	// are the only fields that are relevant to the exporter struct.
	RabbitURL                string            `alloy:"rabbit_url,attr"`
	RabbitUsername           string            `alloy:"rabbit_user,attr,optional"`
	RabbitPassword           alloytypes.Secret `alloy:"rabbit_pass,attr,optional"`
	RabbitConnection         string            `alloy:"rabbit_connection,attr,optional"`
	OutputFormat             string            `alloy:"output_format,attr,optional"`
	CAFile                   string            `alloy:"ca_file,attr,optional"`
	CertFile                 string            `alloy:"cert_file,attr,optional"`
	KeyFile                  string            `alloy:"key_file,attr,optional"`
	InsecureSkipVerify       bool              `alloy:"insecure_skip_verify,attr,optional"`
	ExcludeMetrics           []string          `alloy:"exlude_metrics,attr,optional"`
	IncludeExchangesString   string            `alloy:"include_exchanges,attr,optional"`
	SkipExchangesString      string            `alloy:"skip_exchanges,attr,optional"`
	IncludeQueuesString      string            `alloy:"include_queues,attr,optional"`
	SkipQueuesString         string            `alloy:"skip_queues,attr,optional"`
	SkipVHostString          string            `alloy:"skip_vhost,attr,optional"`
	IncludeVHostString       string            `alloy:"include_vhost,attr,optional"`
	RabbitCapabilitiesString string            `alloy:"rabbit_capabilities,attr,optional"`
	AlivenessVhost           string            `alloy:"aliveness_vhost,attr,optional"`
	EnabledExporters         []string          `alloy:"enabled_exporters,attr,optional"`
	Timeout                  int               `alloy:"timeout,attr,optional"`
	MaxQueues                int               `alloy:"max_queues,attr,optional"`
}

// SetToDefault implements syntax.Defaulter.
func (a *Arguments) SetToDefault() {
	*a = DefaultArguments
}

// Validate implements syntax.Validator.
func (a *Arguments) Validate() error {
	return nil
}

func (a *Arguments) Convert() *rabbitmq_exporter.Config {
	return &rabbitmq_exporter.Config{
		IncludeExporterMetrics:   a.IncludeExporterMetrics,
		RabbitURL:                a.RabbitURL,
		RabbitUsername:           a.RabbitUsername,
		RabbitPassword:           config_util.Secret(a.RabbitPassword),
		RabbitConnection:         a.RabbitConnection,
		OutputFormat:             a.OutputFormat,
		CAFile:                   a.CAFile,
		CertFile:                 a.CertFile,
		KeyFile:                  a.KeyFile,
		InsecureSkipVerify:       a.InsecureSkipVerify,
		ExcludeMetrics:           a.ExcludeMetrics,
		IncludeExchangesString:   a.IncludeExchangesString,
		SkipExchangesString:      a.SkipExchangesString,
		IncludeQueuesString:      a.IncludeQueuesString,
		SkipQueuesString:         a.SkipQueuesString,
		SkipVHostString:          a.SkipVHostString,
		IncludeVHostString:       a.IncludeVHostString,
		RabbitCapabilitiesString: a.RabbitCapabilitiesString,
		AlivenessVhost:           a.AlivenessVhost,
		EnabledExporters:         a.EnabledExporters,
		Timeout:                  a.Timeout,
		MaxQueues:                a.MaxQueues,
	}
}
