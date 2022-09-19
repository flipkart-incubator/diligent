// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: diligent_boss.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// BossClient is the client API for Boss service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BossClient interface {
	// Ping
	Ping(ctx context.Context, in *BossPingRequest, opts ...grpc.CallOption) (*BossPingResponse, error)
	// Minions Management
	RegisterMinion(ctx context.Context, in *BossRegisterMinionRequest, opts ...grpc.CallOption) (*BossRegisterMinionResponse, error)
	UnregisterMinion(ctx context.Context, in *BossUnregisterMinionRequest, opts ...grpc.CallOption) (*BossUnregisterMinionResponse, error)
	GetMinionInfo(ctx context.Context, in *BossGetMinionInfoRequest, opts ...grpc.CallOption) (*BossGetMinionInfoResponse, error)
	// Job Control
	PrepareJob(ctx context.Context, in *BossPrepareJobRequest, opts ...grpc.CallOption) (*BossPrepareJobResponse, error)
	RunJob(ctx context.Context, in *BossRunJobRequest, opts ...grpc.CallOption) (*BossRunJobResponse, error)
	AbortJob(ctx context.Context, in *BossAbortJobRequest, opts ...grpc.CallOption) (*BossAbortJobResponse, error)
	GetJobInfo(ctx context.Context, in *BossGetJobInfoRequest, opts ...grpc.CallOption) (*BossGetJobInfoResponse, error)
	// Experiment Management
	BeginExperiment(ctx context.Context, in *BossBeginExperimentRequest, opts ...grpc.CallOption) (*BossBeginExperimentResponse, error)
	EndExperiment(ctx context.Context, in *BossEndExperimentRequest, opts ...grpc.CallOption) (*BossEndExperimentResponse, error)
	GetExperimentInfo(ctx context.Context, in *BossGetExperimentInfoRequest, opts ...grpc.CallOption) (*BossGetExperimentInfoResponse, error)
}

type bossClient struct {
	cc grpc.ClientConnInterface
}

func NewBossClient(cc grpc.ClientConnInterface) BossClient {
	return &bossClient{cc}
}

