package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/ubiqueworks/go-interface-usage/repo"
	"github.com/ubiqueworks/go-interface-usage/service"
	"github.com/ubiqueworks/go-interface-usage/handler/httphandler"
	"github.com/ubiqueworks/go-interface-usage/handler/rpchandler"
)

type config struct {
	BindAddr string `env:"BIND_ADDR" envDefault:"0.0.0.0" required:"true"`
	HttpPort int    `env:"HTTP_PORT" envDefault:"8000" required:"true"`
	RpcPort  int    `env:"RPC_PORT" envDefault:"9000" required:"true"`
}

func (c config) Validate() error {
	if c.BindAddr == "" {
		return fmt.Errorf("missing required HTTP_ADDR")
	}

	if c.HttpPort < 1 || c.HttpPort > 65535 {
		return fmt.Errorf("invalid HTTP_PORT: %d", c.HttpPort)
	}

	if c.RpcPort < 1 || c.RpcPort > 65535 {
		return fmt.Errorf("invalid RPC_PORT: %d", c.RpcPort)
	}

	return nil
}

func main() {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	// initialize repository
	r := repo.NewUserRepository()

	// initialize service
	s, err := service.New(r)
	if err != nil {
		log.Fatal(err)
	}

	// handle shutdown trigger
	shutdownCh := make(chan struct{}, 1)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		fmt.Println()
		close(shutdownCh)
	}()

	var wg sync.WaitGroup

	// initialize HTTP handler
	wg.Add(1)
	go func() {
		defer wg.Done()

		hh := httphandler.New(s, cfg.BindAddr, cfg.HttpPort)

		shutdownHttp, err := hh.Bootstrap()
		if err != nil {
			log.Fatal(err)
		}

		<-shutdownCh
		_ = shutdownHttp()
	}()


	// initialize RPC handler
	wg.Add(1)
	go func() {
		defer wg.Done()

		rh := rpchandler.New(s, cfg.BindAddr, cfg.RpcPort)

		shutdownRpc, err := rh.Bootstrap()
		if err != nil {
			log.Fatal(err)
		}

		<-shutdownCh
		_ = shutdownRpc()
	}()

	wg.Wait()
	log.Println("service shutdown completed...")
}
