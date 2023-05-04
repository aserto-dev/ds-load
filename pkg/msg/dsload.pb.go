// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: dsload.proto

package msg

import (
	v2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	v1 "github.com/aserto-dev/go-grpc/aserto/common/info/v1"
	status "google.golang.org/genproto/googleapis/rpc/status"
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

type ConfigElementType int32

const (
	ConfigElementType_CONFIG_ELEMENT_TYPE_UNKNOWN ConfigElementType = 0 // Unknown configuration element type
	ConfigElementType_CONFIG_ELEMENT_TYPE_STRING  ConfigElementType = 1 // String configuration element type
	ConfigElementType_CONFIG_ELEMENT_TYPE_INTEGER ConfigElementType = 2 // Integer configuration element type
	ConfigElementType_CONFIG_ELEMENT_TYPE_BOOLEAN ConfigElementType = 3 // Boolean configuration element type
)

// Enum value maps for ConfigElementType.
var (
	ConfigElementType_name = map[int32]string{
		0: "CONFIG_ELEMENT_TYPE_UNKNOWN",
		1: "CONFIG_ELEMENT_TYPE_STRING",
		2: "CONFIG_ELEMENT_TYPE_INTEGER",
		3: "CONFIG_ELEMENT_TYPE_BOOLEAN",
	}
	ConfigElementType_value = map[string]int32{
		"CONFIG_ELEMENT_TYPE_UNKNOWN": 0,
		"CONFIG_ELEMENT_TYPE_STRING":  1,
		"CONFIG_ELEMENT_TYPE_INTEGER": 2,
		"CONFIG_ELEMENT_TYPE_BOOLEAN": 3,
	}
)

func (x ConfigElementType) Enum() *ConfigElementType {
	p := new(ConfigElementType)
	*p = x
	return p
}

func (x ConfigElementType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ConfigElementType) Descriptor() protoreflect.EnumDescriptor {
	return file_dsload_proto_enumTypes[0].Descriptor()
}

func (ConfigElementType) Type() protoreflect.EnumType {
	return &file_dsload_proto_enumTypes[0]
}

func (x ConfigElementType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ConfigElementType.Descriptor instead.
func (ConfigElementType) EnumDescriptor() ([]byte, []int) {
	return file_dsload_proto_rawDescGZIP(), []int{0}
}

type ConfigElement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type        ConfigElementType `protobuf:"varint,1,opt,name=type,proto3,enum=aserto.dsload.ConfigElementType" json:"type,omitempty"`
	Name        string            `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description string            `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Usage       string            `protobuf:"bytes,4,opt,name=usage,proto3" json:"usage,omitempty"`
}

func (x *ConfigElement) Reset() {
	*x = ConfigElement{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dsload_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigElement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigElement) ProtoMessage() {}

func (x *ConfigElement) ProtoReflect() protoreflect.Message {
	mi := &file_dsload_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigElement.ProtoReflect.Descriptor instead.
func (*ConfigElement) Descriptor() ([]byte, []int) {
	return file_dsload_proto_rawDescGZIP(), []int{0}
}

func (x *ConfigElement) GetType() ConfigElementType {
	if x != nil {
		return x.Type
	}
	return ConfigElementType_CONFIG_ELEMENT_TYPE_UNKNOWN
}

func (x *ConfigElement) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ConfigElement) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *ConfigElement) GetUsage() string {
	if x != nil {
		return x.Usage
	}
	return ""
}

type Info struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Build       *v1.BuildInfo    `protobuf:"bytes,1,opt,name=build,proto3" json:"build,omitempty"`
	Description string           `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Configs     []*ConfigElement `protobuf:"bytes,3,rep,name=configs,proto3" json:"configs,omitempty"`
}

func (x *Info) Reset() {
	*x = Info{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dsload_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Info) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Info) ProtoMessage() {}

func (x *Info) ProtoReflect() protoreflect.Message {
	mi := &file_dsload_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Info.ProtoReflect.Descriptor instead.
func (*Info) Descriptor() ([]byte, []int) {
	return file_dsload_proto_rawDescGZIP(), []int{1}
}

func (x *Info) GetBuild() *v1.BuildInfo {
	if x != nil {
		return x.Build
	}
	return nil
}

func (x *Info) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Info) GetConfigs() []*ConfigElement {
	if x != nil {
		return x.Configs
	}
	return nil
}

type Batch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Data:
	//
	//	*Batch_Begin
	//	*Batch_End
	Data isBatch_Data `protobuf_oneof:"data"`
}

func (x *Batch) Reset() {
	*x = Batch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dsload_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Batch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Batch) ProtoMessage() {}

func (x *Batch) ProtoReflect() protoreflect.Message {
	mi := &file_dsload_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Batch.ProtoReflect.Descriptor instead.
func (*Batch) Descriptor() ([]byte, []int) {
	return file_dsload_proto_rawDescGZIP(), []int{2}
}

func (m *Batch) GetData() isBatch_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *Batch) GetBegin() bool {
	if x, ok := x.GetData().(*Batch_Begin); ok {
		return x.Begin
	}
	return false
}

func (x *Batch) GetEnd() bool {
	if x, ok := x.GetData().(*Batch_End); ok {
		return x.End
	}
	return false
}

type isBatch_Data interface {
	isBatch_Data()
}

type Batch_Begin struct {
	Begin bool `protobuf:"varint,1,opt,name=begin,proto3,oneof"`
}

type Batch_End struct {
	End bool `protobuf:"varint,2,opt,name=end,proto3,oneof"`
}

func (*Batch_Begin) isBatch_Data() {}

func (*Batch_End) isBatch_Data() {}

type PluginMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Data:
	//
	//	*PluginMessage_Object
	//	*PluginMessage_Relation
	//	*PluginMessage_Batch
	//	*PluginMessage_Error
	Data isPluginMessage_Data `protobuf_oneof:"data"`
}

func (x *PluginMessage) Reset() {
	*x = PluginMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dsload_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PluginMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PluginMessage) ProtoMessage() {}

func (x *PluginMessage) ProtoReflect() protoreflect.Message {
	mi := &file_dsload_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PluginMessage.ProtoReflect.Descriptor instead.
func (*PluginMessage) Descriptor() ([]byte, []int) {
	return file_dsload_proto_rawDescGZIP(), []int{3}
}

func (m *PluginMessage) GetData() isPluginMessage_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *PluginMessage) GetObject() *v2.Object {
	if x, ok := x.GetData().(*PluginMessage_Object); ok {
		return x.Object
	}
	return nil
}

func (x *PluginMessage) GetRelation() *v2.Relation {
	if x, ok := x.GetData().(*PluginMessage_Relation); ok {
		return x.Relation
	}
	return nil
}

func (x *PluginMessage) GetBatch() *Batch {
	if x, ok := x.GetData().(*PluginMessage_Batch); ok {
		return x.Batch
	}
	return nil
}

func (x *PluginMessage) GetError() *status.Status {
	if x, ok := x.GetData().(*PluginMessage_Error); ok {
		return x.Error
	}
	return nil
}

type isPluginMessage_Data interface {
	isPluginMessage_Data()
}

type PluginMessage_Object struct {
	Object *v2.Object `protobuf:"bytes,1,opt,name=object,proto3,oneof"`
}

type PluginMessage_Relation struct {
	Relation *v2.Relation `protobuf:"bytes,2,opt,name=relation,proto3,oneof"`
}

type PluginMessage_Batch struct {
	Batch *Batch `protobuf:"bytes,3,opt,name=batch,proto3,oneof"`
}

type PluginMessage_Error struct {
	Error *status.Status `protobuf:"bytes,4,opt,name=error,proto3,oneof"`
}

func (*PluginMessage_Object) isPluginMessage_Data() {}

func (*PluginMessage_Relation) isPluginMessage_Data() {}

func (*PluginMessage_Batch) isPluginMessage_Data() {}

func (*PluginMessage_Error) isPluginMessage_Data() {}

var File_dsload_proto protoreflect.FileDescriptor

var file_dsload_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x64, 0x73, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d,
	0x61, 0x73, 0x65, 0x72, 0x74, 0x6f, 0x2e, 0x64, 0x73, 0x6c, 0x6f, 0x61, 0x64, 0x1a, 0x17, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x61, 0x73, 0x65, 0x72, 0x74, 0x6f, 0x2f, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x69, 0x6e, 0x66, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x6e,
	0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x27, 0x61, 0x73, 0x65, 0x72, 0x74, 0x6f,
	0x2f, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2f, 0x76, 0x32, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x91, 0x01, 0x0a, 0x0d, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x45, 0x6c, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x12, 0x34, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x20, 0x2e, 0x61, 0x73, 0x65, 0x72, 0x74, 0x6f, 0x2e, 0x64, 0x73, 0x6c, 0x6f, 0x61,
	0x64, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x14, 0x0a, 0x05, 0x75, 0x73, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x75, 0x73, 0x61, 0x67, 0x65, 0x22, 0x98, 0x01, 0x0a, 0x04, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x36,
	0x0a, 0x05, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e,
	0x61, 0x73, 0x65, 0x72, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x69, 0x6e,
	0x66, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x05, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x36, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x61, 0x73, 0x65, 0x72,
	0x74, 0x6f, 0x2e, 0x64, 0x73, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73,
	0x22, 0x3b, 0x0a, 0x05, 0x42, 0x61, 0x74, 0x63, 0x68, 0x12, 0x16, 0x0a, 0x05, 0x62, 0x65, 0x67,
	0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x05, 0x62, 0x65, 0x67, 0x69,
	0x6e, 0x12, 0x12, 0x0a, 0x03, 0x65, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00,
	0x52, 0x03, 0x65, 0x6e, 0x64, 0x42, 0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0xf3, 0x01,
	0x0a, 0x0d, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x3c, 0x0a, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x22, 0x2e, 0x61, 0x73, 0x65, 0x72, 0x74, 0x6f, 0x2e, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x32, 0x2e, 0x4f, 0x62, 0x6a,
	0x65, 0x63, 0x74, 0x48, 0x00, 0x52, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x42, 0x0a,
	0x08, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x24, 0x2e, 0x61, 0x73, 0x65, 0x72, 0x74, 0x6f, 0x2e, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x32, 0x2e, 0x52, 0x65, 0x6c,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x08, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x2c, 0x0a, 0x05, 0x62, 0x61, 0x74, 0x63, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x61, 0x73, 0x65, 0x72, 0x74, 0x6f, 0x2e, 0x64, 0x73, 0x6c, 0x6f, 0x61, 0x64,
	0x2e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x48, 0x00, 0x52, 0x05, 0x62, 0x61, 0x74, 0x63, 0x68, 0x12,
	0x2a, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x48, 0x00, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x42, 0x06, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x2a, 0x96, 0x01, 0x0a, 0x11, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x45, 0x6c,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1f, 0x0a, 0x1b, 0x43, 0x4f, 0x4e,
	0x46, 0x49, 0x47, 0x5f, 0x45, 0x4c, 0x45, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x1e, 0x0a, 0x1a, 0x43, 0x4f,
	0x4e, 0x46, 0x49, 0x47, 0x5f, 0x45, 0x4c, 0x45, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x53, 0x54, 0x52, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x1f, 0x0a, 0x1b, 0x43, 0x4f,
	0x4e, 0x46, 0x49, 0x47, 0x5f, 0x45, 0x4c, 0x45, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x47, 0x45, 0x52, 0x10, 0x02, 0x12, 0x1f, 0x0a, 0x1b, 0x43,
	0x4f, 0x4e, 0x46, 0x49, 0x47, 0x5f, 0x45, 0x4c, 0x45, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x42, 0x4f, 0x4f, 0x4c, 0x45, 0x41, 0x4e, 0x10, 0x03, 0x42, 0x2b, 0x5a, 0x29,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x73, 0x65, 0x72, 0x74,
	0x6f, 0x2d, 0x64, 0x65, 0x76, 0x2f, 0x64, 0x73, 0x2d, 0x6c, 0x6f, 0x61, 0x64, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x6d, 0x73, 0x67, 0x3b, 0x6d, 0x73, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_dsload_proto_rawDescOnce sync.Once
	file_dsload_proto_rawDescData = file_dsload_proto_rawDesc
)

func file_dsload_proto_rawDescGZIP() []byte {
	file_dsload_proto_rawDescOnce.Do(func() {
		file_dsload_proto_rawDescData = protoimpl.X.CompressGZIP(file_dsload_proto_rawDescData)
	})
	return file_dsload_proto_rawDescData
}

var file_dsload_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_dsload_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_dsload_proto_goTypes = []interface{}{
	(ConfigElementType)(0), // 0: aserto.dsload.ConfigElementType
	(*ConfigElement)(nil),  // 1: aserto.dsload.ConfigElement
	(*Info)(nil),           // 2: aserto.dsload.Info
	(*Batch)(nil),          // 3: aserto.dsload.Batch
	(*PluginMessage)(nil),  // 4: aserto.dsload.PluginMessage
	(*v1.BuildInfo)(nil),   // 5: aserto.common.info.v1.BuildInfo
	(*v2.Object)(nil),      // 6: aserto.directory.common.v2.Object
	(*v2.Relation)(nil),    // 7: aserto.directory.common.v2.Relation
	(*status.Status)(nil),  // 8: google.rpc.Status
}
var file_dsload_proto_depIdxs = []int32{
	0, // 0: aserto.dsload.ConfigElement.type:type_name -> aserto.dsload.ConfigElementType
	5, // 1: aserto.dsload.Info.build:type_name -> aserto.common.info.v1.BuildInfo
	1, // 2: aserto.dsload.Info.configs:type_name -> aserto.dsload.ConfigElement
	6, // 3: aserto.dsload.PluginMessage.object:type_name -> aserto.directory.common.v2.Object
	7, // 4: aserto.dsload.PluginMessage.relation:type_name -> aserto.directory.common.v2.Relation
	3, // 5: aserto.dsload.PluginMessage.batch:type_name -> aserto.dsload.Batch
	8, // 6: aserto.dsload.PluginMessage.error:type_name -> google.rpc.Status
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_dsload_proto_init() }
func file_dsload_proto_init() {
	if File_dsload_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_dsload_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfigElement); i {
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
		file_dsload_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Info); i {
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
		file_dsload_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Batch); i {
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
		file_dsload_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PluginMessage); i {
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
	file_dsload_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*Batch_Begin)(nil),
		(*Batch_End)(nil),
	}
	file_dsload_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*PluginMessage_Object)(nil),
		(*PluginMessage_Relation)(nil),
		(*PluginMessage_Batch)(nil),
		(*PluginMessage_Error)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_dsload_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_dsload_proto_goTypes,
		DependencyIndexes: file_dsload_proto_depIdxs,
		EnumInfos:         file_dsload_proto_enumTypes,
		MessageInfos:      file_dsload_proto_msgTypes,
	}.Build()
	File_dsload_proto = out.File
	file_dsload_proto_rawDesc = nil
	file_dsload_proto_goTypes = nil
	file_dsload_proto_depIdxs = nil
}
