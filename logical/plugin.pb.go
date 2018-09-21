// Code generated by protoc-gen-go. DO NOT EDIT.
// source: logical/plugin.proto

package logical // import "github.com/hashicorp/vault/logical"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type PluginEnvironment struct {
	// VaultVersion is the version of the Vault server
	VaultVersion         string   `protobuf:"bytes,1,opt,name=vault_version,json=vaultVersion,proto3" json:"vault_version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PluginEnvironment) Reset()         { *m = PluginEnvironment{} }
func (m *PluginEnvironment) String() string { return proto.CompactTextString(m) }
func (*PluginEnvironment) ProtoMessage()    {}
func (*PluginEnvironment) Descriptor() ([]byte, []int) {
	return fileDescriptor_plugin_c3e74d5a6c13acf1, []int{0}
}
func (m *PluginEnvironment) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PluginEnvironment.Unmarshal(m, b)
}
func (m *PluginEnvironment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PluginEnvironment.Marshal(b, m, deterministic)
}
func (dst *PluginEnvironment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PluginEnvironment.Merge(dst, src)
}
func (m *PluginEnvironment) XXX_Size() int {
	return xxx_messageInfo_PluginEnvironment.Size(m)
}
func (m *PluginEnvironment) XXX_DiscardUnknown() {
	xxx_messageInfo_PluginEnvironment.DiscardUnknown(m)
}

var xxx_messageInfo_PluginEnvironment proto.InternalMessageInfo

func (m *PluginEnvironment) GetVaultVersion() string {
	if m != nil {
		return m.VaultVersion
	}
	return ""
}

func init() {
	proto.RegisterType((*PluginEnvironment)(nil), "logical.PluginEnvironment")
}

func init() { proto.RegisterFile("logical/plugin.proto", fileDescriptor_plugin_c3e74d5a6c13acf1) }

var fileDescriptor_plugin_c3e74d5a6c13acf1 = []byte{
	// 133 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xc9, 0xc9, 0x4f, 0xcf,
	0x4c, 0x4e, 0xcc, 0xd1, 0x2f, 0xc8, 0x29, 0x4d, 0xcf, 0xcc, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0x62, 0x87, 0x8a, 0x2a, 0x59, 0x70, 0x09, 0x06, 0x80, 0x25, 0x5c, 0xf3, 0xca, 0x32, 0x8b,
	0xf2, 0xf3, 0x72, 0x53, 0xf3, 0x4a, 0x84, 0x94, 0xb9, 0x78, 0xcb, 0x12, 0x4b, 0x73, 0x4a, 0xe2,
	0xcb, 0x52, 0x8b, 0x8a, 0x33, 0xf3, 0xf3, 0x24, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0x78, 0xc0,
	0x82, 0x61, 0x10, 0x31, 0x27, 0x95, 0x28, 0xa5, 0xf4, 0xcc, 0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4,
	0xfc, 0x5c, 0xfd, 0x8c, 0xc4, 0xe2, 0x8c, 0xcc, 0xe4, 0xfc, 0xa2, 0x02, 0x7d, 0xb0, 0x22, 0x7d,
	0xa8, 0xf9, 0x49, 0x6c, 0x60, 0xfb, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0xa3, 0xff, 0x48,
	0xa9, 0x87, 0x00, 0x00, 0x00,
}
