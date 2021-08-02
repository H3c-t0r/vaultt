// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: sdk/logical/identity.proto

package logical

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

type Entity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ID is the unique identifier for the entity
	ID string `sentinel:"" protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	// Name is the human-friendly unique identifier for the entity
	Name string `sentinel:"" protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Aliases contains thhe alias mappings for the given entity
	Aliases []*Alias `sentinel:"" protobuf:"bytes,3,rep,name=aliases,proto3" json:"aliases,omitempty"`
	// Metadata represents the custom data tied to this entity
	Metadata map[string]string `sentinel:"" protobuf:"bytes,4,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Disabled is true if the entity is disabled.
	Disabled bool `sentinel:"" protobuf:"varint,5,opt,name=disabled,proto3" json:"disabled,omitempty"`
	// NamespaceID is the identifier of the namespace to which this entity
	// belongs to.
	NamespaceID string `sentinel:"" protobuf:"bytes,6,opt,name=namespace_id,json=namespaceID,proto3" json:"namespace_id,omitempty"`
}

func (x *Entity) Reset() {
	*x = Entity{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sdk_logical_identity_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Entity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entity) ProtoMessage() {}

func (x *Entity) ProtoReflect() protoreflect.Message {
	mi := &file_sdk_logical_identity_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entity.ProtoReflect.Descriptor instead.
func (*Entity) Descriptor() ([]byte, []int) {
	return file_sdk_logical_identity_proto_rawDescGZIP(), []int{0}
}

func (x *Entity) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Entity) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Entity) GetAliases() []*Alias {
	if x != nil {
		return x.Aliases
	}
	return nil
}

func (x *Entity) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *Entity) GetDisabled() bool {
	if x != nil {
		return x.Disabled
	}
	return false
}

func (x *Entity) GetNamespaceID() string {
	if x != nil {
		return x.NamespaceID
	}
	return ""
}

type Alias struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// MountType is the backend mount's type to which this identity belongs
	MountType string `sentinel:"" protobuf:"bytes,1,opt,name=mount_type,json=mountType,proto3" json:"mount_type,omitempty"`
	// MountAccessor is the identifier of the mount entry to which this
	// identity belongs
	MountAccessor string `sentinel:"" protobuf:"bytes,2,opt,name=mount_accessor,json=mountAccessor,proto3" json:"mount_accessor,omitempty"`
	// Name is the identifier of this identity in its authentication source
	Name string `sentinel:"" protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	// Metadata represents the custom data tied to this alias. Fields added
	// to it should have a low rate of change (or no change) because each
	// change incurs a storage write, so quickly-changing fields can have
	// a significant performance impact at scale. See the SDK's
	// "aliasmetadata" package for a helper that eases and standardizes
	// using this safely.
	Metadata map[string]string `sentinel:"" protobuf:"bytes,4,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// ID is the unique identifier for the alias
	ID string `sentinel:"" protobuf:"bytes,5,opt,name=ID,proto3" json:"ID,omitempty"`
	// NamespaceID is the identifier of the namespace to which this alias
	// belongs.
	NamespaceID string `sentinel:"" protobuf:"bytes,6,opt,name=namespace_id,json=namespaceID,proto3" json:"namespace_id,omitempty"`
}

func (x *Alias) Reset() {
	*x = Alias{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sdk_logical_identity_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Alias) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Alias) ProtoMessage() {}

func (x *Alias) ProtoReflect() protoreflect.Message {
	mi := &file_sdk_logical_identity_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Alias.ProtoReflect.Descriptor instead.
func (*Alias) Descriptor() ([]byte, []int) {
	return file_sdk_logical_identity_proto_rawDescGZIP(), []int{1}
}

func (x *Alias) GetMountType() string {
	if x != nil {
		return x.MountType
	}
	return ""
}

func (x *Alias) GetMountAccessor() string {
	if x != nil {
		return x.MountAccessor
	}
	return ""
}

func (x *Alias) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Alias) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *Alias) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Alias) GetNamespaceID() string {
	if x != nil {
		return x.NamespaceID
	}
	return ""
}

type Group struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ID is the unique identifier for the group
	ID string `sentinel:"" protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	// Name is the human-friendly unique identifier for the group
	Name string `sentinel:"" protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Metadata represents the custom data tied to this group
	Metadata map[string]string `sentinel:"" protobuf:"bytes,3,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// NamespaceID is the identifier of the namespace to which this group
	// belongs to.
	NamespaceID string `sentinel:"" protobuf:"bytes,4,opt,name=namespace_id,json=namespaceID,proto3" json:"namespace_id,omitempty"`
}

func (x *Group) Reset() {
	*x = Group{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sdk_logical_identity_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Group) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Group) ProtoMessage() {}

