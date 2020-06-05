// Code generated by protoc-gen-go. DO NOT EDIT.
// source: grpc/lb/v1/load_balancer.proto

package grpc_lb_v1

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	duration "github.com/golang/protobuf/ptypes/duration"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type LoadBalanceRequest struct {
	// Types that are valid to be assigned to LoadBalanceRequestType:
	//	*LoadBalanceRequest_InitialRequest
	//	*LoadBalanceRequest_ClientStats
	LoadBalanceRequestType isLoadBalanceRequest_LoadBalanceRequestType `protobuf_oneof:"load_balance_request_type"`
	XXX_NoUnkeyedLiteral   struct{}                                    `json:"-"`
	XXX_unrecognized       []byte                                      `json:"-"`
	XXX_sizecache          int32                                       `json:"-"`
}

func (m *LoadBalanceRequest) Reset()         { *m = LoadBalanceRequest{} }
func (m *LoadBalanceRequest) String() string { return proto.CompactTextString(m) }
func (*LoadBalanceRequest) ProtoMessage()    {}
func (*LoadBalanceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7cd3f6d792743fdf, []int{0}
}

func (m *LoadBalanceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoadBalanceRequest.Unmarshal(m, b)
}
func (m *LoadBalanceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoadBalanceRequest.Marshal(b, m, deterministic)
}
func (m *LoadBalanceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoadBalanceRequest.Merge(m, src)
}
func (m *LoadBalanceRequest) XXX_Size() int {
	return xxx_messageInfo_LoadBalanceRequest.Size(m)
}
func (m *LoadBalanceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LoadBalanceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LoadBalanceRequest proto.InternalMessageInfo

type isLoadBalanceRequest_LoadBalanceRequestType interface {
	isLoadBalanceRequest_LoadBalanceRequestType()
}

type LoadBalanceRequest_InitialRequest struct {
	InitialRequest *InitialLoadBalanceRequest `protobuf:"bytes,1,opt,name=initial_request,json=initialRequest,proto3,oneof"`
}

type LoadBalanceRequest_ClientStats struct {
	ClientStats *ClientStats `protobuf:"bytes,2,opt,name=client_stats,json=clientStats,proto3,oneof"`
}

func (*LoadBalanceRequest_InitialRequest) isLoadBalanceRequest_LoadBalanceRequestType() {}

func (*LoadBalanceRequest_ClientStats) isLoadBalanceRequest_LoadBalanceRequestType() {}

func (m *LoadBalanceRequest) GetLoadBalanceRequestType() isLoadBalanceRequest_LoadBalanceRequestType {
	if m != nil {
		return m.LoadBalanceRequestType
	}
	return nil
}

func (m *LoadBalanceRequest) GetInitialRequest() *InitialLoadBalanceRequest {
	if x, ok := m.GetLoadBalanceRequestType().(*LoadBalanceRequest_InitialRequest); ok {
		return x.InitialRequest
	}
	return nil
}

func (m *LoadBalanceRequest) GetClientStats() *ClientStats {
	if x, ok := m.GetLoadBalanceRequestType().(*LoadBalanceRequest_ClientStats); ok {
		return x.ClientStats
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*LoadBalanceRequest) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*LoadBalanceRequest_InitialRequest)(nil),
		(*LoadBalanceRequest_ClientStats)(nil),
	}
}

type InitialLoadBalanceRequest struct {
	// The name of the load balanced service (e.g., service.googleapis.com). Its
	// length should be less than 256 bytes.
	// The name might include a port number. How to handle the port number is up
	// to the balancer.
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InitialLoadBalanceRequest) Reset()         { *m = InitialLoadBalanceRequest{} }
func (m *InitialLoadBalanceRequest) String() string { return proto.CompactTextString(m) }
func (*InitialLoadBalanceRequest) ProtoMessage()    {}
func (*InitialLoadBalanceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7cd3f6d792743fdf, []int{1}
}

