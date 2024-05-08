// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        (unknown)
// source: builtin/logical/pki/metadata.proto

package pki

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CertificateMetadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IssuerId       string                 `protobuf:"bytes,1,opt,name=issuer_id,json=issuerId,proto3" json:"issuer_id,omitempty"`
	Role           string                 `protobuf:"bytes,2,opt,name=role,proto3" json:"role,omitempty"`
	Expiration     *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=expiration,proto3" json:"expiration,omitempty"`
	ClientMetadata []byte                 `protobuf:"bytes,4,opt,name=client_metadata,json=clientMetadata,proto3,oneof" json:"client_metadata,omitempty"`
}

func (x *CertificateMetadata) Reset() {
	*x = CertificateMetadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_builtin_logical_pki_metadata_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CertificateMetadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CertificateMetadata) ProtoMessage() {}

func (x *CertificateMetadata) ProtoReflect() protoreflect.Message {
	mi := &file_builtin_logical_pki_metadata_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CertificateMetadata.ProtoReflect.Descriptor instead.
func (*CertificateMetadata) Descriptor() ([]byte, []int) {
	return file_builtin_logical_pki_metadata_proto_rawDescGZIP(), []int{0}
}

func (x *CertificateMetadata) GetIssuerId() string {
	if x != nil {
		return x.IssuerId
	}
	return ""
}

func (x *CertificateMetadata) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

func (x *CertificateMetadata) GetExpiration() *timestamppb.Timestamp {
	if x != nil {
		return x.Expiration
	}
	return nil
}

func (x *CertificateMetadata) GetClientMetadata() []byte {
	if x != nil {
		return x.ClientMetadata
	}
	return nil
}

var File_builtin_logical_pki_metadata_proto protoreflect.FileDescriptor

var file_builtin_logical_pki_metadata_proto_rawDesc = []byte{
	0x0a, 0x22, 0x62, 0x75, 0x69, 0x6c, 0x74, 0x69, 0x6e, 0x2f, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x61,
	0x6c, 0x2f, 0x70, 0x6b, 0x69, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x70, 0x6b, 0x69, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc4, 0x01, 0x0a, 0x13, 0x43,
	0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x73, 0x75, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x73, 0x73, 0x75, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x72,
	0x6f, 0x6c, 0x65, 0x12, 0x3a, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x2c, 0x0a, 0x0f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x0e, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x88, 0x01, 0x01, 0x42, 0x12, 0x0a,
	0x10, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x68, 0x61, 0x73, 0x68, 0x69, 0x63, 0x6f, 0x72, 0x70, 0x2f, 0x76, 0x61, 0x75, 0x6c, 0x74, 0x2f,
	0x62, 0x75, 0x69, 0x6c, 0x74, 0x69, 0x6e, 0x2f, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x61, 0x6c, 0x2f,
	0x70, 0x6b, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_builtin_logical_pki_metadata_proto_rawDescOnce sync.Once
	file_builtin_logical_pki_metadata_proto_rawDescData = file_builtin_logical_pki_metadata_proto_rawDesc
)

func file_builtin_logical_pki_metadata_proto_rawDescGZIP() []byte {
	file_builtin_logical_pki_metadata_proto_rawDescOnce.Do(func() {
		file_builtin_logical_pki_metadata_proto_rawDescData = protoimpl.X.CompressGZIP(file_builtin_logical_pki_metadata_proto_rawDescData)
	})
	return file_builtin_logical_pki_metadata_proto_rawDescData
}

var file_builtin_logical_pki_metadata_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_builtin_logical_pki_metadata_proto_goTypes = []interface{}{
	(*CertificateMetadata)(nil),   // 0: pki.CertificateMetadata
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_builtin_logical_pki_metadata_proto_depIdxs = []int32{
	1, // 0: pki.CertificateMetadata.expiration:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_builtin_logical_pki_metadata_proto_init() }
func file_builtin_logical_pki_metadata_proto_init() {
	if File_builtin_logical_pki_metadata_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_builtin_logical_pki_metadata_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CertificateMetadata); i {
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
	file_builtin_logical_pki_metadata_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_builtin_logical_pki_metadata_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_builtin_logical_pki_metadata_proto_goTypes,
		DependencyIndexes: file_builtin_logical_pki_metadata_proto_depIdxs,
		MessageInfos:      file_builtin_logical_pki_metadata_proto_msgTypes,
	}.Build()
	File_builtin_logical_pki_metadata_proto = out.File
	file_builtin_logical_pki_metadata_proto_rawDesc = nil
	file_builtin_logical_pki_metadata_proto_goTypes = nil
	file_builtin_logical_pki_metadata_proto_depIdxs = nil
}
