package configuration

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/mrtroian/payservice/internal/gateway"
)

type Config struct {
	Key      string                   `yaml:"ssl_key_path"`
	Cert     string                   `yaml:"ssl_cert_path"`
	Endpoint string                   `yaml:"endpoint"`
	Host     string                   `yaml:"host"`
	Port     int                      `yaml:"port"`
	Gateways []gateway.PaymentGateway `yaml:"payment_providers"`
}

func ReadConfig(path string) (*Config, error) {
	rawConfig, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, errors.Wrap(err, "cannot open file")
	}
	config := new(Config)

	if err := yaml.Unmarshal(rawConfig, config); err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal")
	}

	return config, nil
}

func GetConfig() (*Config, error) {
	path := os.Getenv("PAYSERVICE_CONFIG_PATH")

	if len(path) <= 0 {
		return nil, errors.New("cannot read config from env")
	}

	conf, err := ReadConfig(path)

	if err != nil {
		return nil, errors.Wrap(err, "cannot read config:")
	}

	conf.Key = os.Getenv("SSL_KEY")
	conf.Cert = os.Getenv("SSL_CERT")

	return conf, nil
}
