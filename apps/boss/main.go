package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"net"
)

const (
	defaultGrpcHost = ""
	defaultGrpcPort = "5710"
)

func main() {
	defaultGrpcAddr := net.JoinHostPort(defaultGrpcHost, defaultGrpcPort)

	// Parse command line options
	grpcAddr := flag.String("grpc-addr", defaultGrpcAddr, "listening host[:port] for gRPC connections")
	flag.Parse()

	log.Printf("Starting server process. Arguments: grpc-addr=%s", *grpcAddr)

	// If no port is specified for gRPC listener, use default gRPC port
	if _, _, err := net.SplitHostPort(*grpcAddr); err != nil {
		*grpcAddr = net.JoinHostPort(*grpcAddr, defaultGrpcPort)
	}

	log.Printf("Creating boss server instance. Arguments: grpc-addr=%s", *grpcAddr)
	boss := NewBossServer(*grpcAddr)
	err := boss.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
