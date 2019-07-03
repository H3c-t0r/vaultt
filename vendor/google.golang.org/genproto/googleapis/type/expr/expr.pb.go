// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/type/expr.proto

package expr

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
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

// Represents an expression text. Example:
//
//     title: "User account presence"
//     description: "Determines whether the request has a user account"
//     expression: "size(request.user) > 0"
type Expr struct {
	// Textual representation of an expression in
	// Common Expression Language syntax.
	//
	// The application context of the containing message determines which
	// well-known feature set of CEL is supported.
	Expression string `protobuf:"bytes,1,opt,name=expression,proto3" json:"expression,omitempty"`
	// An optional title for the expression, i.e. a short string describing
	// its purpose. This can be used e.g. in UIs which allow to enter the
	// expression.
	Title string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	// An optional description of the expression. This is a longer text which
	// describes the expression, e.g. when hovered over it in a UI.
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	// An optional string indicating the location of the expression for error
	// reporting, e.g. a file name and a position in the file.
	Location             string   `protobuf:"bytes,4,opt,name=location,proto3" json:"location,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Expr) Reset()         { *m = Expr{} }
func (m *Expr) String() string { return proto.CompactTextString(m) }
func (*Expr) ProtoMessage()    {}
func (*Expr) Descriptor() ([]byte, []int) {
	return fileDescriptor_d7920f1ae7a2722f, []int{0}
}

func (m *Expr) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Expr.Unmarshal(m, b)
}
func (m *Expr) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Expr.Marshal(b, m, deterministic)
}
func (m *Expr) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Expr.Merge(m, src)
}
func (m *Expr) XXX_Size() int {
	return xxx_messageInfo_Expr.Size(m)
}
func (m *Expr) XXX_DiscardUnknown() {
	xxx_messageInfo_Expr.DiscardUnknown(m)
}

var xxx_messageInfo_Expr proto.InternalMessageInfo

func (m *Expr) GetExpression() string {
	if m != nil {
		return m.Expression
	}
	return ""
}

func (m *Expr) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Expr) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Expr) GetLocation() string {
	if m != nil {
		return m.Location
	}
	return ""
}

func init() {
	proto.RegisterType((*Expr)(nil), "google.type.Expr")
}

func init() { proto.RegisterFile("google/type/expr.proto", fileDescriptor_d7920f1ae7a2722f) }

var fileDescriptor_d7920f1ae7a2722f = []byte{
	// 195 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4b, 0xcf, 0xcf, 0x4f,
	0xcf, 0x49, 0xd5, 0x2f, 0xa9, 0x2c, 0x48, 0xd5, 0x4f, 0xad, 0x28, 0x28, 0xd2, 0x2b, 0x28, 0xca,
	0x2f, 0xc9, 0x17, 0xe2, 0x86, 0x88, 0xeb, 0x81, 0xc4, 0x95, 0xaa, 0xb8, 0x58, 0x5c, 0x2b, 0x0a,
	0x8a, 0x84, 0xe4, 0xb8, 0xb8, 0x40, 0x4a, 0x52, 0x8b, 0x8b, 0x33, 0xf3, 0xf3, 0x24, 0x18, 0x15,
	0x18, 0x35, 0x38, 0x83, 0x90, 0x44, 0x84, 0x44, 0xb8, 0x58, 0x4b, 0x32, 0x4b, 0x72, 0x52, 0x25,
	0x98, 0xc0, 0x52, 0x10, 0x8e, 0x90, 0x02, 0x17, 0x77, 0x4a, 0x6a, 0x71, 0x72, 0x51, 0x66, 0x41,
	0x09, 0x48, 0x1b, 0x33, 0x58, 0x0e, 0x59, 0x48, 0x48, 0x8a, 0x8b, 0x23, 0x27, 0x3f, 0x39, 0x11,
	0x2c, 0xcd, 0x02, 0x96, 0x86, 0xf3, 0x9d, 0xa2, 0xb8, 0xf8, 0x93, 0xf3, 0x73, 0xf5, 0x90, 0x9c,
	0xe3, 0xc4, 0x09, 0x72, 0x4c, 0x00, 0xc8, 0x99, 0x01, 0x8c, 0x51, 0x26, 0x50, 0x99, 0xf4, 0xfc,
	0x9c, 0xc4, 0xbc, 0x74, 0xbd, 0xfc, 0xa2, 0x74, 0xfd, 0xf4, 0xd4, 0x3c, 0xb0, 0x27, 0xf4, 0x21,
	0x52, 0x89, 0x05, 0x99, 0xc5, 0x08, 0xff, 0x59, 0x83, 0x88, 0x45, 0x4c, 0xcc, 0xee, 0x21, 0x01,
	0x49, 0x6c, 0x60, 0x65, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xe7, 0x67, 0x9e, 0xf5, 0x05,
	0x01, 0x00, 0x00,
}