func (x *Group) ProtoReflect() protoreflect.Message {
	mi := &file_sdk_logical_identity_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Group.ProtoReflect.Descriptor instead.
func (*Group) Descriptor() ([]byte, []int) {
	return file_sdk_logical_identity_proto_rawDescGZIP(), []int{2}
}

func (x *Group) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Group) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Group) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *Group) GetNamespaceID() string {
	if x != nil {
		return x.NamespaceID
	}
	return ""
}

var File_sdk_logical_identity_proto protoreflect.FileDescriptor

var file_sdk_logical_identity_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x73, 0x64, 0x6b, 0x2f, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x61, 0x6c, 0x2f, 0x69, 0x64,
	0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6c, 0x6f,
	0x67, 0x69, 0x63, 0x61, 0x6c, 0x22, 0x8d, 0x02, 0x0a, 0x06, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x28, 0x0a, 0x07, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x65, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x61, 0x6c, 0x2e,
	0x41, 0x6c, 0x69, 0x61, 0x73, 0x52, 0x07, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x65, 0x73, 0x12, 0x39,
	0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x1d, 0x2e, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x61, 0x6c, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x69, 0x73,
	0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x64, 0x69, 0x73,
	0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61,
	0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6e, 0x61, 0x6d,
	0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x49, 0x64, 0x1a, 0x3b, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x8b, 0x02, 0x0a, 0x05, 0x41, 0x6c, 0x69, 0x61, 0x73, 0x12,
	0x1d, 0x0a, 0x0a, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x25,
	0x0a, 0x0e, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x41, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x6f, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x38, 0x0a, 0x08, 0x6d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x6c, 0x6f,
	0x67, 0x69, 0x63, 0x61, 0x6c, 0x2e, 0x41, 0x6c, 0x69, 0x61, 0x73, 0x2e, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x49, 0x44, 0x12, 0x21, 0x0a, 0x0c, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6e, 0x61, 0x6d, 0x65, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x49, 0x64, 0x1a, 0x3b, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x22, 0xc5, 0x01, 0x0a, 0x05, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x0e, 0x0a,
	0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x38, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x61, 0x6c, 0x2e, 0x47, 0x72,
	0x6f, 0x75, 0x70, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x21, 0x0a, 0x0c, 0x6e,
	0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x49, 0x64, 0x1a, 0x3b,
	0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x28, 0x5a, 0x26, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x61, 0x73, 0x68, 0x69, 0x63,
	0x6f, 0x72, 0x70, 0x2f, 0x76, 0x61, 0x75, 0x6c, 0x74, 0x2f, 0x73, 0x64, 0x6b, 0x2f, 0x6c, 0x6f,
	0x67, 0x69, 0x63, 0x61, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_sdk_logical_identity_proto_rawDescOnce sync.Once
	file_sdk_logical_identity_proto_rawDescData = file_sdk_logical_identity_proto_rawDesc
)

func file_sdk_logical_identity_proto_rawDescGZIP() []byte {
	file_sdk_logical_identity_proto_rawDescOnce.Do(func() {
		file_sdk_logical_identity_proto_rawDescData = protoimpl.X.CompressGZIP(file_sdk_logical_identity_proto_rawDescData)
	})
	return file_sdk_logical_identity_proto_rawDescData
}

var file_sdk_logical_identity_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_sdk_logical_identity_proto_goTypes = []interface{}{
	(*Entity)(nil), // 0: logical.Entity
	(*Alias)(nil),  // 1: logical.Alias
	(*Group)(nil),  // 2: logical.Group
	nil,            // 3: logical.Entity.MetadataEntry
	nil,            // 4: logical.Alias.MetadataEntry
	nil,            // 5: logical.Group.MetadataEntry
}
var file_sdk_logical_identity_proto_depIDxs = []int32{
	1, // 0: logical.Entity.aliases:type_name -> logical.Alias
	3, // 1: logical.Entity.metadata:type_name -> logical.Entity.MetadataEntry
	4, // 2: logical.Alias.metadata:type_name -> logical.Alias.MetadataEntry
	5, // 3: logical.Group.metadata:type_name -> logical.Group.MetadataEntry
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_sdk_logical_identity_proto_init() }
func file_sdk_logical_identity_proto_init() {
	if File_sdk_logical_identity_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_sdk_logical_identity_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Entity); i {
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
		file_sdk_logical_identity_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Alias); i {
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
		file_sdk_logical_identity_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Group); i {
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
			RawDescriptor: file_sdk_logical_identity_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_sdk_logical_identity_proto_goTypes,
		DependencyIndexes: file_sdk_logical_identity_proto_depIDxs,
		MessageInfos:      file_sdk_logical_identity_proto_msgTypes,
	}.Build()
	File_sdk_logical_identity_proto = out.File
	file_sdk_logical_identity_proto_rawDesc = nil
	file_sdk_logical_identity_proto_goTypes = nil
	file_sdk_logical_identity_proto_depIDxs = nil
}
