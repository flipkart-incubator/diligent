package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
)

const (
	defaultHost     = ""
	defaultGrpcPort = ":9191"
	defaultMetricsPort = ":2121"
)

func main() {
	// Parse command line options
	host := flag.String("host", defaultHost, "listening host")
	grpcPort := flag.String("grpc-port", defaultGrpcPort, "grpc port")
	metricsPort := flag.String("metrics-port", defaultMetricsPort, "metrics port")
	advertiseHost := flag.String("advertiseHost", defaultHost, "self host to advertise externally")
	advertisePort := flag.String("advertisePort", defaultGrpcPort, "self port to advertise externally")
	boss := flag.String("boss", "", "boss host:port")
	flag.Parse()

	log.Printf("Starting server. Host=%s, grpcPort=%s, metricsPort=%s, boss=%s, advertiseHost=%s",
		*host, *grpcPort, *metricsPort, *boss, *advertiseHost)

	if *boss == "" {
		log.Fatal("required argument boss not provided")
	}
	if *advertiseHost == "" {
		log.Fatal("required argument advertiseHost not provided")
	}

	minion := NewMinionServer(*host, *grpcPort, *metricsPort, *boss, *advertiseHost, *advertisePort)
	err := minion.RegisterWithBoss()
	if err != nil {
		log.Fatal(err)
	}
	err = minion.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
