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
	UnregisterMinion(ctx context.Context, in *BossUnregisterMinonRequest, opts ...grpc.CallOption) (*BossUnregisterMinionResponse, error)
	ShowMinions(ctx context.Context, in *BossShowMinionRequest, opts ...grpc.CallOption) (*BossShowMinionResponse, error)
	// Workload Related
	RunWorkload(ctx context.Context, in *BossRunWorkloadRequest, opts ...grpc.CallOption) (*BossRunWorkloadResponse, error)
	StopWorkload(ctx context.Context, in *BossStopWorkloadRequest, opts ...grpc.CallOption) (*BossStopWorkloadResponse, error)
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

func (c *bossClient) UnregisterMinion(ctx context.Context, in *BossUnregisterMinonRequest, opts ...grpc.CallOption) (*BossUnregisterMinionResponse, error) {
	out := new(BossUnregisterMinionResponse)
	err := c.cc.Invoke(ctx, "/proto.Boss/UnregisterMinion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bossClient) ShowMinions(ctx context.Context, in *BossShowMinionRequest, opts ...grpc.CallOption) (*BossShowMinionResponse, error) {
	out := new(BossShowMinionResponse)
	err := c.cc.Invoke(ctx, "/proto.Boss/ShowMinions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bossClient) RunWorkload(ctx context.Context, in *BossRunWorkloadRequest, opts ...grpc.CallOption) (*BossRunWorkloadResponse, error) {
	out := new(BossRunWorkloadResponse)
	err := c.cc.Invoke(ctx, "/proto.Boss/RunWorkload", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bossClient) StopWorkload(ctx context.Context, in *BossStopWorkloadRequest, opts ...grpc.CallOption) (*BossStopWorkloadResponse, error) {
	out := new(BossStopWorkloadResponse)
	err := c.cc.Invoke(ctx, "/proto.Boss/StopWorkload", in, out, opts...)
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
	UnregisterMinion(context.Context, *BossUnregisterMinonRequest) (*BossUnregisterMinionResponse, error)
	ShowMinions(context.Context, *BossShowMinionRequest) (*BossShowMinionResponse, error)
	// Workload Related
	RunWorkload(context.Context, *BossRunWorkloadRequest) (*BossRunWorkloadResponse, error)
	StopWorkload(context.Context, *BossStopWorkloadRequest) (*BossStopWorkloadResponse, error)
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
func (UnimplementedBossServer) UnregisterMinion(context.Context, *BossUnregisterMinonRequest) (*BossUnregisterMinionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnregisterMinion not implemented")
}
func (UnimplementedBossServer) ShowMinions(context.Context, *BossShowMinionRequest) (*BossShowMinionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShowMinions not implemented")
}
func (UnimplementedBossServer) RunWorkload(context.Context, *BossRunWorkloadRequest) (*BossRunWorkloadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RunWorkload not implemented")
}
func (UnimplementedBossServer) StopWorkload(context.Context, *BossStopWorkloadRequest) (*BossStopWorkloadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopWorkload not implemented")
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
	in := new(BossUnregisterMinonRequest)
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
		return srv.(BossServer).UnregisterMinion(ctx, req.(*BossUnregisterMinonRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Boss_ShowMinions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BossShowMinionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BossServer).ShowMinions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Boss/ShowMinions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BossServer).ShowMinions(ctx, req.(*BossShowMinionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Boss_RunWorkload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BossRunWorkloadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BossServer).RunWorkload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Boss/RunWorkload",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BossServer).RunWorkload(ctx, req.(*BossRunWorkloadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Boss_StopWorkload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BossStopWorkloadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BossServer).StopWorkload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Boss/StopWorkload",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BossServer).StopWorkload(ctx, req.(*BossStopWorkloadRequest))
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
			MethodName: "ShowMinions",
			Handler:    _Boss_ShowMinions_Handler,
		},
		{
			MethodName: "RunWorkload",
			Handler:    _Boss_RunWorkload_Handler,
		},
		{
			MethodName: "StopWorkload",
			Handler:    _Boss_StopWorkload_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "diligent_boss.proto",
}