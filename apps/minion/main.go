package main

import (
	"flag"
	"github.com/flipkart-incubator/diligent/pkg/buildinfo"
	log "github.com/sirupsen/logrus"
	"net"
)

const (
	defaultGrpcHost    = ""
	defaultGrpcPort    = "5711"
	defaultMetricsPort = "9090"
	defaultBossPort    = "5710"
	defaultLogLevel    = "warn"
)

func main() {
	// Construct default address from host and port
	defaultGrpcAddr := net.JoinHostPort(defaultGrpcHost, defaultGrpcPort)
	defaultMetricsAddr := net.JoinHostPort(defaultGrpcHost, defaultMetricsPort)

	// Parse command line options
	grpcAddr := flag.String("grpc-addr", defaultGrpcAddr, "listening host[:port] for gRPC connections")
	metricsAddr := flag.String("metrics-addr", defaultMetricsAddr, "listening host[:port] for metrics scraping")
	advertiseAddr := flag.String("advertise-addr", defaultGrpcAddr, "gRPC host[:port] to advertise externally")
	bossAddr := flag.String("boss", "", "boss host[:port]")
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

	log.Infof("Starting diligent-minion process\n")
	log.Infof("Build information:")
	log.Infof("go-version : %s\n", buildinfo.GoVersion)
	log.Infof("commit-hash: %s\n", buildinfo.CommitHash)
	log.Infof("build-time : %s\n", buildinfo.BuildTime)
	log.Infof("Arguments:")
	log.Infof("grpc-addr  : %s\n", *grpcAddr)
	log.Infof("metrics-addr  : %s\n", *metricsAddr)
	log.Infof("advertise-addr  : %s\n", *advertiseAddr)
	log.Infof("boss  : %s\n", *bossAddr)

	if *bossAddr == "" {
		log.Fatal("required argument boss not provided")
	}
	if *advertiseAddr == "" {
		log.Fatal("required argument advertise-addr not provided")
	}

	// If no port is specified for gRPC listener, use default gRPC port
	if _, _, err := net.SplitHostPort(*grpcAddr); err != nil {
		*grpcAddr = net.JoinHostPort(*grpcAddr, defaultGrpcPort)
	}

	// If no port is specified for advertise addr, use default gRPC port
	if _, _, err := net.SplitHostPort(*advertiseAddr); err != nil {
		*advertiseAddr = net.JoinHostPort(*advertiseAddr, defaultGrpcPort)
	}

	// If no port is specified for metrics listener, use default metrics port
	if _, _, err := net.SplitHostPort(*metricsAddr); err != nil {
		*metricsAddr = net.JoinHostPort(*metricsAddr, defaultMetricsPort)
	}

	// If no port is specified for boss, use default boss port
	if _, _, err := net.SplitHostPort(*bossAddr); err != nil {
		*bossAddr = net.JoinHostPort(*bossAddr, defaultBossPort)
	}

	log.Infof("Creating minion server instance")
	log.Infof("grpc-addr  : %s\n", *grpcAddr)
	log.Infof("metrics-addr  : %s\n", *metricsAddr)
	log.Infof("advertise-addr  : %s\n", *advertiseAddr)
	log.Infof("boss  : %s\n", *bossAddr)

	minion := NewMinionServer(*grpcAddr, *metricsAddr, *advertiseAddr, *bossAddr)
	err := minion.RegisterWithBoss()
	if err != nil {
		log.Fatal(err)
	}
	err = minion.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
