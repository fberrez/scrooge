package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/fberrez/scrooge/transaction/api"
	"github.com/ovh/configstore"
	log "github.com/sirupsen/logrus"
)

var (
	portKey        = "port"
	defaultPort    = "3001"
	scroogeKey     = "scrooge"
	defaultScrooge = "localhost:3000"
)

func port() string {
	port, err := configstore.GetItemValue(portKey)
	if err != nil || len(port) == 0 {
		port = defaultPort
	}

	port = fmt.Sprintf(":%s", port)
	log.Infof("Server running on :%s", port)

	return port
}

func scrooge() string {
	scrooge, err := configstore.GetItemValue(scroogeKey)
	if err != nil || port == nil {
		scrooge = defaultScrooge
	}

	return scrooge
}

func main() {
	configstore.InitFromEnvironment()

	api, err := api.New()
	if err != nil {
		panic(err)
	}

	err = api.NewGrpcConnection(scrooge())
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:    port(),
		Handler: api,
	}

	// Runs the server
	go func() {
		if errListen := srv.ListenAndServe(); errListen != nil && errListen != http.ErrServerClosed {
			panic(errListen)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM)
	signal.Notify(quit, syscall.SIGINT)

	// Wait for a SIGTERM or SIGINT
	<-quit
	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}

	log.Info("Graceful shutdown")
	os.Exit(0)
}
