package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"net"
)

const (
	defaultHost     = ""
	defaultGrpcPort = "5710"
)

func main() {
	// Parse command line options
	host := flag.String("host", defaultHost, "listening host")
	grpcPort := flag.String("grpc-port", defaultGrpcPort, "grpc port")
	flag.Parse()

	log.Printf("Starting server. Host=%s, grpcPort=%s", *host, *grpcPort)
	boss := NewBossServer(net.JoinHostPort(*host, *grpcPort))
	err := boss.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
