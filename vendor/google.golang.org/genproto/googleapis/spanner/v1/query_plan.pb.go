// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/spanner/v1/query_plan.proto

package spanner

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	_struct "github.com/golang/protobuf/ptypes/struct"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// The kind of [PlanNode][google.spanner.v1.PlanNode]. Distinguishes between the two different kinds of
// nodes that can appear in a query plan.
type PlanNode_Kind int32

const (
	// Not specified.
	PlanNode_KIND_UNSPECIFIED PlanNode_Kind = 0
	// Denotes a Relational operator node in the expression tree. Relational
	// operators represent iterative processing of rows during query execution.
	// For example, a `TableScan` operation that reads rows from a table.
	PlanNode_RELATIONAL PlanNode_Kind = 1
	// Denotes a Scalar node in the expression tree. Scalar nodes represent
	// non-iterable entities in the query plan. For example, constants or
	// arithmetic operators appearing inside predicate expressions or references
	// to column names.
	PlanNode_SCALAR PlanNode_Kind = 2
)

var PlanNode_Kind_name = map[int32]string{
	0: "KIND_UNSPECIFIED",
	1: "RELATIONAL",
	2: "SCALAR",
}

var PlanNode_Kind_value = map[string]int32{
	"KIND_UNSPECIFIED": 0,
	"RELATIONAL":       1,
	"SCALAR":           2,
}

func (x PlanNode_Kind) String() string {
	return proto.EnumName(PlanNode_Kind_name, int32(x))
}

func (PlanNode_Kind) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_14e5ab0c0507d181, []int{0, 0}
}

