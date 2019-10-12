package main

import (
	"fmt"
	"net"
	"os"

	"github.com/fberrez/scrooge/balance"
	"github.com/fberrez/scrooge/scrooge"
	"github.com/ovh/configstore"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	portKey     = "port"
	defaultPort = "3000"
)

func port() string {
	port, err := configstore.GetItemValue(portKey)
	if err != nil || len(port) == 0 {
		port = defaultPort
	}

	port = fmt.Sprintf(":%s", port)
	log.Infof("Server running on %s (gRPC)", port)

	return port
}

func main() {
	configstore.InitFromEnvironment()

	lis, err := net.Listen("tcp", port())
	if err != nil {
		panic(err)
	}

	server, err := balance.New()
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	scrooge.RegisterScroogeServer(grpcServer, server)

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}

	log.Info("Graceful shutdown")
	os.Exit(0)
}
