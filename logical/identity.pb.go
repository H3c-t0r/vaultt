// Code generated by protoc-gen-go. DO NOT EDIT.
// source: logical/identity.proto

package logical

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

type Entity struct {
	// ID is the unique identifier for the entity
	ID string `sentinel:"" protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	// Name is the human-friendly unique identifier for the entity
	Name string `sentinel:"" protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	// Aliases contains thhe alias mappings for the given entity
	Aliases              []*Alias `sentinel:"" protobuf:"bytes,3,rep,name=aliases" json:"aliases,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Entity) Reset()         { *m = Entity{} }
func (m *Entity) String() string { return proto.CompactTextString(m) }
func (*Entity) ProtoMessage()    {}
func (*Entity) Descriptor() ([]byte, []int) {
	return fileDescriptor_identity_67c761c38836b2d0, []int{0}
}
func (m *Entity) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Entity.Unmarshal(m, b)
}
func (m *Entity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Entity.Marshal(b, m, deterministic)
}
func (dst *Entity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Entity.Merge(dst, src)
}
func (m *Entity) XXX_Size() int {
	return xxx_messageInfo_Entity.Size(m)
}
func (m *Entity) XXX_DiscardUnknown() {
	xxx_messageInfo_Entity.DiscardUnknown(m)
}

var xxx_messageInfo_Entity proto.InternalMessageInfo

func (m *Entity) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *Entity) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Entity) GetAliases() []*Alias {
	if m != nil {
		return m.Aliases
	}
	return nil
}

type Alias struct {
	// MountType is the backend mount's type to which this identity belongs
	MountType string `sentinel:"" protobuf:"bytes,1,opt,name=mount_type,json=mountType" json:"mount_type,omitempty"`
	// MountAccessor is the identifier of the mount entry to which this
	// identity belongs
	MountAccessor string `sentinel:"" protobuf:"bytes,2,opt,name=mount_accessor,json=mountAccessor" json:"mount_accessor,omitempty"`
	// Name is the identifier of this identity in its authentication source
	Name                 string   `sentinel:"" protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Alias) Reset()         { *m = Alias{} }
func (m *Alias) String() string { return proto.CompactTextString(m) }
func (*Alias) ProtoMessage()    {}
func (*Alias) Descriptor() ([]byte, []int) {
	return fileDescriptor_identity_67c761c38836b2d0, []int{1}
}
func (m *Alias) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Alias.Unmarshal(m, b)
}
func (m *Alias) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Alias.Marshal(b, m, deterministic)
}
func (dst *Alias) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Alias.Merge(dst, src)
}
func (m *Alias) XXX_Size() int {
	return xxx_messageInfo_Alias.Size(m)
}
func (m *Alias) XXX_DiscardUnknown() {
	xxx_messageInfo_Alias.DiscardUnknown(m)
}

var xxx_messageInfo_Alias proto.InternalMessageInfo

func (m *Alias) GetMountType() string {
	if m != nil {
		return m.MountType
	}
	return ""
}

func (m *Alias) GetMountAccessor() string {
	if m != nil {
		return m.MountAccessor
	}
	return ""
}

func (m *Alias) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*Entity)(nil), "logical.Entity")
	proto.RegisterType((*Alias)(nil), "logical.Alias")
}

func init() { proto.RegisterFile("logical/identity.proto", fileDescriptor_identity_67c761c38836b2d0) }

var fileDescriptor_identity_67c761c38836b2d0 = []byte{
	// 179 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xcb, 0xc9, 0x4f, 0xcf,
	0x4c, 0x4e, 0xcc, 0xd1, 0xcf, 0x4c, 0x49, 0xcd, 0x2b, 0xc9, 0x2c, 0xa9, 0xd4, 0x2b, 0x28, 0xca,
	0x2f, 0xc9, 0x17, 0x62, 0x87, 0x8a, 0x2b, 0x85, 0x71, 0xb1, 0xb9, 0x82, 0x25, 0x84, 0xf8, 0xb8,
	0x98, 0x32, 0x53, 0x24, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0x98, 0x32, 0x53, 0x84, 0x84, 0xb8,
	0x58, 0xf2, 0x12, 0x73, 0x53, 0x25, 0x98, 0xc0, 0x22, 0x60, 0xb6, 0x90, 0x06, 0x17, 0x7b, 0x62,
	0x4e, 0x66, 0x62, 0x71, 0x6a, 0xb1, 0x04, 0xb3, 0x02, 0xb3, 0x06, 0xb7, 0x11, 0x9f, 0x1e, 0xd4,
	0x20, 0x3d, 0x47, 0x90, 0x78, 0x10, 0x4c, 0x5a, 0x29, 0x91, 0x8b, 0x15, 0x2c, 0x22, 0x24, 0xcb,
	0xc5, 0x95, 0x9b, 0x5f, 0x9a, 0x57, 0x12, 0x5f, 0x52, 0x59, 0x90, 0x0a, 0x35, 0x9e, 0x13, 0x2c,
	0x12, 0x52, 0x59, 0x90, 0x2a, 0xa4, 0xca, 0xc5, 0x07, 0x91, 0x4e, 0x4c, 0x4e, 0x4e, 0x2d, 0x2e,
	0xce, 0x2f, 0x82, 0xda, 0xc7, 0x0b, 0x16, 0x75, 0x84, 0x0a, 0xc2, 0x1d, 0xc3, 0x8c, 0x70, 0x4c,
	0x12, 0x1b, 0xd8, 0x2b, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xea, 0xd0, 0x42, 0xf9, 0xe4,
	0x00, 0x00, 0x00,
}
