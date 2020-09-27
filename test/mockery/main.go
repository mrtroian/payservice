package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"

	"syscall"

	"github.com/mrtroian/payservice/internal/api"
	"github.com/mrtroian/payservice/test/mockery/mockapi"

	"gopkg.in/yaml.v2"
)

var sigChannel = make(chan os.Signal)

type MockConfig struct {
	Providers []ProviderConf `yaml:"payment_providers"`
}

type ProviderConf struct {
	Port     int    `yaml:"port"`
	Host     string `yaml:"host"`
	Endpoint string `yaml:"endpoint"`
	Name     string `yaml:"name"`
}

func handleSignals(cancel context.CancelFunc) {
	go func() {
		signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

		for {
			sig := <-sigChannel
			switch sig {
			case os.Interrupt:
				fmt.Print("\r")
				log.Println("got", os.Interrupt)
				cancel()
				return
			}
		}
	}()
}

func main() {
	var (
		cert     string
		key      string
		confPath string
	)
	path := os.Getenv("PAYSERVICE_CONFIGS_DIR")

	if len(path) <= 0 {
		log.Fatalln("'PAYSERVICE_CONFIGS_DIR' not set in env")
	}
	files, err := ioutil.ReadDir(path)

	if err != nil {
		log.Fatalln("cannot find 'configs/':", err)
	}

	for _, f := range files {
		if f.Name() == "mockconfig.yaml" {
			confPath = path + f.Name()
		}
	}

	files, err = ioutil.ReadDir(path + "ssl/")

	if err != nil {
		log.Fatalln("cannot find 'configs/':", err)
	}

	for _, f := range files {
		if strings.Contains(f.Name(), ".crt") {
			cert = path + "ssl/" + f.Name()
		}
		if strings.Contains(f.Name(), ".key") {
			key = path + "ssl/" + f.Name()
		}
	}

	rawConfig, err := ioutil.ReadFile(confPath)

	if err != nil {
		log.Fatalln("cannot read from '/configs'", err)
	}

	config := new(MockConfig)

	if err := yaml.Unmarshal(rawConfig, config); err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	handleSignals(cancel)
	servers := make([]api.Server, 0, len(config.Providers))

	for _, p := range config.Providers {
		router := mockapi.NewRouter(p.Endpoint, p.Name)
		srv := api.NewServer()
		srv.SetAddr(p.Host, p.Port)
		srv.SetRouter(router)

		if err = srv.SetTLS(cert, key); err != nil {
			log.Fatalln(err)
		}

		srv.Start()
		servers = append(servers, srv)
	}

	<-ctx.Done()
	for _, srv := range servers {
		srv.Stop()
	}
	os.Exit(0)
}