// Node information for nodes appearing in a [QueryPlan.plan_nodes][google.spanner.v1.QueryPlan.plan_nodes].
type PlanNode struct {
	// The `PlanNode`'s index in [node list][google.spanner.v1.QueryPlan.plan_nodes].
	Index int32 `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	// Used to determine the type of node. May be needed for visualizing
	// different kinds of nodes differently. For example, If the node is a
	// [SCALAR][google.spanner.v1.PlanNode.Kind.SCALAR] node, it will have a condensed representation
	// which can be used to directly embed a description of the node in its
	// parent.
	Kind PlanNode_Kind `protobuf:"varint,2,opt,name=kind,proto3,enum=google.spanner.v1.PlanNode_Kind" json:"kind,omitempty"`
	// The display name for the node.
	DisplayName string `protobuf:"bytes,3,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	// List of child node `index`es and their relationship to this parent.
	ChildLinks []*PlanNode_ChildLink `protobuf:"bytes,4,rep,name=child_links,json=childLinks,proto3" json:"child_links,omitempty"`
	// Condensed representation for [SCALAR][google.spanner.v1.PlanNode.Kind.SCALAR] nodes.
	ShortRepresentation *PlanNode_ShortRepresentation `protobuf:"bytes,5,opt,name=short_representation,json=shortRepresentation,proto3" json:"short_representation,omitempty"`
	// Attributes relevant to the node contained in a group of key-value pairs.
	// For example, a Parameter Reference node could have the following
	// information in its metadata:
	//
	//     {
	//       "parameter_reference": "param1",
	//       "parameter_type": "array"
	//     }
	Metadata *_struct.Struct `protobuf:"bytes,6,opt,name=metadata,proto3" json:"metadata,omitempty"`
	// The execution statistics associated with the node, contained in a group of
	// key-value pairs. Only present if the plan was returned as a result of a
	// profile query. For example, number of executions, number of rows/time per
	// execution etc.
	ExecutionStats       *_struct.Struct `protobuf:"bytes,7,opt,name=execution_stats,json=executionStats,proto3" json:"execution_stats,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *PlanNode) Reset()         { *m = PlanNode{} }
func (m *PlanNode) String() string { return proto.CompactTextString(m) }
func (*PlanNode) ProtoMessage()    {}
func (*PlanNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_14e5ab0c0507d181, []int{0}
}

func (m *PlanNode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlanNode.Unmarshal(m, b)
}
func (m *PlanNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlanNode.Marshal(b, m, deterministic)
}
func (m *PlanNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlanNode.Merge(m, src)
}
func (m *PlanNode) XXX_Size() int {
	return xxx_messageInfo_PlanNode.Size(m)
}
func (m *PlanNode) XXX_DiscardUnknown() {
	xxx_messageInfo_PlanNode.DiscardUnknown(m)
}

var xxx_messageInfo_PlanNode proto.InternalMessageInfo

func (m *PlanNode) GetIndex() int32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *PlanNode) GetKind() PlanNode_Kind {
	if m != nil {
		return m.Kind
	}
	return PlanNode_KIND_UNSPECIFIED
}

func (m *PlanNode) GetDisplayName() string {
	if m != nil {
		return m.DisplayName
	}
	return ""
}

func (m *PlanNode) GetChildLinks() []*PlanNode_ChildLink {
	if m != nil {
		return m.ChildLinks
	}
	return nil
}

func (m *PlanNode) GetShortRepresentation() *PlanNode_ShortRepresentation {
	if m != nil {
		return m.ShortRepresentation
	}
	return nil
}

func (m *PlanNode) GetMetadata() *_struct.Struct {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *PlanNode) GetExecutionStats() *_struct.Struct {
	if m != nil {
		return m.ExecutionStats
	}
	return nil
}

// Metadata associated with a parent-child relationship appearing in a
// [PlanNode][google.spanner.v1.PlanNode].
type PlanNode_ChildLink struct {
	// The node to which the link points.
	ChildIndex int32 `protobuf:"varint,1,opt,name=child_index,json=childIndex,proto3" json:"child_index,omitempty"`
	// The type of the link. For example, in Hash Joins this could be used to
	// distinguish between the build child and the probe child, or in the case
	// of the child being an output variable, to represent the tag associated
	// with the output variable.
	Type string `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	// Only present if the child node is [SCALAR][google.spanner.v1.PlanNode.Kind.SCALAR] and corresponds
	// to an output variable of the parent node. The field carries the name of
	// the output variable.
	// For example, a `TableScan` operator that reads rows from a table will
	// have child links to the `SCALAR` nodes representing the output variables
	// created for each column that is read by the operator. The corresponding
	// `variable` fields will be set to the variable names assigned to the
	// columns.
	Variable             string   `protobuf:"bytes,3,opt,name=variable,proto3" json:"variable,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PlanNode_ChildLink) Reset()         { *m = PlanNode_ChildLink{} }
func (m *PlanNode_ChildLink) String() string { return proto.CompactTextString(m) }
func (*PlanNode_ChildLink) ProtoMessage()    {}
func (*PlanNode_ChildLink) Descriptor() ([]byte, []int) {
	return fileDescriptor_14e5ab0c0507d181, []int{0, 0}
}

func (m *PlanNode_ChildLink) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlanNode_ChildLink.Unmarshal(m, b)
}
func (m *PlanNode_ChildLink) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlanNode_ChildLink.Marshal(b, m, deterministic)
}
func (m *PlanNode_ChildLink) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlanNode_ChildLink.Merge(m, src)
}
func (m *PlanNode_ChildLink) XXX_Size() int {
	return xxx_messageInfo_PlanNode_ChildLink.Size(m)
}
func (m *PlanNode_ChildLink) XXX_DiscardUnknown() {
	xxx_messageInfo_PlanNode_ChildLink.DiscardUnknown(m)
}

var xxx_messageInfo_PlanNode_ChildLink proto.InternalMessageInfo

func (m *PlanNode_ChildLink) GetChildIndex() int32 {
	if m != nil {
		return m.ChildIndex
	}
	return 0
}

func (m *PlanNode_ChildLink) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *PlanNode_ChildLink) GetVariable() string {
	if m != nil {
		return m.Variable
	}
	return ""
}

// Condensed representation of a node and its subtree. Only present for
// `SCALAR` [PlanNode(s)][google.spanner.v1.PlanNode].
type PlanNode_ShortRepresentation struct {
	// A string representation of the expression subtree rooted at this node.
	Description string `protobuf:"bytes,1,opt,name=description,proto3" json:"description,omitempty"`
	// A mapping of (subquery variable name) -> (subquery node id) for cases
	// where the `description` string of this node references a `SCALAR`
	// subquery contained in the expression subtree rooted at this node. The
	// referenced `SCALAR` subquery may not necessarily be a direct child of
	// this node.
	Subqueries           map[string]int32 `protobuf:"bytes,2,rep,name=subqueries,proto3" json:"subqueries,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *PlanNode_ShortRepresentation) Reset()         { *m = PlanNode_ShortRepresentation{} }
func (m *PlanNode_ShortRepresentation) String() string { return proto.CompactTextString(m) }
func (*PlanNode_ShortRepresentation) ProtoMessage()    {}
func (*PlanNode_ShortRepresentation) Descriptor() ([]byte, []int) {
	return fileDescriptor_14e5ab0c0507d181, []int{0, 1}
}

func (m *PlanNode_ShortRepresentation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlanNode_ShortRepresentation.Unmarshal(m, b)
}
func (m *PlanNode_ShortRepresentation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlanNode_ShortRepresentation.Marshal(b, m, deterministic)
}
func (m *PlanNode_ShortRepresentation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlanNode_ShortRepresentation.Merge(m, src)
}
func (m *PlanNode_ShortRepresentation) XXX_Size() int {
	return xxx_messageInfo_PlanNode_ShortRepresentation.Size(m)
}
func (m *PlanNode_ShortRepresentation) XXX_DiscardUnknown() {
	xxx_messageInfo_PlanNode_ShortRepresentation.DiscardUnknown(m)
}

var xxx_messageInfo_PlanNode_ShortRepresentation proto.InternalMessageInfo

func (m *PlanNode_ShortRepresentation) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *PlanNode_ShortRepresentation) GetSubqueries() map[string]int32 {
	if m != nil {
		return m.Subqueries
	}
	return nil
}

