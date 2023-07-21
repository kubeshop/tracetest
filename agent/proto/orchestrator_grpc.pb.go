// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: proto/orchestrator.proto

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

const (
	Orchestrator_Connect_FullMethodName              = "/proto.Orchestrator/Connect"
	Orchestrator_RegisterTriggerAgent_FullMethodName = "/proto.Orchestrator/RegisterTriggerAgent"
	Orchestrator_SendTriggerResult_FullMethodName    = "/proto.Orchestrator/SendTriggerResult"
)

// OrchestratorClient is the client API for Orchestrator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrchestratorClient interface {
	// Connects an agent and returns the configuration that must be used by that agent
	Connect(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (*ConnectResponse, error)
	// Register an agent as a trigger agent, once connected, the server will be able to send
	// multiple trigger requests to the agent.
	RegisterTriggerAgent(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (Orchestrator_RegisterTriggerAgentClient, error)
	// Sends the trigger result back to the server
	SendTriggerResult(ctx context.Context, in *TriggerResponse, opts ...grpc.CallOption) (*Empty, error)
}

type orchestratorClient struct {
	cc grpc.ClientConnInterface
}

func NewOrchestratorClient(cc grpc.ClientConnInterface) OrchestratorClient {
	return &orchestratorClient{cc}
}

func (c *orchestratorClient) Connect(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (*ConnectResponse, error) {
	out := new(ConnectResponse)
	err := c.cc.Invoke(ctx, Orchestrator_Connect_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orchestratorClient) RegisterTriggerAgent(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (Orchestrator_RegisterTriggerAgentClient, error) {
	stream, err := c.cc.NewStream(ctx, &Orchestrator_ServiceDesc.Streams[0], Orchestrator_RegisterTriggerAgent_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &orchestratorRegisterTriggerAgentClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Orchestrator_RegisterTriggerAgentClient interface {
	Recv() (*TriggerRequest, error)
	grpc.ClientStream
}

type orchestratorRegisterTriggerAgentClient struct {
	grpc.ClientStream
}

func (x *orchestratorRegisterTriggerAgentClient) Recv() (*TriggerRequest, error) {
	m := new(TriggerRequest)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *orchestratorClient) SendTriggerResult(ctx context.Context, in *TriggerResponse, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, Orchestrator_SendTriggerResult_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrchestratorServer is the server API for Orchestrator service.
// All implementations must embed UnimplementedOrchestratorServer
// for forward compatibility
type OrchestratorServer interface {
	// Connects an agent and returns the configuration that must be used by that agent
	Connect(context.Context, *ConnectRequest) (*ConnectResponse, error)
	// Register an agent as a trigger agent, once connected, the server will be able to send
	// multiple trigger requests to the agent.
	RegisterTriggerAgent(*ConnectRequest, Orchestrator_RegisterTriggerAgentServer) error
	// Sends the trigger result back to the server
	SendTriggerResult(context.Context, *TriggerResponse) (*Empty, error)
	mustEmbedUnimplementedOrchestratorServer()
}

// UnimplementedOrchestratorServer must be embedded to have forward compatible implementations.
type UnimplementedOrchestratorServer struct {
}

func (UnimplementedOrchestratorServer) Connect(context.Context, *ConnectRequest) (*ConnectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Connect not implemented")
}
func (UnimplementedOrchestratorServer) RegisterTriggerAgent(*ConnectRequest, Orchestrator_RegisterTriggerAgentServer) error {
	return status.Errorf(codes.Unimplemented, "method RegisterTriggerAgent not implemented")
}
func (UnimplementedOrchestratorServer) SendTriggerResult(context.Context, *TriggerResponse) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendTriggerResult not implemented")
}
func (UnimplementedOrchestratorServer) mustEmbedUnimplementedOrchestratorServer() {}

// UnsafeOrchestratorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrchestratorServer will
// result in compilation errors.
type UnsafeOrchestratorServer interface {
	mustEmbedUnimplementedOrchestratorServer()
}

func RegisterOrchestratorServer(s grpc.ServiceRegistrar, srv OrchestratorServer) {
	s.RegisterService(&Orchestrator_ServiceDesc, srv)
}

func _Orchestrator_Connect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).Connect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Orchestrator_Connect_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).Connect(ctx, req.(*ConnectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orchestrator_RegisterTriggerAgent_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ConnectRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(OrchestratorServer).RegisterTriggerAgent(m, &orchestratorRegisterTriggerAgentServer{stream})
}

type Orchestrator_RegisterTriggerAgentServer interface {
	Send(*TriggerRequest) error
	grpc.ServerStream
}

type orchestratorRegisterTriggerAgentServer struct {
	grpc.ServerStream
}

func (x *orchestratorRegisterTriggerAgentServer) Send(m *TriggerRequest) error {
	return x.ServerStream.SendMsg(m)
}

func _Orchestrator_SendTriggerResult_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TriggerResponse)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).SendTriggerResult(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Orchestrator_SendTriggerResult_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).SendTriggerResult(ctx, req.(*TriggerResponse))
	}
	return interceptor(ctx, in, info, handler)
}

// Orchestrator_ServiceDesc is the grpc.ServiceDesc for Orchestrator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Orchestrator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Orchestrator",
	HandlerType: (*OrchestratorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Connect",
			Handler:    _Orchestrator_Connect_Handler,
		},
		{
			MethodName: "SendTriggerResult",
			Handler:    _Orchestrator_SendTriggerResult_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RegisterTriggerAgent",
			Handler:       _Orchestrator_RegisterTriggerAgent_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/orchestrator.proto",
}