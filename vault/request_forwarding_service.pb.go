// Code generated by protoc-gen-go. DO NOT EDIT.
// source: request_forwarding_service.proto

/*
Package vault is a generated protocol buffer package.

It is generated from these files:
	request_forwarding_service.proto

It has these top-level messages:
	EchoRequest
	EchoReply
*/
package vault

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import forwarding "github.com/hashicorp/vault/helper/forwarding"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type EchoRequest struct {
	Message string `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
	// ClusterAddr is used to send up a standby node's address to the active
	// node upon heartbeat
	ClusterAddr string `protobuf:"bytes,2,opt,name=cluster_addr,json=clusterAddr" json:"cluster_addr,omitempty"`
	// ClusterAddrs is used to send up a list of cluster addresses to a dr
	// primary from a dr secondary
	ClusterAddrs []string `protobuf:"bytes,3,rep,name=cluster_addrs,json=clusterAddrs" json:"cluster_addrs,omitempty"`
}

func (m *EchoRequest) Reset()                    { *m = EchoRequest{} }
func (m *EchoRequest) String() string            { return proto.CompactTextString(m) }
func (*EchoRequest) ProtoMessage()               {}
func (*EchoRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *EchoRequest) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *EchoRequest) GetClusterAddr() string {
	if m != nil {
		return m.ClusterAddr
	}
	return ""
}

func (m *EchoRequest) GetClusterAddrs() []string {
	if m != nil {
		return m.ClusterAddrs
	}
	return nil
}

type EchoReply struct {
	Message          string   `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
	ClusterAddrs     []string `protobuf:"bytes,2,rep,name=cluster_addrs,json=clusterAddrs" json:"cluster_addrs,omitempty"`
	ReplicationState uint32   `protobuf:"varint,3,opt,name=replication_state,json=replicationState" json:"replication_state,omitempty"`
}

func (m *EchoReply) Reset()                    { *m = EchoReply{} }
func (m *EchoReply) String() string            { return proto.CompactTextString(m) }
func (*EchoReply) ProtoMessage()               {}
func (*EchoReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *EchoReply) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *EchoReply) GetClusterAddrs() []string {
	if m != nil {
		return m.ClusterAddrs
	}
	return nil
}

func (m *EchoReply) GetReplicationState() uint32 {
	if m != nil {
		return m.ReplicationState
	}
	return 0
}

