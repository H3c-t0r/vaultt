// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/api/resource.proto

package annotations

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
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

// A description of the historical or future-looking state of the
// resource pattern.
type ResourceDescriptor_History int32

const (
	// The "unset" value.
	ResourceDescriptor_HISTORY_UNSPECIFIED ResourceDescriptor_History = 0
	// The resource originally had one pattern and launched as such, and
	// additional patterns were added later.
	ResourceDescriptor_ORIGINALLY_SINGLE_PATTERN ResourceDescriptor_History = 1
	// The resource has one pattern, but the API owner expects to add more
	// later. (This is the inverse of ORIGINALLY_SINGLE_PATTERN, and prevents
	// that from being necessary once there are multiple patterns.)
	ResourceDescriptor_FUTURE_MULTI_PATTERN ResourceDescriptor_History = 2
)

var ResourceDescriptor_History_name = map[int32]string{
	0: "HISTORY_UNSPECIFIED",
	1: "ORIGINALLY_SINGLE_PATTERN",
	2: "FUTURE_MULTI_PATTERN",
}

var ResourceDescriptor_History_value = map[string]int32{
	"HISTORY_UNSPECIFIED":       0,
	"ORIGINALLY_SINGLE_PATTERN": 1,
	"FUTURE_MULTI_PATTERN":      2,
}

func (x ResourceDescriptor_History) String() string {
	return proto.EnumName(ResourceDescriptor_History_name, int32(x))
}

func (ResourceDescriptor_History) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_465e9122405d1bb5, []int{0, 0}
}

// A simple descriptor of a resource type.
//
// ResourceDescriptor annotates a resource message (either by means of a
// protobuf annotation or use in the service config), and associates the
// resource's schema, the resource type, and the pattern of the resource name.
//
// Example:
//
//   message Topic {
//     // Indicates this message defines a resource schema.
//     // Declares the resource type in the format of {service}/{kind}.
//     // For Kubernetes resources, the format is {api group}/{kind}.
//     option (google.api.resource) = {
//       type: "pubsub.googleapis.com/Topic"
//       pattern: "projects/{project}/topics/{topic}"
//     };
//   }
//
// Sometimes, resources have multiple patterns, typically because they can
// live under multiple parents.
//
// Example:
//
//   message LogEntry {
//     option (google.api.resource) = {
//       type: "logging.googleapis.com/LogEntry"
//       pattern: "projects/{project}/logs/{log}"
//       pattern: "organizations/{organization}/logs/{log}"
//       pattern: "folders/{folder}/logs/{log}"
//       pattern: "billingAccounts/{billing_account}/logs/{log}"
//     };
//   }
type ResourceDescriptor struct {
	// The full name of the resource type. It must be in the format of
	// {service_name}/{resource_type_kind}. The resource type names are
	// singular and do not contain version numbers.
	//
	// For example: `storage.googleapis.com/Bucket`
	//
	// The value of the resource_type_kind must follow the regular expression
	// /[A-Z][a-zA-Z0-9]+/. It must start with upper case character and
	// recommended to use PascalCase (UpperCamelCase). The maximum number of
	// characters allowed for the resource_type_kind is 100.
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	// Required. The valid pattern or patterns for this resource's names.
	//
	// Examples:
	//   - "projects/{project}/topics/{topic}"
	//   - "projects/{project}/knowledgeBases/{knowledge_base}"
	//
	// The components in braces correspond to the IDs for each resource in the
	// hierarchy. It is expected that, if multiple patterns are provided,
	// the same component name (e.g. "project") refers to IDs of the same
	// type of resource.
	Pattern []string `protobuf:"bytes,2,rep,name=pattern,proto3" json:"pattern,omitempty"`
	// Optional. The field on the resource that designates the resource name
	// field. If omitted, this is assumed to be "name".
	NameField string `protobuf:"bytes,3,opt,name=name_field,json=nameField,proto3" json:"name_field,omitempty"`
	// Optional. The historical or future-looking state of the resource pattern.
	//
	// Example:
	//   // The InspectTemplate message originally only supported resource
	//   // names with organization, and project was added later.
	//   message InspectTemplate {
	//     option (google.api.resource) = {
	//       type: "dlp.googleapis.com/InspectTemplate"
	//       pattern: "organizations/{organization}/inspectTemplates/{inspect_template}"
	//       pattern: "projects/{project}/inspectTemplates/{inspect_template}"
	//       history: ORIGINALLY_SINGLE_PATTERN
	//     };
	//   }
	History              ResourceDescriptor_History `protobuf:"varint,4,opt,name=history,proto3,enum=google.api.ResourceDescriptor_History" json:"history,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *ResourceDescriptor) Reset()         { *m = ResourceDescriptor{} }
func (m *ResourceDescriptor) String() string { return proto.CompactTextString(m) }
func (*ResourceDescriptor) ProtoMessage()    {}
func (*ResourceDescriptor) Descriptor() ([]byte, []int) {
	return fileDescriptor_465e9122405d1bb5, []int{0}
}

func (m *ResourceDescriptor) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResourceDescriptor.Unmarshal(m, b)
}
func (m *ResourceDescriptor) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResourceDescriptor.Marshal(b, m, deterministic)
}
func (m *ResourceDescriptor) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResourceDescriptor.Merge(m, src)
}
func (m *ResourceDescriptor) XXX_Size() int {
	return xxx_messageInfo_ResourceDescriptor.Size(m)
}
func (m *ResourceDescriptor) XXX_DiscardUnknown() {
	xxx_messageInfo_ResourceDescriptor.DiscardUnknown(m)
}

var xxx_messageInfo_ResourceDescriptor proto.InternalMessageInfo

func (m *ResourceDescriptor) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *ResourceDescriptor) GetPattern() []string {
	if m != nil {
		return m.Pattern
	}
	return nil
}

func (m *ResourceDescriptor) GetNameField() string {
	if m != nil {
		return m.NameField
	}
	return ""
}

func (m *ResourceDescriptor) GetHistory() ResourceDescriptor_History {
	if m != nil {
		return m.History
	}
	return ResourceDescriptor_HISTORY_UNSPECIFIED
}

// An annotation designating that this field is a reference to a resource
// defined by another message.
type ResourceReference struct {
	// The unified resource type name of the type that this field references.
	// Marks this as a field referring to a resource in another message.
	//
	// Example:
	//
	//   message Subscription {
	//     string topic = 2 [(google.api.resource_reference) = {
	//       type = "pubsub.googleapis.com/Topic"
	//     }];
	//   }
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	// The fully-qualified message name of a child of the type that this field
	// references.
	//
	// This is useful for `parent` fields where a resource has more than one
	// possible type of parent.
	//
	// Example:
	//
	//   message ListLogEntriesRequest {
	//     string parent = 1 [(google.api.resource_reference) = {
	//       child_type: "logging.googleapis.com/LogEntry"
	//     };
	//   }
	//
	// If the referenced message is in the same proto package, the service name
	// may be omitted:
	//
	//   message ListLogEntriesRequest {
	//     string parent = 1
	//       [(google.api.resource_reference).child_type = "LogEntry"];
	//   }
	ChildType            string   `protobuf:"bytes,2,opt,name=child_type,json=childType,proto3" json:"child_type,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResourceReference) Reset()         { *m = ResourceReference{} }
func (m *ResourceReference) String() string { return proto.CompactTextString(m) }
func (*ResourceReference) ProtoMessage()    {}
func (*ResourceReference) Descriptor() ([]byte, []int) {
	return fileDescriptor_465e9122405d1bb5, []int{1}
}

func (m *ResourceReference) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResourceReference.Unmarshal(m, b)
}
func (m *ResourceReference) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResourceReference.Marshal(b, m, deterministic)
}
func (m *ResourceReference) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResourceReference.Merge(m, src)
}
func (m *ResourceReference) XXX_Size() int {
	return xxx_messageInfo_ResourceReference.Size(m)
}
func (m *ResourceReference) XXX_DiscardUnknown() {
	xxx_messageInfo_ResourceReference.DiscardUnknown(m)
}

