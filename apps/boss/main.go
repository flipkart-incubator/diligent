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
	defaultLogLevel = "warn"
)

func main() {
	defaultGrpcAddr := net.JoinHostPort(defaultGrpcHost, defaultGrpcPort)

	// Parse command line options
	grpcAddr := flag.String("grpc-addr", defaultGrpcAddr, "listening host[:port] for gRPC connections")
	logLevel := flag.String("log-level", defaultLogLevel, "log level [debug|info|error|warn")
	flag.Parse()

	switch *logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.Errorf("bad log level setting: %s", *logLevel)
	}

	log.Printf("Starting diligent-boss process\n")
	log.Printf("Build information:")
	log.Printf("go-version : %s\n", buildinfo.GoVersion)
	log.Printf("commit-hash: %s\n", buildinfo.CommitHash)
	log.Printf("build-time : %s\n", buildinfo.BuildTime)
	log.Printf("Arguments:")
	log.Printf("grpc-addr  : %s\n", *grpcAddr)
	log.Printf("log-level  : %s\n", *logLevel)

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
