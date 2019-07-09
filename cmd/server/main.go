package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ubiqueworks/go-solid-tutorial/handler/httphandler"
	"github.com/ubiqueworks/go-solid-tutorial/handler/rpchandler"
	"github.com/ubiqueworks/go-solid-tutorial/repo"
	"github.com/ubiqueworks/go-solid-tutorial/service"
)

func main() {
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
	hh, err := httphandler.New(s)
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		shutdownHttp, err := hh.Bootstrap()
		if err != nil {
			log.Fatal(err)
		}

		<-shutdownCh
		_ = shutdownHttp()
	}()

	// initialize RPC handler
	rh, err := rpchandler.New(s)
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

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
