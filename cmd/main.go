package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mrtroian/payservice/internal/api"
	"github.com/mrtroian/payservice/internal/configuration"
	"github.com/mrtroian/payservice/internal/gateway"
	"github.com/mrtroian/payservice/internal/server"
)

var sigChannel = make(chan os.Signal)

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
	var config *configuration.Config
	var err error

	config, err = configuration.GetConfig()

	if err != nil {
		configPath := flag.String("config", "./configs/test.yaml",
			fmt.Sprintf("%s -config path/to/config.yaml", os.Args[0]))
		flag.Parse()

		if configPath == nil {
			log.Fatalln("config not provided")
		}

		config, err = configuration.ReadConfig(*configPath)

		if err != nil {
			log.Fatalln("cannot read config:", err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	router := api.NewRouter(config.Host, config.Endpoint)
	gm := gateway.NewGatewayManager()

	for _, p := range config.Gateways {
		gm.AddGateway(p)
	}
	srv := server.New()
	srv.SetAddr(config.Host, config.Port)
	srv.SetRouter(router)
	handleSignals(cancel)

	if err = srv.SetTLS(config.Cert, config.Key); err != nil {
		log.Fatalln(err)
	}
	srv.Start()

	<-ctx.Done()
	srv.Stop()
	os.Exit(0)
}
