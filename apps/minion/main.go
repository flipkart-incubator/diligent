package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
)

const (
	defaultHost     = ""
	defaultGrpcPort = ":9119"
	defaultMetricsPort = ":2112"
)

func main() {
	// Parse command line options
	host := flag.String("host", defaultHost, "listening host")
	grpcPort := flag.String("grpc-port", defaultGrpcPort, "grpc port")
	metricsPort := flag.String("metrics-port", defaultMetricsPort, "metrics port")
	flag.Parse()

	log.Printf("Starting server. Host=%s, grpcPort=%s, metricsPort=%s", *host, *grpcPort, *metricsPort)
	minion := NewMinionServer(*host, *grpcPort, *metricsPort)
	err := minion.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
