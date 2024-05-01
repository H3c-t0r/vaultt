// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.0
// 	protoc        (unknown)
// source: vault/hcp_link/proto/link_control/link_control.proto

package link_control

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

type PurgePolicyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PurgePolicyRequest) Reset() {
	*x = PurgePolicyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vault_hcp_link_proto_link_control_link_control_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PurgePolicyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PurgePolicyRequest) ProtoMessage() {}

func (x *PurgePolicyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_vault_hcp_link_proto_link_control_link_control_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PurgePolicyRequest.ProtoReflect.Descriptor instead.
func (*PurgePolicyRequest) Descriptor() ([]byte, []int) {
	return file_vault_hcp_link_proto_link_control_link_control_proto_rawDescGZIP(), []int{0}
}

type PurgePolicyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PurgePolicyResponse) Reset() {
	*x = PurgePolicyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vault_hcp_link_proto_link_control_link_control_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PurgePolicyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PurgePolicyResponse) ProtoMessage() {}

func (x *PurgePolicyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_vault_hcp_link_proto_link_control_link_control_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PurgePolicyResponse.ProtoReflect.Descriptor instead.
func (*PurgePolicyResponse) Descriptor() ([]byte, []int) {
	return file_vault_hcp_link_proto_link_control_link_control_proto_rawDescGZIP(), []int{1}
}

var File_vault_hcp_link_proto_link_control_link_control_proto protoreflect.FileDescriptor

var file_vault_hcp_link_proto_link_control_link_control_proto_rawDesc = []byte{
	0x0a, 0x34, 0x76, 0x61, 0x75, 0x6c, 0x74, 0x2f, 0x68, 0x63, 0x70, 0x5f, 0x6c, 0x69, 0x6e, 0x6b,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6c, 0x69, 0x6e, 0x6b, 0x5f, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x2f, 0x6c, 0x69, 0x6e, 0x6b, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x6c, 0x69, 0x6e, 0x6b, 0x5f, 0x63, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x22, 0x14, 0x0a, 0x12, 0x50, 0x75, 0x72, 0x67, 0x65, 0x50, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x15, 0x0a, 0x13, 0x50, 0x75,
	0x72, 0x67, 0x65, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x32, 0x64, 0x0a, 0x0e, 0x48, 0x43, 0x50, 0x4c, 0x69, 0x6e, 0x6b, 0x43, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x12, 0x52, 0x0a, 0x0b, 0x50, 0x75, 0x72, 0x67, 0x65, 0x50, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x12, 0x20, 0x2e, 0x6c, 0x69, 0x6e, 0x6b, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f,
	0x6c, 0x2e, 0x50, 0x75, 0x72, 0x67, 0x65, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x6c, 0x69, 0x6e, 0x6b, 0x5f, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x2e, 0x50, 0x75, 0x72, 0x67, 0x65, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x3e, 0x5a, 0x3c, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x61, 0x73, 0x68, 0x69, 0x63, 0x6f, 0x72, 0x70, 0x2f,
	0x76, 0x61, 0x75, 0x6c, 0x74, 0x2f, 0x76, 0x61, 0x75, 0x6c, 0x74, 0x2f, 0x68, 0x63, 0x70, 0x5f,
	0x6c, 0x69, 0x6e, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6c, 0x69, 0x6e, 0x6b, 0x5f,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_vault_hcp_link_proto_link_control_link_control_proto_rawDescOnce sync.Once
	file_vault_hcp_link_proto_link_control_link_control_proto_rawDescData = file_vault_hcp_link_proto_link_control_link_control_proto_rawDesc
)

func file_vault_hcp_link_proto_link_control_link_control_proto_rawDescGZIP() []byte {
	file_vault_hcp_link_proto_link_control_link_control_proto_rawDescOnce.Do(func() {
		file_vault_hcp_link_proto_link_control_link_control_proto_rawDescData = protoimpl.X.CompressGZIP(file_vault_hcp_link_proto_link_control_link_control_proto_rawDescData)
	})
	return file_vault_hcp_link_proto_link_control_link_control_proto_rawDescData
}

var file_vault_hcp_link_proto_link_control_link_control_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_vault_hcp_link_proto_link_control_link_control_proto_goTypes = []interface{}{
	(*PurgePolicyRequest)(nil),  // 0: link_control.PurgePolicyRequest
	(*PurgePolicyResponse)(nil), // 1: link_control.PurgePolicyResponse
}
var file_vault_hcp_link_proto_link_control_link_control_proto_depIdxs = []int32{
	0, // 0: link_control.HCPLinkControl.PurgePolicy:input_type -> link_control.PurgePolicyRequest
	1, // 1: link_control.HCPLinkControl.PurgePolicy:output_type -> link_control.PurgePolicyResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_vault_hcp_link_proto_link_control_link_control_proto_init() }
func file_vault_hcp_link_proto_link_control_link_control_proto_init() {
	if File_vault_hcp_link_proto_link_control_link_control_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_vault_hcp_link_proto_link_control_link_control_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PurgePolicyRequest); i {
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
		file_vault_hcp_link_proto_link_control_link_control_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PurgePolicyResponse); i {
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
			RawDescriptor: file_vault_hcp_link_proto_link_control_link_control_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_vault_hcp_link_proto_link_control_link_control_proto_goTypes,
		DependencyIndexes: file_vault_hcp_link_proto_link_control_link_control_proto_depIdxs,
		MessageInfos:      file_vault_hcp_link_proto_link_control_link_control_proto_msgTypes,
	}.Build()
	File_vault_hcp_link_proto_link_control_link_control_proto = out.File
	file_vault_hcp_link_proto_link_control_link_control_proto_rawDesc = nil
	file_vault_hcp_link_proto_link_control_link_control_proto_goTypes = nil
	file_vault_hcp_link_proto_link_control_link_control_proto_depIdxs = nil
}