func (c *bossClient) Ping(ctx context.Context, in *BossPingRequest, opts ...grpc.CallOption) (*BossPingResponse, error) {
	out := new(BossPingResponse)
	err := c.cc.Invoke(ctx, "/proto.Boss/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bossClient) RegisterMinion(ctx context.Context, in *BossRegisterMinionRequest, opts ...grpc.CallOption) (*BossRegisterMinionResponse, error) {
	out := new(BossRegisterMinionResponse)
	err := c.cc.Invoke(ctx, "/proto.Boss/RegisterMinion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bossClient) UnregisterMinion(ctx context.Context, in *BossUnregisterMinionRequest, opts ...grpc.CallOption) (*BossUnregisterMinionResponse, error) {
	out := new(BossUnregisterMinionResponse)
	err := c.cc.Invoke(ctx, "/proto.Boss/UnregisterMinion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bossClient) GetMinionInfo(ctx context.Context, in *BossGetMinionInfoRequest, opts ...grpc.CallOption) (*BossGetMinionInfoResponse, error) {
	out := new(BossGetMinionInfoResponse)
	err := c.cc.Invoke(ctx, "/proto.Boss/GetMinionInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bossClient) PrepareJob(ctx context.Context, in *BossPrepareJobRequest, opts ...grpc.CallOption) (*BossPrepareJobResponse, error) {
	out := new(BossPrepareJobResponse)
	err := c.cc.Invoke(ctx, "/proto.Boss/PrepareJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bossClient) RunJob(ctx context.Context, in *BossRunJobRequest, opts ...grpc.CallOption) (*BossRunJobResponse, error) {
	out := new(BossRunJobResponse)
	err := c.cc.Invoke(ctx, "/proto.Boss/RunJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bossClient) AbortJob(ctx context.Context, in *BossAbortJobRequest, opts ...grpc.CallOption) (*BossAbortJobResponse, error) {
	out := new(BossAbortJobResponse)
	err := c.cc.Invoke(ctx, "/proto.Boss/AbortJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bossClient) GetJobInfo(ctx context.Context, in *BossGetJobInfoRequest, opts ...grpc.CallOption) (*BossGetJobInfoResponse, error) {
	out := new(BossGetJobInfoResponse)
	err := c.cc.Invoke(ctx, "/proto.Boss/GetJobInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bossClient) BeginExperiment(ctx context.Context, in *BossBeginExperimentRequest, opts ...grpc.CallOption) (*BossBeginExperimentResponse, error) {
	out := new(BossBeginExperimentResponse)
	err := c.cc.Invoke(ctx, "/proto.Boss/BeginExperiment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bossClient) EndExperiment(ctx context.Context, in *BossEndExperimentRequest, opts ...grpc.CallOption) (*BossEndExperimentResponse, error) {
	out := new(BossEndExperimentResponse)
	err := c.cc.Invoke(ctx, "/proto.Boss/EndExperiment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bossClient) GetExperimentInfo(ctx context.Context, in *BossGetExperimentInfoRequest, opts ...grpc.CallOption) (*BossGetExperimentInfoResponse, error) {
	out := new(BossGetExperimentInfoResponse)
	err := c.cc.Invoke(ctx, "/proto.Boss/GetExperimentInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BossServer is the server API for Boss service.
// All implementations must embed UnimplementedBossServer
// for forward compatibility
type BossServer interface {
	// Ping
	Ping(context.Context, *BossPingRequest) (*BossPingResponse, error)
	// Minions Management
	RegisterMinion(context.Context, *BossRegisterMinionRequest) (*BossRegisterMinionResponse, error)
	UnregisterMinion(context.Context, *BossUnregisterMinionRequest) (*BossUnregisterMinionResponse, error)
	GetMinionInfo(context.Context, *BossGetMinionInfoRequest) (*BossGetMinionInfoResponse, error)
	// Job Control
	PrepareJob(context.Context, *BossPrepareJobRequest) (*BossPrepareJobResponse, error)
	RunJob(context.Context, *BossRunJobRequest) (*BossRunJobResponse, error)
	AbortJob(context.Context, *BossAbortJobRequest) (*BossAbortJobResponse, error)
	GetJobInfo(context.Context, *BossGetJobInfoRequest) (*BossGetJobInfoResponse, error)
	// Experiment Management
	BeginExperiment(context.Context, *BossBeginExperimentRequest) (*BossBeginExperimentResponse, error)
	EndExperiment(context.Context, *BossEndExperimentRequest) (*BossEndExperimentResponse, error)
	GetExperimentInfo(context.Context, *BossGetExperimentInfoRequest) (*BossGetExperimentInfoResponse, error)
	mustEmbedUnimplementedBossServer()
}

// UnimplementedBossServer must be embedded to have forward compatible implementations.
type UnimplementedBossServer struct {
}

func (UnimplementedBossServer) Ping(context.Context, *BossPingRequest) (*BossPingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedBossServer) RegisterMinion(context.Context, *BossRegisterMinionRequest) (*BossRegisterMinionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterMinion not implemented")
}
func (UnimplementedBossServer) UnregisterMinion(context.Context, *BossUnregisterMinionRequest) (*BossUnregisterMinionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnregisterMinion not implemented")
}
func (UnimplementedBossServer) GetMinionInfo(context.Context, *BossGetMinionInfoRequest) (*BossGetMinionInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMinionInfo not implemented")
}
func (UnimplementedBossServer) PrepareJob(context.Context, *BossPrepareJobRequest) (*BossPrepareJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PrepareJob not implemented")
}
func (UnimplementedBossServer) RunJob(context.Context, *BossRunJobRequest) (*BossRunJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RunJob not implemented")
}
func (UnimplementedBossServer) AbortJob(context.Context, *BossAbortJobRequest) (*BossAbortJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AbortJob not implemented")
}
func (UnimplementedBossServer) GetJobInfo(context.Context, *BossGetJobInfoRequest) (*BossGetJobInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetJobInfo not implemented")
}
func (UnimplementedBossServer) BeginExperiment(context.Context, *BossBeginExperimentRequest) (*BossBeginExperimentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BeginExperiment not implemented")
}
func (UnimplementedBossServer) EndExperiment(context.Context, *BossEndExperimentRequest) (*BossEndExperimentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EndExperiment not implemented")
}
func (UnimplementedBossServer) GetExperimentInfo(context.Context, *BossGetExperimentInfoRequest) (*BossGetExperimentInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetExperimentInfo not implemented")
}
func (UnimplementedBossServer) mustEmbedUnimplementedBossServer() {}

// UnsafeBossServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BossServer will
// result in compilation errors.
type UnsafeBossServer interface {
	mustEmbedUnimplementedBossServer()
}

func RegisterBossServer(s grpc.ServiceRegistrar, srv BossServer) {
	s.RegisterService(&Boss_ServiceDesc, srv)
}

func _Boss_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BossPingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BossServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Boss/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BossServer).Ping(ctx, req.(*BossPingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Boss_RegisterMinion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BossRegisterMinionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BossServer).RegisterMinion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Boss/RegisterMinion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BossServer).RegisterMinion(ctx, req.(*BossRegisterMinionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Boss_UnregisterMinion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BossUnregisterMinionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BossServer).UnregisterMinion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Boss/UnregisterMinion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BossServer).UnregisterMinion(ctx, req.(*BossUnregisterMinionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Boss_GetMinionInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BossGetMinionInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BossServer).GetMinionInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Boss/GetMinionInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BossServer).GetMinionInfo(ctx, req.(*BossGetMinionInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Boss_PrepareJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BossPrepareJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BossServer).PrepareJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Boss/PrepareJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BossServer).PrepareJob(ctx, req.(*BossPrepareJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Boss_RunJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BossRunJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BossServer).RunJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Boss/RunJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BossServer).RunJob(ctx, req.(*BossRunJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Boss_AbortJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BossAbortJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BossServer).AbortJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Boss/AbortJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BossServer).AbortJob(ctx, req.(*BossAbortJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Boss_GetJobInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BossGetJobInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BossServer).GetJobInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Boss/GetJobInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BossServer).GetJobInfo(ctx, req.(*BossGetJobInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Boss_BeginExperiment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BossBeginExperimentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BossServer).BeginExperiment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Boss/BeginExperiment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BossServer).BeginExperiment(ctx, req.(*BossBeginExperimentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Boss_EndExperiment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BossEndExperimentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BossServer).EndExperiment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Boss/EndExperiment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BossServer).EndExperiment(ctx, req.(*BossEndExperimentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Boss_GetExperimentInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BossGetExperimentInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BossServer).GetExperimentInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Boss/GetExperimentInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BossServer).GetExperimentInfo(ctx, req.(*BossGetExperimentInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Boss_ServiceDesc is the grpc.ServiceDesc for Boss service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Boss_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Boss",
	HandlerType: (*BossServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Boss_Ping_Handler,
		},
		{
			MethodName: "RegisterMinion",
			Handler:    _Boss_RegisterMinion_Handler,
		},
		{
			MethodName: "UnregisterMinion",
			Handler:    _Boss_UnregisterMinion_Handler,
		},
		{
			MethodName: "GetMinionInfo",
			Handler:    _Boss_GetMinionInfo_Handler,
		},
		{
			MethodName: "PrepareJob",
			Handler:    _Boss_PrepareJob_Handler,
		},
		{
			MethodName: "RunJob",
			Handler:    _Boss_RunJob_Handler,
		},
		{
			MethodName: "AbortJob",
			Handler:    _Boss_AbortJob_Handler,
		},
		{
			MethodName: "GetJobInfo",
			Handler:    _Boss_GetJobInfo_Handler,
		},
		{
			MethodName: "BeginExperiment",
			Handler:    _Boss_BeginExperiment_Handler,
		},
		{
			MethodName: "EndExperiment",
			Handler:    _Boss_EndExperiment_Handler,
		},
		{
			MethodName: "GetExperimentInfo",
			Handler:    _Boss_GetExperimentInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "diligent_boss.proto",
}
