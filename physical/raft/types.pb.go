// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: physical/raft/types.proto

package raft

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type LogOperation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// OpType is the Operation type
	OpType uint32 `protobuf:"varint,1,opt,name=op_type,json=opType,proto3" json:"op_type,omitempty"`
	// Flags is an opaque value, currently unused. Reserved.
	Flags uint64 `protobuf:"varint,2,opt,name=flags,proto3" json:"flags,omitempty"`
	// Key that is being affected
	Key string `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	// Value is optional, corresponds to the key
	Value []byte `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *LogOperation) Reset() {
	*x = LogOperation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_physical_raft_types_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogOperation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogOperation) ProtoMessage() {}

func (x *LogOperation) ProtoReflect() protoreflect.Message {
	mi := &file_physical_raft_types_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogOperation.ProtoReflect.Descriptor instead.
func (*LogOperation) Descriptor() ([]byte, []int) {
	return file_physical_raft_types_proto_rawDescGZIP(), []int{0}
}

func (x *LogOperation) GetOpType() uint32 {
	if x != nil {
		return x.OpType
	}
	return 0
}

func (x *LogOperation) GetFlags() uint64 {
	if x != nil {
		return x.Flags
	}
	return 0
}

func (x *LogOperation) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *LogOperation) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type LogData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Operations []*LogOperation `protobuf:"bytes,1,rep,name=operations,proto3" json:"operations,omitempty"`
}

func (x *LogData) Reset() {
	*x = LogData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_physical_raft_types_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogData) ProtoMessage() {}

func (x *LogData) ProtoReflect() protoreflect.Message {
	mi := &file_physical_raft_types_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogData.ProtoReflect.Descriptor instead.
func (*LogData) Descriptor() ([]byte, []int) {
	return file_physical_raft_types_proto_rawDescGZIP(), []int{1}
}

func (x *LogData) GetOperations() []*LogOperation {
	if x != nil {
		return x.Operations
	}
	return nil
}

type IndexValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term  uint64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	Index uint64 `protobuf:"varint,2,opt,name=index,proto3" json:"index,omitempty"`
}

func (x *IndexValue) Reset() {
	*x = IndexValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_physical_raft_types_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IndexValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IndexValue) ProtoMessage() {}

func (x *IndexValue) ProtoReflect() protoreflect.Message {
	mi := &file_physical_raft_types_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IndexValue.ProtoReflect.Descriptor instead.
func (*IndexValue) Descriptor() ([]byte, []int) {
	return file_physical_raft_types_proto_rawDescGZIP(), []int{2}
}

func (x *IndexValue) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *IndexValue) GetIndex() uint64 {
	if x != nil {
		return x.Index
	}
	return 0
}

type Server struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Suffrage int32  `protobuf:"varint,1,opt,name=suffrage,proto3" json:"suffrage,omitempty"`
	Id       string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Address  string `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty"`
}

func (x *Server) Reset() {
	*x = Server{}
	if protoimpl.UnsafeEnabled {
		mi := &file_physical_raft_types_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Server) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Server) ProtoMessage() {}

