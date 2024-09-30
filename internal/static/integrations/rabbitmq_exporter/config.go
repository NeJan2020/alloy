package rabbitmq_exporter

import (
	"fmt"
	"regexp"
	"strings"

	re "github.com/NeJan2020/rabbitmq_exporter"
	"github.com/go-kit/log"
	"github.com/grafana/alloy/internal/static/integrations"
)

var DefaultConfig = Config{
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

type Config struct {
	RabbitURL                string   `yaml:"rabbit_url"`
	RabbitUsername           string   `yaml:"rabbit_user"`
	RabbitPassword           string   `yaml:"rabbit_pass"`
	RabbitConnection         string   `yaml:"rabbit_connection"`
	OutputFormat             string   `yaml:"output_format"`
	CAFile                   string   `yaml:"ca_file"`
	CertFile                 string   `yaml:"cert_file"`
	KeyFile                  string   `yaml:"key_file"`
	InsecureSkipVerify       bool     `yaml:"insecure_skip_verify"`
	ExcludeMetrics           []string `yaml:"exlude_metrics"`
	IncludeExchangesString   string   `yaml:"include_exchanges"`
	SkipExchangesString      string   `yaml:"skip_exchanges"`
	IncludeQueuesString      string   `yaml:"include_queues"`
	SkipQueuesString         string   `yaml:"skip_queues"`
	SkipVHostString          string   `yaml:"skip_vhost"`
	IncludeVHostString       string   `yaml:"include_vhost"`
	RabbitCapabilitiesString string   `yaml:"rabbit_capabilities"`
	AlivenessVhost           string   `yaml:"aliveness_vhost"`
	EnabledExporters         []string `yaml:"enabled_exporters"`
	Timeout                  int      `yaml:"timeout"`
	MaxQueues                int      `yaml:"max_queues"`
}

// UnmarshalYAML 通过重写UnmarshalYAML设置默认值
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultConfig

	type plain Config
	return unmarshal((*plain)(c))
}

func (c *Config) Name() string {
	return "rabbitmq_exporter"
}

func (c *Config) InstanceKey(agentKey string) (string, error) {
	if c.RabbitURL == "" {
		c.RabbitURL = "http://127.0.0.1:15672"
	}

	return c.RabbitURL, nil
}

func (c *Config) NewIntegration(l log.Logger) (integrations.Integration, error) {
	return New(l, c)
}

func initConfig(c *Config) (*re.RabbitExporterConfig, error) {
	var config re.RabbitExporterConfig
	if valid, _ := regexp.MatchString("https?://[a-zA-Z.0-9]+", strings.ToLower(c.RabbitURL)); valid {
		config.RabbitURL = c.RabbitURL
	} else {
		return nil, fmt.Errorf("rabbit URL must start with http:// or https://")
	}

	if valid, _ := regexp.MatchString("(direct|loadbalancer)", c.RabbitConnection); valid {
		config.RabbitConnection = c.RabbitConnection
	} else {
		return nil, fmt.Errorf("rabbit connection must be direct or loadbalancer")
	}

	config.RabbitUsername = c.RabbitUsername
	config.RabbitPassword = c.RabbitPassword
	config.OutputFormat = c.OutputFormat
	config.CAFile = c.CAFile
	config.CertFile = c.CertFile
	config.KeyFile = c.KeyFile
	config.InsecureSkipVerify = c.InsecureSkipVerify
	config.ExcludeMetrics = c.ExcludeMetrics
	config.SkipExchanges = regexp.MustCompile(c.SkipExchangesString)
	config.IncludeExchanges = regexp.MustCompile(c.IncludeExchangesString)
	config.SkipQueues = regexp.MustCompile(c.SkipQueuesString)
	config.IncludeQueues = regexp.MustCompile(c.IncludeQueuesString)
	config.SkipVHost = regexp.MustCompile(c.SkipVHostString)
	config.IncludeVHost = regexp.MustCompile(c.IncludeVHostString)
	config.RabbitCapabilities = re.ParseCapabilities(c.RabbitCapabilitiesString)
	config.EnabledExporters = c.EnabledExporters
	config.AlivenessVhost = c.AlivenessVhost
	config.Timeout = c.Timeout
	config.MaxQueues = c.MaxQueues
	return &config, nil
}