func (m *InitialLoadBalanceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InitialLoadBalanceRequest.Unmarshal(m, b)
}
func (m *InitialLoadBalanceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InitialLoadBalanceRequest.Marshal(b, m, deterministic)
}
func (m *InitialLoadBalanceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InitialLoadBalanceRequest.Merge(m, src)
}
func (m *InitialLoadBalanceRequest) XXX_Size() int {
	return xxx_messageInfo_InitialLoadBalanceRequest.Size(m)
}
func (m *InitialLoadBalanceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_InitialLoadBalanceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_InitialLoadBalanceRequest proto.InternalMessageInfo

func (m *InitialLoadBalanceRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// Contains the number of calls finished for a particular load balance token.
type ClientStatsPerToken struct {
	// See Server.load_balance_token.
	LoadBalanceToken string `protobuf:"bytes,1,opt,name=load_balance_token,json=loadBalanceToken,proto3" json:"load_balance_token,omitempty"`
	// The total number of RPCs that finished associated with the token.
	NumCalls             int64    `protobuf:"varint,2,opt,name=num_calls,json=numCalls,proto3" json:"num_calls,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ClientStatsPerToken) Reset()         { *m = ClientStatsPerToken{} }
func (m *ClientStatsPerToken) String() string { return proto.CompactTextString(m) }
func (*ClientStatsPerToken) ProtoMessage()    {}
func (*ClientStatsPerToken) Descriptor() ([]byte, []int) {
	return fileDescriptor_7cd3f6d792743fdf, []int{2}
}

func (m *ClientStatsPerToken) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClientStatsPerToken.Unmarshal(m, b)
}
func (m *ClientStatsPerToken) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClientStatsPerToken.Marshal(b, m, deterministic)
}
func (m *ClientStatsPerToken) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClientStatsPerToken.Merge(m, src)
}
func (m *ClientStatsPerToken) XXX_Size() int {
	return xxx_messageInfo_ClientStatsPerToken.Size(m)
}
func (m *ClientStatsPerToken) XXX_DiscardUnknown() {
	xxx_messageInfo_ClientStatsPerToken.DiscardUnknown(m)
}

var xxx_messageInfo_ClientStatsPerToken proto.InternalMessageInfo

func (m *ClientStatsPerToken) GetLoadBalanceToken() string {
	if m != nil {
		return m.LoadBalanceToken
	}
	return ""
}

func (m *ClientStatsPerToken) GetNumCalls() int64 {
	if m != nil {
		return m.NumCalls
	}
	return 0
}

// Contains client level statistics that are useful to load balancing. Each
// count except the timestamp should be reset to zero after reporting the stats.
type ClientStats struct {
	// The timestamp of generating the report.
	Timestamp *timestamp.Timestamp `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// The total number of RPCs that started.
	NumCallsStarted int64 `protobuf:"varint,2,opt,name=num_calls_started,json=numCallsStarted,proto3" json:"num_calls_started,omitempty"`
	// The total number of RPCs that finished.
	NumCallsFinished int64 `protobuf:"varint,3,opt,name=num_calls_finished,json=numCallsFinished,proto3" json:"num_calls_finished,omitempty"`
	// The total number of RPCs that failed to reach a server except dropped RPCs.
	NumCallsFinishedWithClientFailedToSend int64 `protobuf:"varint,6,opt,name=num_calls_finished_with_client_failed_to_send,json=numCallsFinishedWithClientFailedToSend,proto3" json:"num_calls_finished_with_client_failed_to_send,omitempty"`
	// The total number of RPCs that finished and are known to have been received
	// by a server.
	NumCallsFinishedKnownReceived int64 `protobuf:"varint,7,opt,name=num_calls_finished_known_received,json=numCallsFinishedKnownReceived,proto3" json:"num_calls_finished_known_received,omitempty"`
	// The list of dropped calls.
	CallsFinishedWithDrop []*ClientStatsPerToken `protobuf:"bytes,8,rep,name=calls_finished_with_drop,json=callsFinishedWithDrop,proto3" json:"calls_finished_with_drop,omitempty"`
	XXX_NoUnkeyedLiteral  struct{}               `json:"-"`
	XXX_unrecognized      []byte                 `json:"-"`
	XXX_sizecache         int32                  `json:"-"`
}

func (m *ClientStats) Reset()         { *m = ClientStats{} }
func (m *ClientStats) String() string { return proto.CompactTextString(m) }
func (*ClientStats) ProtoMessage()    {}
func (*ClientStats) Descriptor() ([]byte, []int) {
	return fileDescriptor_7cd3f6d792743fdf, []int{3}
}

func (m *ClientStats) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClientStats.Unmarshal(m, b)
}
func (m *ClientStats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClientStats.Marshal(b, m, deterministic)
}
func (m *ClientStats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClientStats.Merge(m, src)
}
func (m *ClientStats) XXX_Size() int {
	return xxx_messageInfo_ClientStats.Size(m)
}
func (m *ClientStats) XXX_DiscardUnknown() {
	xxx_messageInfo_ClientStats.DiscardUnknown(m)
}

var xxx_messageInfo_ClientStats proto.InternalMessageInfo

func (m *ClientStats) GetTimestamp() *timestamp.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *ClientStats) GetNumCallsStarted() int64 {
	if m != nil {
		return m.NumCallsStarted
	}
	return 0
}