func (x *Server) ProtoReflect() protoreflect.Message {
	mi := &file_physical_raft_types_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Server.ProtoReflect.Descriptor instead.
func (*Server) Descriptor() ([]byte, []int) {
	return file_physical_raft_types_proto_rawDescGZIP(), []int{3}
}

func (x *Server) GetSuffrage() int32 {
	if x != nil {
		return x.Suffrage
	}
	return 0
}

func (x *Server) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Server) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type ConfigurationValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index   uint64    `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	Servers []*Server `protobuf:"bytes,2,rep,name=servers,proto3" json:"servers,omitempty"`
}

func (x *ConfigurationValue) Reset() {
	*x = ConfigurationValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_physical_raft_types_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigurationValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigurationValue) ProtoMessage() {}

func (x *ConfigurationValue) ProtoReflect() protoreflect.Message {
	mi := &file_physical_raft_types_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigurationValue.ProtoReflect.Descriptor instead.
func (*ConfigurationValue) Descriptor() ([]byte, []int) {
	return file_physical_raft_types_proto_rawDescGZIP(), []int{4}
}

func (x *ConfigurationValue) GetIndex() uint64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *ConfigurationValue) GetServers() []*Server {
	if x != nil {
		return x.Servers
	}
	return nil
}

type LocalNodeConfigValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DesiredSuffrage string `protobuf:"bytes,1,opt,name=desired_suffrage,json=desiredSuffrage,proto3" json:"desired_suffrage,omitempty"`
}

func (x *LocalNodeConfigValue) Reset() {
	*x = LocalNodeConfigValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_physical_raft_types_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LocalNodeConfigValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocalNodeConfigValue) ProtoMessage() {}

func (x *LocalNodeConfigValue) ProtoReflect() protoreflect.Message {
	mi := &file_physical_raft_types_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocalNodeConfigValue.ProtoReflect.Descriptor instead.
func (*LocalNodeConfigValue) Descriptor() ([]byte, []int) {
	return file_physical_raft_types_proto_rawDescGZIP(), []int{5}
}

func (x *LocalNodeConfigValue) GetDesiredSuffrage() string {
	if x != nil {
		return x.DesiredSuffrage
	}
	return ""
}

var File_physical_raft_types_proto protoreflect.FileDescriptor

var file_physical_raft_types_proto_rawDesc = []byte{
	0x0a, 0x19, 0x70, 0x68, 0x79, 0x73, 0x69, 0x63, 0x61, 0x6c, 0x2f, 0x72, 0x61, 0x66, 0x74, 0x2f,
	0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x72, 0x61, 0x66,
	0x74, 0x22, 0x65, 0x0a, 0x0c, 0x4c, 0x6f, 0x67, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x17, 0x0a, 0x07, 0x6f, 0x70, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x06, 0x6f, 0x70, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x6c,
	0x61, 0x67, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x66, 0x6c, 0x61, 0x67, 0x73,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x3d, 0x0a, 0x07, 0x4c, 0x6f, 0x67, 0x44,
	0x61, 0x74, 0x61, 0x12, 0x32, 0x0a, 0x0a, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x4c,
	0x6f, 0x67, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x6f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x36, 0x0a, 0x0a, 0x49, 0x6e, 0x64, 0x65, 0x78,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x22,
	0x4e, 0x0a, 0x06, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x75, 0x66,
	0x66, 0x72, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x73, 0x75, 0x66,
	0x66, 0x72, 0x61, 0x67, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22,
	0x52, 0x0a, 0x12, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x26, 0x0a, 0x07, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x72,
	0x61, 0x66, 0x74, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x07, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x73, 0x22, 0x41, 0x0a, 0x14, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x4e, 0x6f, 0x64, 0x65,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x29, 0x0a, 0x10, 0x64,
	0x65, 0x73, 0x69, 0x72, 0x65, 0x64, 0x5f, 0x73, 0x75, 0x66, 0x66, 0x72, 0x61, 0x67, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x64, 0x65, 0x73, 0x69, 0x72, 0x65, 0x64, 0x53, 0x75,
	0x66, 0x66, 0x72, 0x61, 0x67, 0x65, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x61, 0x73, 0x68, 0x69, 0x63, 0x6f, 0x72, 0x70, 0x2f, 0x76,
	0x61, 0x75, 0x6c, 0x74, 0x2f, 0x70, 0x68, 0x79, 0x73, 0x69, 0x63, 0x61, 0x6c, 0x2f, 0x72, 0x61,
	0x66, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_physical_raft_types_proto_rawDescOnce sync.Once
	file_physical_raft_types_proto_rawDescData = file_physical_raft_types_proto_rawDesc
)

func file_physical_raft_types_proto_rawDescGZIP() []byte {
	file_physical_raft_types_proto_rawDescOnce.Do(func() {
		file_physical_raft_types_proto_rawDescData = protoimpl.X.CompressGZIP(file_physical_raft_types_proto_rawDescData)
	})
	return file_physical_raft_types_proto_rawDescData
}

var file_physical_raft_types_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_physical_raft_types_proto_goTypes = []any{
	(*LogOperation)(nil),         // 0: raft.LogOperation
	(*LogData)(nil),              // 1: raft.LogData
	(*IndexValue)(nil),           // 2: raft.IndexValue
	(*Server)(nil),               // 3: raft.Server
	(*ConfigurationValue)(nil),   // 4: raft.ConfigurationValue
	(*LocalNodeConfigValue)(nil), // 5: raft.LocalNodeConfigValue
}
var file_physical_raft_types_proto_depIdxs = []int32{
	0, // 0: raft.LogData.operations:type_name -> raft.LogOperation
	3, // 1: raft.ConfigurationValue.servers:type_name -> raft.Server
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_physical_raft_types_proto_init() }
func file_physical_raft_types_proto_init() {
	if File_physical_raft_types_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_physical_raft_types_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*LogOperation); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_physical_raft_types_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*LogData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_physical_raft_types_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*IndexValue); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_physical_raft_types_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*Server); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_physical_raft_types_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*ConfigurationValue); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_physical_raft_types_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*LocalNodeConfigValue); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_physical_raft_types_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_physical_raft_types_proto_goTypes,
		DependencyIndexes: file_physical_raft_types_proto_depIdxs,
		MessageInfos:      file_physical_raft_types_proto_msgTypes,
	}.Build()
	File_physical_raft_types_proto = out.File
	file_physical_raft_types_proto_rawDesc = nil
	file_physical_raft_types_proto_goTypes = nil
	file_physical_raft_types_proto_depIdxs = nil
}