// Contains an ordered list of nodes appearing in the query plan.
type QueryPlan struct {
	// The nodes in the query plan. Plan nodes are returned in pre-order starting
	// with the plan root. Each [PlanNode][google.spanner.v1.PlanNode]'s `id` corresponds to its index in
	// `plan_nodes`.
	PlanNodes            []*PlanNode `protobuf:"bytes,1,rep,name=plan_nodes,json=planNodes,proto3" json:"plan_nodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *QueryPlan) Reset()         { *m = QueryPlan{} }
func (m *QueryPlan) String() string { return proto.CompactTextString(m) }
func (*QueryPlan) ProtoMessage()    {}
func (*QueryPlan) Descriptor() ([]byte, []int) {
	return fileDescriptor_14e5ab0c0507d181, []int{1}
}

func (m *QueryPlan) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryPlan.Unmarshal(m, b)
}
func (m *QueryPlan) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryPlan.Marshal(b, m, deterministic)
}
func (m *QueryPlan) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryPlan.Merge(m, src)
}
func (m *QueryPlan) XXX_Size() int {
	return xxx_messageInfo_QueryPlan.Size(m)
}
func (m *QueryPlan) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryPlan.DiscardUnknown(m)
}

var xxx_messageInfo_QueryPlan proto.InternalMessageInfo

func (m *QueryPlan) GetPlanNodes() []*PlanNode {
	if m != nil {
		return m.PlanNodes
	}
	return nil
}

func init() {
	proto.RegisterEnum("google.spanner.v1.PlanNode_Kind", PlanNode_Kind_name, PlanNode_Kind_value)
	proto.RegisterType((*PlanNode)(nil), "google.spanner.v1.PlanNode")
	proto.RegisterType((*PlanNode_ChildLink)(nil), "google.spanner.v1.PlanNode.ChildLink")
	proto.RegisterType((*PlanNode_ShortRepresentation)(nil), "google.spanner.v1.PlanNode.ShortRepresentation")
	proto.RegisterMapType((map[string]int32)(nil), "google.spanner.v1.PlanNode.ShortRepresentation.SubqueriesEntry")
	proto.RegisterType((*QueryPlan)(nil), "google.spanner.v1.QueryPlan")
}

func init() { proto.RegisterFile("google/spanner/v1/query_plan.proto", fileDescriptor_14e5ab0c0507d181) }

var fileDescriptor_14e5ab0c0507d181 = []byte{
	// 604 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xdf, 0x6e, 0xd3, 0x3e,
	0x18, 0xfd, 0xa5, 0x6b, 0xfb, 0x5b, 0xbf, 0xa2, 0xae, 0x78, 0x43, 0x8b, 0x0a, 0x12, 0xa1, 0x12,
	0x52, 0xaf, 0x12, 0x75, 0xe3, 0x62, 0x1a, 0x42, 0xd0, 0x75, 0xdd, 0x54, 0xad, 0x2a, 0xc5, 0x81,
	0x5d, 0xa0, 0x4a, 0x91, 0xdb, 0x98, 0xce, 0x6a, 0x6a, 0x87, 0x38, 0xa9, 0xd6, 0x97, 0xe0, 0x96,
	0x77, 0xe0, 0x51, 0x78, 0x01, 0x5e, 0x07, 0xd9, 0xf9, 0xc3, 0x60, 0xa8, 0x12, 0x77, 0x9f, 0x7d,
	0xce, 0x77, 0xe2, 0xef, 0x1c, 0x3b, 0xd0, 0x5e, 0x08, 0xb1, 0x08, 0xa8, 0x23, 0x43, 0xc2, 0x39,
	0x8d, 0x9c, 0x75, 0xd7, 0xf9, 0x9c, 0xd0, 0x68, 0xe3, 0x85, 0x01, 0xe1, 0x76, 0x18, 0x89, 0x58,
	0xa0, 0x87, 0x29, 0xc7, 0xce, 0x38, 0xf6, 0xba, 0xdb, 0x7a, 0x92, 0xb5, 0x69, 0xc2, 0x2c, 0xf9,
	0xe4, 0xc8, 0x38, 0x4a, 0xe6, 0x71, 0xda, 0x50, 0xa0, 0x24, 0x64, 0x0e, 0xe1, 0x5c, 0xc4, 0x24,
	0x66, 0x82, 0xcb, 0x14, 0x6d, 0x7f, 0xa9, 0xc2, 0xee, 0x24, 0x20, 0x7c, 0x2c, 0x7c, 0x8a, 0x0e,
	0xa0, 0xc2, 0xb8, 0x4f, 0x6f, 0x4d, 0xc3, 0x32, 0x3a, 0x15, 0x9c, 0x2e, 0xd0, 0x0b, 0x28, 0x2f,
	0x19, 0xf7, 0xcd, 0x92, 0x65, 0x74, 0x1a, 0x47, 0x96, 0x7d, 0xef, 0x00, 0x76, 0x2e, 0x60, 0x5f,
	0x31, 0xee, 0x63, 0xcd, 0x46, 0xcf, 0xe0, 0x81, 0xcf, 0x64, 0x18, 0x90, 0x8d, 0xc7, 0xc9, 0x8a,
	0x9a, 0x3b, 0x96, 0xd1, 0xa9, 0xe1, 0x7a, 0xb6, 0x37, 0x26, 0x2b, 0x8a, 0x2e, 0xa0, 0x3e, 0xbf,
	0x61, 0x81, 0xef, 0x05, 0x8c, 0x2f, 0xa5, 0x59, 0xb6, 0x76, 0x3a, 0xf5, 0xa3, 0xe7, 0xdb, 0xf4,
	0xfb, 0x8a, 0x3e, 0x62, 0x7c, 0x89, 0x61, 0x9e, 0x97, 0x12, 0xcd, 0xe0, 0x40, 0xde, 0x88, 0x28,
	0xf6, 0x22, 0x1a, 0x46, 0x54, 0x52, 0x9e, 0x8e, 0x68, 0x56, 0x2c, 0xa3, 0x53, 0x3f, 0x72, 0xb6,
	0x09, 0xba, 0xaa, 0x0f, 0xff, 0xd6, 0x86, 0xf7, 0xe5, 0xfd, 0x4d, 0x74, 0x0c, 0xbb, 0x2b, 0x1a,
	0x13, 0x9f, 0xc4, 0xc4, 0xac, 0x6a, 0xdd, 0xc3, 0x5c, 0x37, 0xb7, 0xdd, 0x76, 0xb5, 0xed, 0xb8,
	0x20, 0xa2, 0x37, 0xb0, 0x47, 0x6f, 0xe9, 0x3c, 0x51, 0x0a, 0x9e, 0x8c, 0x49, 0x2c, 0xcd, 0xff,
	0xb7, 0xf7, 0x36, 0x0a, 0xbe, 0xab, 0xe8, 0xad, 0x29, 0xd4, 0x8a, 0x99, 0xd1, 0xd3, 0xdc, 0xaf,
	0xbb, 0x21, 0xa5, 0x46, 0x0c, 0x75, 0x52, 0x08, 0xca, 0xf1, 0x26, 0xa4, 0x3a, 0xa9, 0x1a, 0xd6,
	0x35, 0x6a, 0xc1, 0xee, 0x9a, 0x44, 0x8c, 0xcc, 0x82, 0x3c, 0x83, 0x62, 0xdd, 0xfa, 0x61, 0xc0,
	0xfe, 0x5f, 0x1c, 0x40, 0x16, 0xd4, 0x7d, 0x2a, 0xe7, 0x11, 0x0b, 0xb5, 0x8f, 0x46, 0x16, 0xdd,
	0xaf, 0x2d, 0xe4, 0x01, 0xc8, 0x64, 0xa6, 0x2e, 0x27, 0xa3, 0xd2, 0x2c, 0xe9, 0xe4, 0x5e, 0xff,
	0xa3, 0xd1, 0xb6, 0x5b, 0x28, 0x0c, 0x78, 0x1c, 0x6d, 0xf0, 0x1d, 0xc9, 0xd6, 0x2b, 0xd8, 0xfb,
	0x03, 0x46, 0x4d, 0xd8, 0x59, 0xd2, 0x4d, 0x76, 0x1a, 0x55, 0xaa, 0xfb, 0xba, 0x26, 0x41, 0x92,
	0x0e, 0x5c, 0xc1, 0xe9, 0xe2, 0xb4, 0x74, 0x62, 0xb4, 0x4f, 0xa0, 0xac, 0xee, 0x22, 0x3a, 0x80,
	0xe6, 0xd5, 0x70, 0x7c, 0xee, 0x7d, 0x18, 0xbb, 0x93, 0x41, 0x7f, 0x78, 0x31, 0x1c, 0x9c, 0x37,
	0xff, 0x43, 0x0d, 0x00, 0x3c, 0x18, 0xf5, 0xde, 0x0f, 0xdf, 0x8e, 0x7b, 0xa3, 0xa6, 0x81, 0x00,
	0xaa, 0x6e, 0xbf, 0x37, 0xea, 0xe1, 0x66, 0xa9, 0x7d, 0x09, 0xb5, 0x77, 0xea, 0xcd, 0xa9, 0x93,
	0xa3, 0x53, 0x00, 0xf5, 0xf4, 0x3c, 0x2e, 0x7c, 0x2a, 0x4d, 0x43, 0x8f, 0xf9, 0x78, 0xcb, 0x98,
	0xb8, 0x16, 0x66, 0x95, 0x3c, 0xfb, 0x6a, 0xc0, 0xa3, 0xb9, 0x58, 0xdd, 0x67, 0x9f, 0x35, 0x8a,
	0x0f, 0x4c, 0x54, 0xfc, 0x13, 0xe3, 0xe3, 0x49, 0x46, 0x5a, 0x88, 0x80, 0xf0, 0x85, 0x2d, 0xa2,
	0x85, 0xb3, 0xa0, 0x5c, 0x5f, 0x0e, 0x27, 0x85, 0x48, 0xc8, 0xe4, 0x9d, 0xff, 0xc2, 0xcb, 0xac,
	0xfc, 0x56, 0x3a, 0xbc, 0x4c, 0x5b, 0xfb, 0x81, 0x48, 0x7c, 0xdb, 0xcd, 0xbe, 0x72, 0xdd, 0xfd,
	0x9e, 0x23, 0x53, 0x8d, 0x4c, 0x33, 0x64, 0x7a, 0xdd, 0x9d, 0x55, 0xb5, 0xf0, 0xf1, 0xcf, 0x00,
	0x00, 0x00, 0xff, 0xff, 0xd6, 0xe5, 0x2e, 0x1d, 0x6f, 0x04, 0x00, 0x00,
}
