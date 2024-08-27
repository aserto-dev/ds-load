// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: dsload.proto

package msg

import (
	v2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	v3 "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	v31 "github.com/aserto-dev/go-directory/aserto/directory/importer/v3"
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

type TransformV2 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Objects   []*v2.Object   `protobuf:"bytes,1,rep,name=objects,proto3" json:"objects,omitempty"`
	Relations []*v2.Relation `protobuf:"bytes,2,rep,name=relations,proto3" json:"relations,omitempty"`
}

func (x *TransformV2) Reset() {
	*x = TransformV2{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dsload_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransformV2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransformV2) ProtoMessage() {}

func (x *TransformV2) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use TransformV2.ProtoReflect.Descriptor instead.
func (*TransformV2) Descriptor() ([]byte, []int) {
	return file_dsload_proto_rawDescGZIP(), []int{0}
}

func (x *TransformV2) GetObjects() []*v2.Object {
	if x != nil {
		return x.Objects
	}
	return nil
}

func (x *TransformV2) GetRelations() []*v2.Relation {
	if x != nil {
		return x.Relations
	}
	return nil
}

type Transform struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Objects   []*v3.Object   `protobuf:"bytes,1,rep,name=objects,proto3" json:"objects,omitempty"`
	Relations []*v3.Relation `protobuf:"bytes,2,rep,name=relations,proto3" json:"relations,omitempty"`
	OpCode    v31.Opcode     `protobuf:"varint,3,opt,name=op_code,json=opCode,proto3,enum=aserto.directory.importer.v3.Opcode" json:"op_code,omitempty"`
}

func (x *Transform) Reset() {
	*x = Transform{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dsload_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transform) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transform) ProtoMessage() {}

func (x *Transform) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Transform.ProtoReflect.Descriptor instead.
func (*Transform) Descriptor() ([]byte, []int) {
	return file_dsload_proto_rawDescGZIP(), []int{1}
}

func (x *Transform) GetObjects() []*v3.Object {
	if x != nil {
		return x.Objects
	}
	return nil
}

func (x *Transform) GetRelations() []*v3.Relation {
	if x != nil {
		return x.Relations
	}
	return nil
}

func (x *Transform) GetOpCode() v31.Opcode {
	if x != nil {
		return x.OpCode
	}
	return v31.Opcode(0)
}

var File_dsload_proto protoreflect.FileDescriptor

var file_dsload_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x64, 0x73, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d,
	0x61, 0x73, 0x65, 0x72, 0x74, 0x6f, 0x2e, 0x64, 0x73, 0x6c, 0x6f, 0x61, 0x64, 0x1a, 0x27, 0x61,
	0x73, 0x65, 0x72, 0x74, 0x6f, 0x2f, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x2f,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x32, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x27, 0x61, 0x73, 0x65, 0x72, 0x74, 0x6f, 0x2f, 0x64,
	0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f,
	0x76, 0x33, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x2b, 0x61, 0x73, 0x65, 0x72, 0x74, 0x6f, 0x2f, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x79, 0x2f, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2f, 0x76, 0x33, 0x2f, 0x69, 0x6d,
	0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8f, 0x01, 0x0a,
	0x0b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x56, 0x32, 0x12, 0x3c, 0x0a, 0x07,
	0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e,
	0x61, 0x73, 0x65, 0x72, 0x74, 0x6f, 0x2e, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x32, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x52, 0x07, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x12, 0x42, 0x0a, 0x09, 0x72, 0x65,
	0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e,
	0x61, 0x73, 0x65, 0x72, 0x74, 0x6f, 0x2e, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x32, 0x2e, 0x52, 0x65, 0x6c, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x09, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0xcc,
	0x01, 0x0a, 0x09, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x12, 0x3c, 0x0a, 0x07,
	0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e,
	0x61, 0x73, 0x65, 0x72, 0x74, 0x6f, 0x2e, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x33, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x52, 0x07, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x12, 0x42, 0x0a, 0x09, 0x72, 0x65,
	0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e,
	0x61, 0x73, 0x65, 0x72, 0x74, 0x6f, 0x2e, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x33, 0x2e, 0x52, 0x65, 0x6c, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x09, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x3d,
	0x0a, 0x07, 0x6f, 0x70, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x24, 0x2e, 0x61, 0x73, 0x65, 0x72, 0x74, 0x6f, 0x2e, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x79, 0x2e, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x76, 0x33, 0x2e, 0x4f,
	0x70, 0x63, 0x6f, 0x64, 0x65, 0x52, 0x06, 0x6f, 0x70, 0x43, 0x6f, 0x64, 0x65, 0x42, 0x32, 0x5a,
	0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x73, 0x65, 0x72,
	0x74, 0x6f, 0x2d, 0x64, 0x65, 0x76, 0x2f, 0x64, 0x73, 0x2d, 0x6c, 0x6f, 0x61, 0x64, 0x2f, 0x73,
	0x64, 0x6b, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x6d, 0x73, 0x67, 0x3b, 0x6d, 0x73,
	0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_dsload_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_dsload_proto_goTypes = []any{
	(*TransformV2)(nil), // 0: aserto.dsload.TransformV2
	(*Transform)(nil),   // 1: aserto.dsload.Transform
	(*v2.Object)(nil),   // 2: aserto.directory.common.v2.Object
	(*v2.Relation)(nil), // 3: aserto.directory.common.v2.Relation
	(*v3.Object)(nil),   // 4: aserto.directory.common.v3.Object
	(*v3.Relation)(nil), // 5: aserto.directory.common.v3.Relation
	(v31.Opcode)(0),     // 6: aserto.directory.importer.v3.Opcode
}
var file_dsload_proto_depIdxs = []int32{
	2, // 0: aserto.dsload.TransformV2.objects:type_name -> aserto.directory.common.v2.Object
	3, // 1: aserto.dsload.TransformV2.relations:type_name -> aserto.directory.common.v2.Relation
	4, // 2: aserto.dsload.Transform.objects:type_name -> aserto.directory.common.v3.Object
	5, // 3: aserto.dsload.Transform.relations:type_name -> aserto.directory.common.v3.Relation
	6, // 4: aserto.dsload.Transform.op_code:type_name -> aserto.directory.importer.v3.Opcode
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_dsload_proto_init() }
func file_dsload_proto_init() {
	if File_dsload_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_dsload_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*TransformV2); i {
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
		file_dsload_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Transform); i {
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
			RawDescriptor: file_dsload_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_dsload_proto_goTypes,
		DependencyIndexes: file_dsload_proto_depIdxs,
		MessageInfos:      file_dsload_proto_msgTypes,
	}.Build()
	File_dsload_proto = out.File
	file_dsload_proto_rawDesc = nil
	file_dsload_proto_goTypes = nil
	file_dsload_proto_depIdxs = nil
}
