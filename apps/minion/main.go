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
	bossUrl := flag.String("bossUrl", "", "boss host:port")
	advertiseUrl := flag.String("advertiseUrl", "", "self host:port to advertise externally")
	flag.Parse()

	log.Printf("Starting server. Host=%s, grpcPort=%s, metricsPort=%s, bossUrl=%s, advertiseUrl=%s",
		*host, *grpcPort, *metricsPort, *bossUrl, *advertiseUrl)

	if *bossUrl == "" {
		log.Fatal("required argument bossUrl not provided")
	}
	if *advertiseUrl == "" {
		log.Fatal("required argument advertiseUrl not provided")
	}

	minion := NewMinionServer(*host, *grpcPort, *metricsPort, *bossUrl, *advertiseUrl)
	err := minion.RegisterWithBoss()
	if err != nil {
		log.Fatal(err)
	}
	err = minion.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
