// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.22.0
// 	protoc        v3.11.2
// source: google/monitoring/v3/dropped_labels.proto

package monitoring

import (
	reflect "reflect"
	sync "sync"

	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// A set of (label, value) pairs which were dropped during aggregation, attached
// to google.api.Distribution.Exemplars in google.api.Distribution values during
// aggregation.
//
// These values are used in combination with the label values that remain on the
// aggregated Distribution timeseries to construct the full label set for the
// exemplar values.  The resulting full label set may be used to identify the
// specific task/job/instance (for example) which may be contributing to a
// long-tail, while allowing the storage savings of only storing aggregated
// distribution values for a large group.
//
// Note that there are no guarantees on ordering of the labels from
// exemplar-to-exemplar and from distribution-to-distribution in the same
// stream, and there may be duplicates.  It is up to clients to resolve any
// ambiguities.
type DroppedLabels struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Map from label to its value, for all labels dropped in any aggregation.
	Label map[string]string `protobuf:"bytes,1,rep,name=label,proto3" json:"label,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *DroppedLabels) Reset() {
	*x = DroppedLabels{}
	if protoimpl.UnsafeEnabled {
		mi := &file_google_monitoring_v3_dropped_labels_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DroppedLabels) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DroppedLabels) ProtoMessage() {}

func (x *DroppedLabels) ProtoReflect() protoreflect.Message {
	mi := &file_google_monitoring_v3_dropped_labels_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DroppedLabels.ProtoReflect.Descriptor instead.
func (*DroppedLabels) Descriptor() ([]byte, []int) {
	return file_google_monitoring_v3_dropped_labels_proto_rawDescGZIP(), []int{0}
}

func (x *DroppedLabels) GetLabel() map[string]string {
	if x != nil {
		return x.Label
	}
	return nil
}

var File_google_monitoring_v3_dropped_labels_proto protoreflect.FileDescriptor

var file_google_monitoring_v3_dropped_labels_proto_rawDesc = []byte{
	0x0a, 0x29, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x6d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72,
	0x69, 0x6e, 0x67, 0x2f, 0x76, 0x33, 0x2f, 0x64, 0x72, 0x6f, 0x70, 0x70, 0x65, 0x64, 0x5f, 0x6c,
	0x61, 0x62, 0x65, 0x6c, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x6d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x2e, 0x76,
	0x33, 0x22, 0x8f, 0x01, 0x0a, 0x0d, 0x44, 0x72, 0x6f, 0x70, 0x70, 0x65, 0x64, 0x4c, 0x61, 0x62,
	0x65, 0x6c, 0x73, 0x12, 0x44, 0x0a, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x6d, 0x6f, 0x6e, 0x69,
	0x74, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x33, 0x2e, 0x44, 0x72, 0x6f, 0x70, 0x70, 0x65,
	0x64, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x1a, 0x38, 0x0a, 0x0a, 0x4c, 0x61, 0x62,
	0x65, 0x6c, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x42, 0xca, 0x01, 0x0a, 0x18, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x6d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x33,
	0x42, 0x12, 0x44, 0x72, 0x6f, 0x70, 0x70, 0x65, 0x64, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x3e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x67,
	0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2e, 0x6f, 0x72, 0x67, 0x2f, 0x67, 0x65, 0x6e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x6d, 0x6f,
	0x6e, 0x69, 0x74, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x2f, 0x76, 0x33, 0x3b, 0x6d, 0x6f, 0x6e, 0x69,
	0x74, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0xaa, 0x02, 0x1a, 0x47, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x43, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x69, 0x6e, 0x67,
	0x2e, 0x56, 0x33, 0xca, 0x02, 0x1a, 0x47, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x5c, 0x43, 0x6c, 0x6f,
	0x75, 0x64, 0x5c, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x5c, 0x56, 0x33,
	0xea, 0x02, 0x1d, 0x47, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x3a, 0x3a, 0x43, 0x6c, 0x6f, 0x75, 0x64,
	0x3a, 0x3a, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x3a, 0x3a, 0x56, 0x33,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_google_monitoring_v3_dropped_labels_proto_rawDescOnce sync.Once
	file_google_monitoring_v3_dropped_labels_proto_rawDescData = file_google_monitoring_v3_dropped_labels_proto_rawDesc
)

func file_google_monitoring_v3_dropped_labels_proto_rawDescGZIP() []byte {
	file_google_monitoring_v3_dropped_labels_proto_rawDescOnce.Do(func() {
		file_google_monitoring_v3_dropped_labels_proto_rawDescData = protoimpl.X.CompressGZIP(file_google_monitoring_v3_dropped_labels_proto_rawDescData)
	})
	return file_google_monitoring_v3_dropped_labels_proto_rawDescData
}

var file_google_monitoring_v3_dropped_labels_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_google_monitoring_v3_dropped_labels_proto_goTypes = []interface{}{
	(*DroppedLabels)(nil), // 0: google.monitoring.v3.DroppedLabels
	nil,                   // 1: google.monitoring.v3.DroppedLabels.LabelEntry
}
var file_google_monitoring_v3_dropped_labels_proto_depIdxs = []int32{
	1, // 0: google.monitoring.v3.DroppedLabels.label:type_name -> google.monitoring.v3.DroppedLabels.LabelEntry
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_google_monitoring_v3_dropped_labels_proto_init() }
func file_google_monitoring_v3_dropped_labels_proto_init() {
	if File_google_monitoring_v3_dropped_labels_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_google_monitoring_v3_dropped_labels_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DroppedLabels); i {
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
			RawDescriptor: file_google_monitoring_v3_dropped_labels_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_google_monitoring_v3_dropped_labels_proto_goTypes,
		DependencyIndexes: file_google_monitoring_v3_dropped_labels_proto_depIdxs,
		MessageInfos:      file_google_monitoring_v3_dropped_labels_proto_msgTypes,
	}.Build()
	File_google_monitoring_v3_dropped_labels_proto = out.File
	file_google_monitoring_v3_dropped_labels_proto_rawDesc = nil
	file_google_monitoring_v3_dropped_labels_proto_goTypes = nil
	file_google_monitoring_v3_dropped_labels_proto_depIdxs = nil
}
