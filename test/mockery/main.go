package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"syscall"

	"github.com/mrtroian/payservice/internal/server"
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
	path := os.Getenv("PAYSERVICE_MOCKCONFIG_PATH")
	cert := os.Getenv("SSL_CERT")
	key := os.Getenv("SSL_KEY")
	rawConfig, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatalln(err)
	}
	config := new(MockConfig)

	if err := yaml.Unmarshal(rawConfig, config); err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	handleSignals(cancel)
	servers := make([]server.Server, 0, len(config.Providers))

	for _, p := range config.Providers {
		router := mockapi.NewRouter(p.Endpoint, p.Name)
		srv := server.New()
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
