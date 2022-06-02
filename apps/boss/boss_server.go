package main

import (
	"context"
	"github.com/flipkart-incubator/diligent/pkg/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"sync"
)

type BossServer struct {
	proto.UnimplementedBossServer
	host        string
	grpcPort    string

	mut        *sync.Mutex
	minionUrls map[string]bool
}

func NewBossServer(host, grpcPort string) *BossServer {
	return &BossServer{
		host:        host,
		grpcPort:    grpcPort,
		mut:         &sync.Mutex{},
		minionUrls:  make(map[string]bool),
	}
}

func (s *BossServer) Serve() error {
	// Create listening port
	address := s.host + s.grpcPort
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register our this server object with the gRPC server
	proto.RegisterBossServer(grpcServer, s)

	// Start serving
	err = grpcServer.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}

func (s *BossServer) Ping(context.Context, *proto.BossPingRequest) (*proto.BossPingResponse, error) {
	log.Infof("GRPC: Ping")
	s.mut.Lock()
	defer s.mut.Unlock()

	return &proto.BossPingResponse{}, nil
}

func (s *BossServer) MinionRegister(_ context.Context, in *proto.MinionRegisterRequest) (*proto.MinionRegisterResponse, error) {
	url := in.GetUrl()
	log.Infof("GRPC: MinionRegister(%s)", url)
	s.mut.Lock()
	defer s.mut.Unlock()
	s.minionUrls[url] = true
	return &proto.MinionRegisterResponse{
		IsOk:          true,
		FailureReason: "",
	}, nil
}

func (s *BossServer) MinionUnregister(_ context.Context, in *proto.MinionUnregisterRequest) (*proto.MinionUnregisterResponse, error) {
	url := in.GetUrl()
	log.Infof("GRPC: MinionUnregister(%s)", url)
	s.mut.Lock()
	defer s.mut.Unlock()
	delete(s.minionUrls, url)
	return &proto.MinionUnregisterResponse{
		IsOk:          true,
		FailureReason: "",
	}, nil
}

func (s *BossServer) MinionShow(context.Context, *proto.MinionShowRequest) (*proto.MinionShowResponse, error) {
	log.Infof("GRPC: MinionShow()")
	s.mut.Lock()
	defer s.mut.Unlock()
	minions := make([]*proto.MinionStatus, 0)
	for url := range s.minionUrls {
		minions = append(minions, &proto.MinionStatus{
			Url:         url,
			IsReachable: true, // TODO Make an actual call and check this
		})
	}
	return &proto.MinionShowResponse{
		Minions: minions,
	}, nil
}
