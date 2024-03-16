// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.3
// source: transaction/type.proto

package transaction

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

type Type int32

const (
	Type_UNKNOWN        Type = 0
	Type_TASK_ASSIGNED  Type = 1
	Type_TASK_COMPLETED Type = 2
	Type_PAYOUT         Type = 3
)

// Enum value maps for Type.
var (
	Type_name = map[int32]string{
		0: "UNKNOWN",
		1: "TASK_ASSIGNED",
		2: "TASK_COMPLETED",
		3: "PAYOUT",
	}
	Type_value = map[string]int32{
		"UNKNOWN":        0,
		"TASK_ASSIGNED":  1,
		"TASK_COMPLETED": 2,
		"PAYOUT":         3,
	}
)

func (x Type) Enum() *Type {
	p := new(Type)
	*p = x
	return p
}

func (x Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Type) Descriptor() protoreflect.EnumDescriptor {
	return file_transaction_type_proto_enumTypes[0].Descriptor()
}

func (Type) Type() protoreflect.EnumType {
	return &file_transaction_type_proto_enumTypes[0]
}

func (x Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Type.Descriptor instead.
func (Type) EnumDescriptor() ([]byte, []int) {
	return file_transaction_type_proto_rawDescGZIP(), []int{0}
}

var File_transaction_type_proto protoreflect.FileDescriptor

var file_transaction_type_proto_rawDesc = []byte{
	0x0a, 0x16, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x74, 0x79,
	0x70, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2a, 0x46, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a,
	0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x54, 0x41,
	0x53, 0x4b, 0x5f, 0x41, 0x53, 0x53, 0x49, 0x47, 0x4e, 0x45, 0x44, 0x10, 0x01, 0x12, 0x12, 0x0a,
	0x0e, 0x54, 0x41, 0x53, 0x4b, 0x5f, 0x43, 0x4f, 0x4d, 0x50, 0x4c, 0x45, 0x54, 0x45, 0x44, 0x10,
	0x02, 0x12, 0x0a, 0x0a, 0x06, 0x50, 0x41, 0x59, 0x4f, 0x55, 0x54, 0x10, 0x03, 0x42, 0x45, 0x5a,
	0x43, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6c, 0x79, 0x61,
	0x74, 0x6f, 0x73, 0x2f, 0x61, 0x61, 0x63, 0x2d, 0x74, 0x61, 0x73, 0x6b, 0x2d, 0x74, 0x72, 0x61,
	0x63, 0x6b, 0x65, 0x72, 0x2f, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x5f, 0x72, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x72, 0x79, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_transaction_type_proto_rawDescOnce sync.Once
	file_transaction_type_proto_rawDescData = file_transaction_type_proto_rawDesc
)

func file_transaction_type_proto_rawDescGZIP() []byte {
	file_transaction_type_proto_rawDescOnce.Do(func() {
		file_transaction_type_proto_rawDescData = protoimpl.X.CompressGZIP(file_transaction_type_proto_rawDescData)
	})
	return file_transaction_type_proto_rawDescData
}

var file_transaction_type_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_transaction_type_proto_goTypes = []interface{}{
	(Type)(0), // 0: transaction.Type
}
var file_transaction_type_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_transaction_type_proto_init() }
func file_transaction_type_proto_init() {
	if File_transaction_type_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_transaction_type_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_transaction_type_proto_goTypes,
		DependencyIndexes: file_transaction_type_proto_depIdxs,
		EnumInfos:         file_transaction_type_proto_enumTypes,
	}.Build()
	File_transaction_type_proto = out.File
	file_transaction_type_proto_rawDesc = nil
	file_transaction_type_proto_goTypes = nil
	file_transaction_type_proto_depIdxs = nil
}