var xxx_messageInfo_ResourceReference proto.InternalMessageInfo

func (m *ResourceReference) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *ResourceReference) GetChildType() string {
	if m != nil {
		return m.ChildType
	}
	return ""
}

var E_ResourceReference = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*ResourceReference)(nil),
	Field:         1055,
	Name:          "google.api.resource_reference",
	Tag:           "bytes,1055,opt,name=resource_reference",
	Filename:      "google/api/resource.proto",
}

var E_Resource = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*ResourceDescriptor)(nil),
	Field:         1053,
	Name:          "google.api.resource",
	Tag:           "bytes,1053,opt,name=resource",
	Filename:      "google/api/resource.proto",
}

func init() {
	proto.RegisterEnum("google.api.ResourceDescriptor_History", ResourceDescriptor_History_name, ResourceDescriptor_History_value)
	proto.RegisterType((*ResourceDescriptor)(nil), "google.api.ResourceDescriptor")
	proto.RegisterType((*ResourceReference)(nil), "google.api.ResourceReference")
	proto.RegisterExtension(E_ResourceReference)
	proto.RegisterExtension(E_Resource)
}

func init() { proto.RegisterFile("google/api/resource.proto", fileDescriptor_465e9122405d1bb5) }

var fileDescriptor_465e9122405d1bb5 = []byte{
	// 430 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x52, 0x41, 0x6f, 0xd3, 0x30,
	0x18, 0x25, 0x59, 0x45, 0xd7, 0x0f, 0x31, 0x6d, 0x06, 0x89, 0x0c, 0x29, 0x10, 0xf5, 0x80, 0x7a,
	0x4a, 0xa4, 0x71, 0x1b, 0x17, 0x3a, 0x96, 0x76, 0x91, 0xba, 0x36, 0x72, 0xd3, 0xc3, 0x00, 0x29,
	0xf2, 0xd2, 0xaf, 0x59, 0xa4, 0xcc, 0xb6, 0x9c, 0xec, 0xd0, 0x1b, 0x7f, 0x04, 0x21, 0xf1, 0x2b,
	0x39, 0xa2, 0x3a, 0x71, 0x98, 0xd8, 0xb4, 0x9b, 0xf3, 0xde, 0xfb, 0xbe, 0xf7, 0xfc, 0x1c, 0x38,
	0xce, 0x85, 0xc8, 0x4b, 0x0c, 0x98, 0x2c, 0x02, 0x85, 0x95, 0xb8, 0x53, 0x19, 0xfa, 0x52, 0x89,
	0x5a, 0x10, 0x68, 0x28, 0x9f, 0xc9, 0xe2, 0xad, 0xd7, 0xca, 0x34, 0x73, 0x7d, 0xb7, 0x09, 0xd6,
	0x58, 0x65, 0xaa, 0x90, 0xb5, 0x50, 0x8d, 0x7a, 0xf8, 0xc3, 0x06, 0x42, 0xdb, 0x05, 0xe7, 0x1d,
	0x49, 0x08, 0xf4, 0xea, 0xad, 0x44, 0xc7, 0xf2, 0xac, 0xd1, 0x80, 0xea, 0x33, 0x71, 0xa0, 0x2f,
	0x59, 0x5d, 0xa3, 0xe2, 0x8e, 0xed, 0xed, 0x8d, 0x06, 0xd4, 0x7c, 0x12, 0x17, 0x80, 0xb3, 0x5b,
	0x4c, 0x37, 0x05, 0x96, 0x6b, 0x67, 0x4f, 0xcf, 0x0c, 0x76, 0xc8, 0x64, 0x07, 0x90, 0xcf, 0xd0,
	0xbf, 0x29, 0xaa, 0x5a, 0xa8, 0xad, 0xd3, 0xf3, 0xac, 0xd1, 0xc1, 0xc9, 0x07, 0xff, 0x5f, 0x46,
	0xff, 0xa1, 0xbb, 0x7f, 0xd1, 0xa8, 0xa9, 0x19, 0x1b, 0x7e, 0x83, 0x7e, 0x8b, 0x91, 0x37, 0xf0,
	0xea, 0x22, 0x5a, 0x26, 0x0b, 0x7a, 0x95, 0xae, 0xe6, 0xcb, 0x38, 0xfc, 0x12, 0x4d, 0xa2, 0xf0,
	0xfc, 0xf0, 0x19, 0x71, 0xe1, 0x78, 0x41, 0xa3, 0x69, 0x34, 0x1f, 0xcf, 0x66, 0x57, 0xe9, 0x32,
	0x9a, 0x4f, 0x67, 0x61, 0x1a, 0x8f, 0x93, 0x24, 0xa4, 0xf3, 0x43, 0x8b, 0x38, 0xf0, 0x7a, 0xb2,
	0x4a, 0x56, 0x34, 0x4c, 0x2f, 0x57, 0xb3, 0x24, 0xea, 0x18, 0x7b, 0x38, 0x81, 0x23, 0x93, 0x81,
	0xe2, 0x06, 0x15, 0xf2, 0x0c, 0x1f, 0x2d, 0xc0, 0x05, 0xc8, 0x6e, 0x8a, 0x72, 0x9d, 0x6a, 0xc6,
	0x6e, 0xae, 0xa9, 0x91, 0x64, 0x2b, 0xf1, 0xb4, 0x04, 0x62, 0x9e, 0x22, 0x55, 0xdd, 0x22, 0xd7,
	0xdc, 0xd5, 0xbc, 0x81, 0xaf, 0x4b, 0x59, 0xc8, 0xba, 0x10, 0xbc, 0x72, 0x7e, 0xed, 0x7b, 0xd6,
	0xe8, 0xc5, 0x89, 0xfb, 0x58, 0x23, 0x5d, 0x1a, 0x7a, 0xa4, 0xfe, 0x87, 0x4e, 0xbf, 0xc3, 0xbe,
	0x01, 0xc9, 0xfb, 0x07, 0x1e, 0x97, 0x58, 0x55, 0x2c, 0x47, 0xe3, 0xf2, 0xb3, 0x71, 0x79, 0xf7,
	0x74, 0xef, 0xb4, 0xdb, 0x78, 0xc6, 0xe1, 0x20, 0x13, 0xb7, 0xf7, 0xe4, 0x67, 0x2f, 0x8d, 0x3e,
	0xde, 0x79, 0xc4, 0xd6, 0xd7, 0x71, 0x4b, 0xe6, 0xa2, 0x64, 0x3c, 0xf7, 0x85, 0xca, 0x83, 0x1c,
	0xb9, 0x4e, 0x10, 0x34, 0x14, 0x93, 0x45, 0xa5, 0xff, 0x50, 0xc6, 0xb9, 0xa8, 0x99, 0x8e, 0xf2,
	0xe9, 0xde, 0xf9, 0x8f, 0x65, 0xfd, 0xb6, 0x7b, 0xd3, 0x71, 0x1c, 0x5d, 0x3f, 0xd7, 0x73, 0x1f,
	0xff, 0x06, 0x00, 0x00, 0xff, 0xff, 0xb5, 0x1e, 0x07, 0x80, 0xd8, 0x02, 0x00, 0x00,
}