func (m *ClientStats) GetNumCallsFinished() int64 {
	if m != nil {
		return m.NumCallsFinished
	}
	return 0
}

func (m *ClientStats) GetNumCallsFinishedWithClientFailedToSend() int64 {
	if m != nil {
		return m.NumCallsFinishedWithClientFailedToSend
	}
	return 0
}

func (m *ClientStats) GetNumCallsFinishedKnownReceived() int64 {
	if m != nil {
		return m.NumCallsFinishedKnownReceived
	}
	return 0
}

func (m *ClientStats) GetCallsFinishedWithDrop() []*ClientStatsPerToken {
	if m != nil {
		return m.CallsFinishedWithDrop
	}
	return nil
}

type LoadBalanceResponse struct {
	// Types that are valid to be assigned to LoadBalanceResponseType:
	//	*LoadBalanceResponse_InitialResponse
	//	*LoadBalanceResponse_ServerList
	//	*LoadBalanceResponse_FallbackResponse
	LoadBalanceResponseType isLoadBalanceResponse_LoadBalanceResponseType `protobuf_oneof:"load_balance_response_type"`
	XXX_NoUnkeyedLiteral    struct{}                                      `json:"-"`
	XXX_unrecognized        []byte                                        `json:"-"`
	XXX_sizecache           int32                                         `json:"-"`
}

func (m *LoadBalanceResponse) Reset()         { *m = LoadBalanceResponse{} }
func (m *LoadBalanceResponse) String() string { return proto.CompactTextString(m) }
func (*LoadBalanceResponse) ProtoMessage()    {}
func (*LoadBalanceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7cd3f6d792743fdf, []int{4}
}

func (m *LoadBalanceResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoadBalanceResponse.Unmarshal(m, b)
}
func (m *LoadBalanceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoadBalanceResponse.Marshal(b, m, deterministic)
}
func (m *LoadBalanceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoadBalanceResponse.Merge(m, src)
}
func (m *LoadBalanceResponse) XXX_Size() int {
	return xxx_messageInfo_LoadBalanceResponse.Size(m)
}
func (m *LoadBalanceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_LoadBalanceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_LoadBalanceResponse proto.InternalMessageInfo

type isLoadBalanceResponse_LoadBalanceResponseType interface {
	isLoadBalanceResponse_LoadBalanceResponseType()
}

type LoadBalanceResponse_InitialResponse struct {
	InitialResponse *InitialLoadBalanceResponse `protobuf:"bytes,1,opt,name=initial_response,json=initialResponse,proto3,oneof"`
}

type LoadBalanceResponse_ServerList struct {
	ServerList *ServerList `protobuf:"bytes,2,opt,name=server_list,json=serverList,proto3,oneof"`
}

type LoadBalanceResponse_FallbackResponse struct {
	FallbackResponse *FallbackResponse `protobuf:"bytes,3,opt,name=fallback_response,json=fallbackResponse,proto3,oneof"`
}

func (*LoadBalanceResponse_InitialResponse) isLoadBalanceResponse_LoadBalanceResponseType() {}

func (*LoadBalanceResponse_ServerList) isLoadBalanceResponse_LoadBalanceResponseType() {}

func (*LoadBalanceResponse_FallbackResponse) isLoadBalanceResponse_LoadBalanceResponseType() {}

func (m *LoadBalanceResponse) GetLoadBalanceResponseType() isLoadBalanceResponse_LoadBalanceResponseType {
	if m != nil {
		return m.LoadBalanceResponseType
	}
	return nil
}

func (m *LoadBalanceResponse) GetInitialResponse() *InitialLoadBalanceResponse {
	if x, ok := m.GetLoadBalanceResponseType().(*LoadBalanceResponse_InitialResponse); ok {
		return x.InitialResponse
	}
	return nil
}

func (m *LoadBalanceResponse) GetServerList() *ServerList {
	if x, ok := m.GetLoadBalanceResponseType().(*LoadBalanceResponse_ServerList); ok {
		return x.ServerList
	}
	return nil
}

func (m *LoadBalanceResponse) GetFallbackResponse() *FallbackResponse {
	if x, ok := m.GetLoadBalanceResponseType().(*LoadBalanceResponse_FallbackResponse); ok {
		return x.FallbackResponse
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*LoadBalanceResponse) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*LoadBalanceResponse_InitialResponse)(nil),
		(*LoadBalanceResponse_ServerList)(nil),
		(*LoadBalanceResponse_FallbackResponse)(nil),
	}
}

type FallbackResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FallbackResponse) Reset()         { *m = FallbackResponse{} }
func (m *FallbackResponse) String() string { return proto.CompactTextString(m) }
func (*FallbackResponse) ProtoMessage()    {}
func (*FallbackResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7cd3f6d792743fdf, []int{5}
}

func (m *FallbackResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FallbackResponse.Unmarshal(m, b)
}
func (m *FallbackResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FallbackResponse.Marshal(b, m, deterministic)
}
func (m *FallbackResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FallbackResponse.Merge(m, src)
}
func (m *FallbackResponse) XXX_Size() int {
	return xxx_messageInfo_FallbackResponse.Size(m)
}
func (m *FallbackResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_FallbackResponse.DiscardUnknown(m)
}

var xxx_messageInfo_FallbackResponse proto.InternalMessageInfo

type InitialLoadBalanceResponse struct {
	// This interval defines how often the client should send the client stats
	// to the load balancer. Stats should only be reported when the duration is
	// positive.
	ClientStatsReportInterval *duration.Duration `protobuf:"bytes,2,opt,name=client_stats_report_interval,json=clientStatsReportInterval,proto3" json:"client_stats_report_interval,omitempty"`
	XXX_NoUnkeyedLiteral      struct{}           `json:"-"`
	XXX_unrecognized          []byte             `json:"-"`
	XXX_sizecache             int32              `json:"-"`
}

func (m *InitialLoadBalanceResponse) Reset()         { *m = InitialLoadBalanceResponse{} }
func (m *InitialLoadBalanceResponse) String() string { return proto.CompactTextString(m) }
func (*InitialLoadBalanceResponse) ProtoMessage()    {}
func (*InitialLoadBalanceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7cd3f6d792743fdf, []int{6}
}

