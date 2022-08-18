package main

import (
	"flag"
	"github.com/flipkart-incubator/diligent/pkg/buildinfo"
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

	log.Printf("Starting diligent-boss process\n")
	log.Printf("Build information:")
	log.Printf("go-version : %s\n", buildinfo.GoVersion)
	log.Printf("commit-hash: %s\n", buildinfo.CommitHash)
	log.Printf("build-time : %s\n", buildinfo.BuildTime)
	log.Printf("Arguments:")
	log.Printf("grpc-addr  : %s\n", *grpcAddr)

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
