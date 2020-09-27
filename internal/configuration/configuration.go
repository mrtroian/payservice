package configuration

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

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
	var (
		conf *Config
		err  error
	)
	path := os.Getenv("PAYSERVICE_CONFIGS_DIR")

	if len(path) <= 0 {
		return nil, errors.New("cannot read 'PAYSERVICE_CONFIGS_DIR' from env")
	}

	files, err := ioutil.ReadDir(path)

	if err != nil {
		return nil, errors.Wrap(err, "cannot read 'configs/'")
	}

	for _, f := range files {
		if f.Name() == "config.yaml" {
			conf, err = ReadConfig(path + "/config.yaml")
		}
	}

	if conf == nil {
		return nil, errors.Wrap(err, "cannot find config")
	}

	if err != nil {
		return nil, errors.Wrap(err, "cannot read config")
	}

	files, err = ioutil.ReadDir(path + "/ssl")

	if err != nil {
		return nil, errors.Wrap(err, "cannot find ssl directory")
	}

	for _, f := range files {
		if strings.Contains(f.Name(), ".crt") {
			conf.Cert = path + "/ssl/" + f.Name()
		}
		if strings.Contains(f.Name(), ".key") {
			conf.Key = path + "/ssl/" + f.Name()
		}
	}

	if err := conf.validate(); err != nil {
		return nil, err
	}

	return conf, nil
}

func (c *Config) validate() error {
	if len(c.Key) <= 0 {
		return errors.New("config: missing ssl key")
	}

	if len(c.Cert) <= 0 {
		return errors.New("config: missing ssl certificate")
	}

	if len(c.Endpoint) <= 0 {
		return errors.New("config: missing api 'endpoint'")
	}

	if len(c.Host) <= 0 {
		return errors.New("config: missing 'host' value")
	}

	if c.Port <= 0 {
		return errors.New("config: invalid 'port' value")
	}

	if len(c.Gateways) <= 0 {
		return errors.New("config: missing 'payment_providers'")
	}

	return nil
}

func FindConfigs() {
	files, err := ioutil.ReadDir("./configs/ssl")

	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.Contains(f.Name(), ".crt") {

			fmt.Println(f.Name())
		}
	}
}