func (m *InitialLoadBalanceResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InitialLoadBalanceResponse.Unmarshal(m, b)
}
func (m *InitialLoadBalanceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InitialLoadBalanceResponse.Marshal(b, m, deterministic)
}
func (m *InitialLoadBalanceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InitialLoadBalanceResponse.Merge(m, src)
}
func (m *InitialLoadBalanceResponse) XXX_Size() int {
	return xxx_messageInfo_InitialLoadBalanceResponse.Size(m)
}
func (m *InitialLoadBalanceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_InitialLoadBalanceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_InitialLoadBalanceResponse proto.InternalMessageInfo

func (m *InitialLoadBalanceResponse) GetClientStatsReportInterval() *duration.Duration {
	if m != nil {
		return m.ClientStatsReportInterval
	}
	return nil
}

type ServerList struct {
	// Contains a list of servers selected by the load balancer. The list will
	// be updated when server resolutions change or as needed to balance load
	// across more servers. The client should consume the server list in order
	// unless instructed otherwise via the client_config.
	Servers              []*Server `protobuf:"bytes,1,rep,name=servers,proto3" json:"servers,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *ServerList) Reset()         { *m = ServerList{} }
func (m *ServerList) String() string { return proto.CompactTextString(m) }
func (*ServerList) ProtoMessage()    {}
func (*ServerList) Descriptor() ([]byte, []int) {
	return fileDescriptor_7cd3f6d792743fdf, []int{7}
}

func (m *ServerList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServerList.Unmarshal(m, b)
}
func (m *ServerList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServerList.Marshal(b, m, deterministic)
}
func (m *ServerList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServerList.Merge(m, src)
}
func (m *ServerList) XXX_Size() int {
	return xxx_messageInfo_ServerList.Size(m)
}
func (m *ServerList) XXX_DiscardUnknown() {
	xxx_messageInfo_ServerList.DiscardUnknown(m)
}

var xxx_messageInfo_ServerList proto.InternalMessageInfo

func (m *ServerList) GetServers() []*Server {
	if m != nil {
		return m.Servers
	}
	return nil
}

// Contains server information. When the drop field is not true, use the other
// fields.
type Server struct {
	// A resolved address for the server, serialized in network-byte-order. It may
	// either be an IPv4 or IPv6 address.
	IpAddress []byte `protobuf:"bytes,1,opt,name=ip_address,json=ipAddress,proto3" json:"ip_address,omitempty"`
	// A resolved port number for the server.
	Port int32 `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	// An opaque but printable token for load reporting. The client must include
	// the token of the picked server into the initial metadata when it starts a
	// call to that server. The token is used by the server to verify the request
	// and to allow the server to report load to the gRPC LB system. The token is
	// also used in client stats for reporting dropped calls.
	//
	// Its length can be variable but must be less than 50 bytes.
	LoadBalanceToken string `protobuf:"bytes,3,opt,name=load_balance_token,json=loadBalanceToken,proto3" json:"load_balance_token,omitempty"`
	// Indicates whether this particular request should be dropped by the client.
	// If the request is dropped, there will be a corresponding entry in
	// ClientStats.calls_finished_with_drop.
	Drop                 bool     `protobuf:"varint,4,opt,name=drop,proto3" json:"drop,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Server) Reset()         { *m = Server{} }
func (m *Server) String() string { return proto.CompactTextString(m) }
func (*Server) ProtoMessage()    {}
func (*Server) Descriptor() ([]byte, []int) {
	return fileDescriptor_7cd3f6d792743fdf, []int{8}
}

func (m *Server) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Server.Unmarshal(m, b)
}
func (m *Server) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Server.Marshal(b, m, deterministic)
}
func (m *Server) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Server.Merge(m, src)
}
func (m *Server) XXX_Size() int {
	return xxx_messageInfo_Server.Size(m)
}
func (m *Server) XXX_DiscardUnknown() {
	xxx_messageInfo_Server.DiscardUnknown(m)
}

var xxx_messageInfo_Server proto.InternalMessageInfo

func (m *Server) GetIpAddress() []byte {
	if m != nil {
		return m.IpAddress
	}
	return nil
}

func (m *Server) GetPort() int32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *Server) GetLoadBalanceToken() string {
	if m != nil {
		return m.LoadBalanceToken
	}
	return ""
}

func (m *Server) GetDrop() bool {
	if m != nil {
		return m.Drop
	}
	return false
}

func init() {
	proto.RegisterType((*LoadBalanceRequest)(nil), "grpc.lb.v1.LoadBalanceRequest")
	proto.RegisterType((*InitialLoadBalanceRequest)(nil), "grpc.lb.v1.InitialLoadBalanceRequest")
	proto.RegisterType((*ClientStatsPerToken)(nil), "grpc.lb.v1.ClientStatsPerToken")
	proto.RegisterType((*ClientStats)(nil), "grpc.lb.v1.ClientStats")
	proto.RegisterType((*LoadBalanceResponse)(nil), "grpc.lb.v1.LoadBalanceResponse")
	proto.RegisterType((*FallbackResponse)(nil), "grpc.lb.v1.FallbackResponse")
	proto.RegisterType((*InitialLoadBalanceResponse)(nil), "grpc.lb.v1.InitialLoadBalanceResponse")
	proto.RegisterType((*ServerList)(nil), "grpc.lb.v1.ServerList")
	proto.RegisterType((*Server)(nil), "grpc.lb.v1.Server")
}

func init() { proto.RegisterFile("grpc/lb/v1/load_balancer.proto", fileDescriptor_7cd3f6d792743fdf) }

var fileDescriptor_7cd3f6d792743fdf = []byte{
	// 769 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x55, 0xdd, 0x6e, 0xdb, 0x36,
	0x14, 0x8e, 0x62, 0x25, 0x75, 0x8e, 0xb3, 0x45, 0x61, 0xb1, 0xcd, 0x71, 0xd3, 0x24, 0x13, 0xb0,
	0x22, 0x18, 0x3a, 0x79, 0xc9, 0x6e, 0x36, 0x60, 0x17, 0x9b, 0x5b, 0x04, 0x69, 0xda, 0x8b, 0x80,
	0x0e, 0xd0, 0xa1, 0xc0, 0xc0, 0x51, 0x12, 0xed, 0x10, 0xa1, 0x49, 0x8d, 0xa2, 0x5d, 0xec, 0x66,
	0x37, 0x7b, 0x81, 0x3d, 0xca, 0x5e, 0x61, 0x6f, 0x36, 0x88, 0xa4, 0x2c, 0xd5, 0xae, 0xb1, 0x2b,
	0x91, 0xe7, 0x7c, 0xfc, 0xce, 0xff, 0x11, 0x9c, 0x4c, 0x75, 0x91, 0x0d, 0x45, 0x3a, 0x5c, 0x5c,
	0x0c, 0x85, 0xa2, 0x39, 0x49, 0xa9, 0xa0, 0x32, 0x63, 0x3a, 0x29, 0xb4, 0x32, 0x0a, 0x41, 0xa5,
	0x4f, 0x44, 0x9a, 0x2c, 0x2e, 0x06, 0x27, 0x53, 0xa5, 0xa6, 0x82, 0x0d, 0xad, 0x26, 0x9d, 0x4f,
	0x86, 0xf9, 0x5c, 0x53, 0xc3, 0x95, 0x74, 0xd8, 0xc1, 0xe9, 0xaa, 0xde, 0xf0, 0x19, 0x2b, 0x0d,
	0x9d, 0x15, 0x0e, 0x10, 0xff, 0x1b, 0x00, 0x7a, 0xa3, 0x68, 0x3e, 0x72, 0x36, 0x30, 0xfb, 0x7d,
	0xce, 0x4a, 0x83, 0x6e, 0xe1, 0x80, 0x4b, 0x6e, 0x38, 0x15, 0x44, 0x3b, 0x51, 0x3f, 0x38, 0x0b,
	0xce, 0x7b, 0x97, 0x5f, 0x25, 0x8d, 0xf5, 0xe4, 0x95, 0x83, 0xac, 0xbf, 0xbf, 0xde, 0xc2, 0x9f,
	0xfa, 0xf7, 0x35, 0xe3, 0x8f, 0xb0, 0x9f, 0x09, 0xce, 0xa4, 0x21, 0xa5, 0xa1, 0xa6, 0xec, 0x6f,
	0x5b, 0xba, 0x2f, 0xda, 0x74, 0x2f, 0xac, 0x7e, 0x5c, 0xa9, 0xaf, 0xb7, 0x70, 0x2f, 0x6b, 0xae,
	0xa3, 0x27, 0x70, 0xd4, 0x4e, 0x45, 0xed, 0x14, 0x31, 0x7f, 0x14, 0x2c, 0x1e, 0xc2, 0xd1, 0x46,
	0x4f, 0x10, 0x82, 0x50, 0xd2, 0x19, 0xb3, 0xee, 0xef, 0x61, 0x7b, 0x8e, 0x7f, 0x83, 0xc7, 0x2d,
	0x5b, 0xb7, 0x4c, 0xdf, 0xa9, 0x07, 0x26, 0xd1, 0x73, 0x40, 0x1f, 0x18, 0x31, 0x95, 0xd4, 0x3f,
	0x8c, 0x44, 0x43, 0xed, 0xd0, 0x4f, 0x60, 0x4f, 0xce, 0x67, 0x24, 0xa3, 0x42, 0xb8, 0x68, 0x3a,
	0xb8, 0x2b, 0xe7, 0xb3, 0x17, 0xd5, 0x3d, 0xfe, 0xa7, 0x03, 0xbd, 0x96, 0x09, 0xf4, 0x3d, 0xec,
	0x2d, 0x33, 0xef, 0x33, 0x39, 0x48, 0x5c, 0x6d, 0x92, 0xba, 0x36, 0xc9, 0x5d, 0x8d, 0xc0, 0x0d,
	0x18, 0x7d, 0x0d, 0x87, 0x4b, 0x33, 0x55, 0xea, 0xb4, 0x61, 0xb9, 0x37, 0x77, 0x50, 0x9b, 0x1b,
	0x3b, 0x71, 0x15, 0x40, 0x83, 0x9d, 0x70, 0xc9, 0xcb, 0x7b, 0x96, 0xf7, 0x3b, 0x16, 0x1c, 0xd5,
	0xe0, 0x2b, 0x2f, 0x47, 0xbf, 0xc2, 0x37, 0xeb, 0x68, 0xf2, 0x9e, 0x9b, 0x7b, 0xe2, 0x2b, 0x35,
	0xa1, 0x5c, 0xb0, 0x9c, 0x18, 0x45, 0x4a, 0x26, 0xf3, 0xfe, 0xae, 0x25, 0x7a, 0xb6, 0x4a, 0xf4,
	0x96, 0x9b, 0x7b, 0x17, 0xeb, 0x95, 0xc5, 0xdf, 0xa9, 0x31, 0x93, 0x39, 0xba, 0x86, 0x2f, 0x3f,
	0x42, 0xff, 0x20, 0xd5, 0x7b, 0x49, 0x34, 0xcb, 0x18, 0x5f, 0xb0, 0xbc, 0xff, 0xc8, 0x52, 0x3e,
	0x5d, 0xa5, 0x7c, 0x5d, 0xa1, 0xb0, 0x07, 0xa1, 0x5f, 0xa0, 0xff, 0x31, 0x27, 0x73, 0xad, 0x8a,
	0x7e, 0xf7, 0xac, 0x73, 0xde, 0xbb, 0x3c, 0xdd, 0xd0, 0x46, 0x75, 0x69, 0xf1, 0x67, 0xd9, 0xaa,
	0xc7, 0x2f, 0xb5, 0x2a, 0x6e, 0xc2, 0x6e, 0x18, 0xed, 0xdc, 0x84, 0xdd, 0x9d, 0x68, 0x37, 0xfe,
	0x7b, 0x1b, 0x1e, 0x7f, 0xd0, 0x3f, 0x65, 0xa1, 0x64, 0xc9, 0xd0, 0x18, 0xa2, 0x66, 0x14, 0x9c,
	0xcc, 0x57, 0xf0, 0xd9, 0xff, 0xcd, 0x82, 0x43, 0x5f, 0x6f, 0xe1, 0x83, 0xe5, 0x30, 0x78, 0xd2,
	0x1f, 0xa0, 0x57, 0x32, 0xbd, 0x60, 0x9a, 0x08, 0x5e, 0x1a, 0x3f, 0x0c, 0x9f, 0xb7, 0xf9, 0xc6,
	0x56, 0xfd, 0x86, 0xdb, 0x61, 0x82, 0x72, 0x79, 0x43, 0xaf, 0xe1, 0x70, 0x42, 0x85, 0x48, 0x69,
	0xf6, 0xd0, 0x38, 0xd4, 0xb1, 0x04, 0xc7, 0x6d, 0x82, 0x2b, 0x0f, 0x6a, 0xb9, 0x11, 0x4d, 0x56,
	0x64, 0xa3, 0x63, 0x18, 0xac, 0xcc, 0x95, 0x53, 0xb8, 0xc1, 0x42, 0x10, 0xad, 0xb2, 0xc4, 0x7f,
	0xc2, 0x60, 0x73, 0xa8, 0xe8, 0x1d, 0x1c, 0xb7, 0xa7, 0x9c, 0x68, 0x56, 0x28, 0x6d, 0x08, 0x97,
	0x86, 0xe9, 0x05, 0x15, 0x3e, 0xd0, 0xa3, 0xb5, 0xd6, 0x7f, 0xe9, 0xd7, 0x16, 0x3e, 0x6a, 0x4d,
	0x3d, 0xb6, 0x8f, 0x5f, 0xf9, 0xb7, 0x37, 0x61, 0x37, 0x88, 0xb6, 0xe3, 0x9f, 0x00, 0x9a, 0xd4,
	0xa0, 0xe7, 0xf0, 0xc8, 0xa5, 0xa6, 0xec, 0x07, 0xb6, 0x13, 0xd0, 0x7a, 0x0e, 0x71, 0x0d, 0xb9,
	0x09, 0xbb, 0x9d, 0x28, 0x8c, 0xff, 0x0a, 0x60, 0xd7, 0x69, 0xd0, 0x53, 0x00, 0x5e, 0x10, 0x9a,
	0xe7, 0x9a, 0x95, 0xa5, 0xad, 0xea, 0x3e, 0xde, 0xe3, 0xc5, 0xcf, 0x4e, 0x50, 0xed, 0x8e, 0xca,
	0x03, 0xeb, 0xf5, 0x0e, 0xb6, 0xe7, 0x0d, 0x4b, 0xa2, 0xb3, 0x61, 0x49, 0x20, 0x08, 0x6d, 0x9b,
	0x86, 0x67, 0xc1, 0x79, 0x17, 0xdb, 0xb3, 0x6b, 0xb7, 0xcb, 0x14, 0xf6, 0x5b, 0x09, 0xd4, 0x08,
	0x43, 0xcf, 0x9f, 0x2b, 0x31, 0x3a, 0x69, 0xc7, 0xb1, 0xbe, 0xd6, 0x06, 0xa7, 0x1b, 0xf5, 0xae,
	0x12, 0xe7, 0xc1, 0xb7, 0xc1, 0xe8, 0x2d, 0x7c, 0xc2, 0x55, 0x0b, 0x38, 0x3a, 0x6c, 0x9b, 0xbc,
	0xad, 0x92, 0x7f, 0x1b, 0xbc, 0xbb, 0xf0, 0xc5, 0x98, 0x2a, 0x41, 0xe5, 0x34, 0x51, 0x7a, 0x3a,
	0xb4, 0x7f, 0xa0, 0xfa, 0xb7, 0x63, 0x6f, 0x22, 0xb5, 0x1f, 0x22, 0x52, 0xb2, 0xb8, 0x48, 0x77,
	0x6d, 0xe1, 0xbe, 0xfb, 0x2f, 0x00, 0x00, 0xff, 0xff, 0x3f, 0x47, 0x55, 0xac, 0xab, 0x06, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// LoadBalancerClient is the client API for LoadBalancer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LoadBalancerClient interface {
	// Bidirectional rpc to get a list of servers.
	BalanceLoad(ctx context.Context, opts ...grpc.CallOption) (LoadBalancer_BalanceLoadClient, error)
}

type loadBalancerClient struct {
	cc grpc.ClientConnInterface
}

func NewLoadBalancerClient(cc grpc.ClientConnInterface) LoadBalancerClient {
	return &loadBalancerClient{cc}
}

func (c *loadBalancerClient) BalanceLoad(ctx context.Context, opts ...grpc.CallOption) (LoadBalancer_BalanceLoadClient, error) {
	stream, err := c.cc.NewStream(ctx, &_LoadBalancer_serviceDesc.Streams[0], "/grpc.lb.v1.LoadBalancer/BalanceLoad", opts...)
	if err != nil {
		return nil, err
	}
	x := &loadBalancerBalanceLoadClient{stream}
	return x, nil
}

type LoadBalancer_BalanceLoadClient interface {
	Send(*LoadBalanceRequest) error
	Recv() (*LoadBalanceResponse, error)
	grpc.ClientStream
}

type loadBalancerBalanceLoadClient struct {
	grpc.ClientStream
}

func (x *loadBalancerBalanceLoadClient) Send(m *LoadBalanceRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *loadBalancerBalanceLoadClient) Recv() (*LoadBalanceResponse, error) {
	m := new(LoadBalanceResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// LoadBalancerServer is the server API for LoadBalancer service.
type LoadBalancerServer interface {
	// Bidirectional rpc to get a list of servers.
	BalanceLoad(LoadBalancer_BalanceLoadServer) error
}

// UnimplementedLoadBalancerServer can be embedded to have forward compatible implementations.
type UnimplementedLoadBalancerServer struct {
}

func (*UnimplementedLoadBalancerServer) BalanceLoad(srv LoadBalancer_BalanceLoadServer) error {
	return status.Errorf(codes.Unimplemented, "method BalanceLoad not implemented")
}

func RegisterLoadBalancerServer(s *grpc.Server, srv LoadBalancerServer) {
	s.RegisterService(&_LoadBalancer_serviceDesc, srv)
}

func _LoadBalancer_BalanceLoad_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(LoadBalancerServer).BalanceLoad(&loadBalancerBalanceLoadServer{stream})
}

type LoadBalancer_BalanceLoadServer interface {
	Send(*LoadBalanceResponse) error
	Recv() (*LoadBalanceRequest, error)
	grpc.ServerStream
}

type loadBalancerBalanceLoadServer struct {
	grpc.ServerStream
}

func (x *loadBalancerBalanceLoadServer) Send(m *LoadBalanceResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *loadBalancerBalanceLoadServer) Recv() (*LoadBalanceRequest, error) {
	m := new(LoadBalanceRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _LoadBalancer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.lb.v1.LoadBalancer",
	HandlerType: (*LoadBalancerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "BalanceLoad",
			Handler:       _LoadBalancer_BalanceLoad_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "grpc/lb/v1/load_balancer.proto",
}