func init() {
	proto.RegisterType((*EchoRequest)(nil), "vault.EchoRequest")
	proto.RegisterType((*EchoReply)(nil), "vault.EchoReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for RequestForwarding service

type RequestForwardingClient interface {
	ForwardRequest(ctx context.Context, in *forwarding.Request, opts ...grpc.CallOption) (*forwarding.Response, error)
	Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoReply, error)
}

type requestForwardingClient struct {
	cc *grpc.ClientConn
}

func NewRequestForwardingClient(cc *grpc.ClientConn) RequestForwardingClient {
	return &requestForwardingClient{cc}
}

func (c *requestForwardingClient) ForwardRequest(ctx context.Context, in *forwarding.Request, opts ...grpc.CallOption) (*forwarding.Response, error) {
	out := new(forwarding.Response)
	err := grpc.Invoke(ctx, "/vault.RequestForwarding/ForwardRequest", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *requestForwardingClient) Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoReply, error) {
	out := new(EchoReply)
	err := grpc.Invoke(ctx, "/vault.RequestForwarding/Echo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for RequestForwarding service

type RequestForwardingServer interface {
	ForwardRequest(context.Context, *forwarding.Request) (*forwarding.Response, error)
	Echo(context.Context, *EchoRequest) (*EchoReply, error)
}

func RegisterRequestForwardingServer(s *grpc.Server, srv RequestForwardingServer) {
	s.RegisterService(&_RequestForwarding_serviceDesc, srv)
}

func _RequestForwarding_ForwardRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(forwarding.Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RequestForwardingServer).ForwardRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vault.RequestForwarding/ForwardRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RequestForwardingServer).ForwardRequest(ctx, req.(*forwarding.Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _RequestForwarding_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EchoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RequestForwardingServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vault.RequestForwarding/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RequestForwardingServer).Echo(ctx, req.(*EchoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _RequestForwarding_serviceDesc = grpc.ServiceDesc{
	ServiceName: "vault.RequestForwarding",
	HandlerType: (*RequestForwardingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ForwardRequest",
			Handler:    _RequestForwarding_ForwardRequest_Handler,
		},
		{
			MethodName: "Echo",
			Handler:    _RequestForwarding_Echo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "request_forwarding_service.proto",
}

func init() { proto.RegisterFile("request_forwarding_service.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 287 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x90, 0xbf, 0x4e, 0xc3, 0x30,
	0x10, 0xc6, 0x9b, 0x96, 0x3f, 0xaa, 0xdb, 0xa2, 0xd6, 0x30, 0x44, 0x99, 0x42, 0x58, 0x22, 0x21,
	0x39, 0x12, 0x2c, 0x2c, 0x0c, 0x0c, 0xf0, 0x00, 0xe1, 0x01, 0x22, 0xd7, 0x39, 0x12, 0x4b, 0x6e,
	0x6c, 0x7c, 0x4e, 0xab, 0xac, 0x3c, 0x39, 0x6a, 0x92, 0xd2, 0x54, 0x95, 0x18, 0xef, 0x77, 0xa7,
	0xef, 0xd3, 0xfd, 0x48, 0x68, 0xe1, 0xbb, 0x06, 0x74, 0xd9, 0x97, 0xb6, 0x3b, 0x6e, 0x73, 0x59,
	0x15, 0x19, 0x82, 0xdd, 0x4a, 0x01, 0xcc, 0x58, 0xed, 0x34, 0xbd, 0xdc, 0xf2, 0x5a, 0xb9, 0xe0,
	0xa5, 0x90, 0xae, 0xac, 0xd7, 0x4c, 0xe8, 0x4d, 0x52, 0x72, 0x2c, 0xa5, 0xd0, 0xd6, 0x24, 0xed,
	0x2e, 0x29, 0x41, 0x19, 0xb0, 0xc9, 0x31, 0x22, 0x71, 0x8d, 0x01, 0xec, 0x02, 0x22, 0x4d, 0x66,
	0xef, 0xa2, 0xd4, 0x69, 0x57, 0x44, 0x7d, 0x72, 0xbd, 0x01, 0x44, 0x5e, 0x80, 0xef, 0x85, 0x5e,
	0x3c, 0x4d, 0x0f, 0x23, 0xbd, 0x27, 0x73, 0xa1, 0x6a, 0x74, 0x60, 0x33, 0x9e, 0xe7, 0xd6, 0x1f,
	0xb7, 0xeb, 0x59, 0xcf, 0xde, 0xf2, 0xdc, 0xd2, 0x07, 0xb2, 0x18, 0x9e, 0xa0, 0x3f, 0x09, 0x27,
	0xf1, 0x34, 0x9d, 0x0f, 0x6e, 0x30, 0xda, 0x91, 0x69, 0x57, 0x68, 0x54, 0xf3, 0x4f, 0xdd, 0x59,
	0xd6, 0xf8, 0x3c, 0x8b, 0x3e, 0x92, 0x95, 0x05, 0xa3, 0xa4, 0xe0, 0x4e, 0xea, 0x2a, 0x43, 0xc7,
	0x1d, 0xf8, 0x93, 0xd0, 0x8b, 0x17, 0xe9, 0x72, 0xb0, 0xf8, 0xdc, 0xf3, 0xa7, 0x1f, 0x8f, 0xac,
	0xfa, 0x37, 0x3f, 0xfe, 0x5c, 0xd0, 0x57, 0x72, 0xd3, 0x4f, 0x07, 0x05, 0xb7, 0xec, 0xa8, 0x8a,
	0xf5, 0x30, 0xb8, 0x3b, 0x85, 0x68, 0x74, 0x85, 0x10, 0x8d, 0x28, 0x23, 0x17, 0xfb, 0x6f, 0x28,
	0x65, 0xad, 0x6c, 0x36, 0x70, 0x19, 0x2c, 0x4f, 0x98, 0x51, 0x4d, 0x34, 0x5a, 0x5f, 0xb5, 0xd6,
	0x9f, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0x94, 0x1d, 0xe9, 0x21, 0xda, 0x01, 0x00, 0x00,
}
