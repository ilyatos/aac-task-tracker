// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.3
// source: tracker.proto

package tracker

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type CreateTaskRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Description string `protobuf:"bytes,1,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *CreateTaskRequest) Reset() {
	*x = CreateTaskRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tracker_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTaskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTaskRequest) ProtoMessage() {}

func (x *CreateTaskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tracker_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTaskRequest.ProtoReflect.Descriptor instead.
func (*CreateTaskRequest) Descriptor() ([]byte, []int) {
	return file_tracker_proto_rawDescGZIP(), []int{0}
}

func (x *CreateTaskRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

type CreateTaskResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CreateTaskResponse) Reset() {
	*x = CreateTaskResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tracker_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTaskResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTaskResponse) ProtoMessage() {}

func (x *CreateTaskResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tracker_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTaskResponse.ProtoReflect.Descriptor instead.
func (*CreateTaskResponse) Descriptor() ([]byte, []int) {
	return file_tracker_proto_rawDescGZIP(), []int{1}
}

type ReassignTasksRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ReassignTasksRequest) Reset() {
	*x = ReassignTasksRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tracker_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReassignTasksRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReassignTasksRequest) ProtoMessage() {}

func (x *ReassignTasksRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tracker_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReassignTasksRequest.ProtoReflect.Descriptor instead.
func (*ReassignTasksRequest) Descriptor() ([]byte, []int) {
	return file_tracker_proto_rawDescGZIP(), []int{2}
}

type ReassignTasksResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ReassignTasksResponse) Reset() {
	*x = ReassignTasksResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tracker_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReassignTasksResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReassignTasksResponse) ProtoMessage() {}

func (x *ReassignTasksResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tracker_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReassignTasksResponse.ProtoReflect.Descriptor instead.
func (*ReassignTasksResponse) Descriptor() ([]byte, []int) {
	return file_tracker_proto_rawDescGZIP(), []int{3}
}

type CompleteTaskRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskId string `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
}

func (x *CompleteTaskRequest) Reset() {
	*x = CompleteTaskRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tracker_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CompleteTaskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompleteTaskRequest) ProtoMessage() {}

func (x *CompleteTaskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tracker_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompleteTaskRequest.ProtoReflect.Descriptor instead.
func (*CompleteTaskRequest) Descriptor() ([]byte, []int) {
	return file_tracker_proto_rawDescGZIP(), []int{4}
}

func (x *CompleteTaskRequest) GetTaskId() string {
	if x != nil {
		return x.TaskId
	}
	return ""
}

type CompleteTaskResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CompleteTaskResponse) Reset() {
	*x = CompleteTaskResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tracker_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CompleteTaskResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompleteTaskResponse) ProtoMessage() {}

func (x *CompleteTaskResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tracker_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompleteTaskResponse.ProtoReflect.Descriptor instead.
func (*CompleteTaskResponse) Descriptor() ([]byte, []int) {
	return file_tracker_proto_rawDescGZIP(), []int{5}
}

var File_tracker_proto protoreflect.FileDescriptor

var file_tracker_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x35, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x14, 0x0a,
	0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x16, 0x0a, 0x14, 0x52, 0x65, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x54,
	0x61, 0x73, 0x6b, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x17, 0x0a, 0x15, 0x52,
	0x65, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x54, 0x61, 0x73, 0x6b, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x2e, 0x0a, 0x13, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65,
	0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x74,
	0x61, 0x73, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x61,
	0x73, 0x6b, 0x49, 0x64, 0x22, 0x16, 0x0a, 0x14, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65,
	0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xc5, 0x02, 0x0a,
	0x07, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x12, 0x61, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x1a, 0x2e, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x1a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x14, 0x22, 0x0f, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x61, 0x73,
	0x6b, 0x2f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x3a, 0x01, 0x2a, 0x12, 0x6c, 0x0a, 0x0d, 0x52,
	0x65, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x54, 0x61, 0x73, 0x6b, 0x73, 0x12, 0x1d, 0x2e, 0x74,
	0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x54,
	0x61, 0x73, 0x6b, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x74, 0x72,
	0x61, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x54, 0x61,
	0x73, 0x6b, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1c, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x16, 0x22, 0x11, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x61, 0x73, 0x6b, 0x2f, 0x72, 0x65,
	0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x3a, 0x01, 0x2a, 0x12, 0x69, 0x0a, 0x0c, 0x43, 0x6f, 0x6d,
	0x70, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x1c, 0x2e, 0x74, 0x72, 0x61, 0x63,
	0x6b, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x61, 0x73, 0x6b,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65,
	0x72, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x16, 0x22, 0x11,
	0x2f, 0x76, 0x31, 0x2f, 0x74, 0x61, 0x73, 0x6b, 0x2f, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74,
	0x65, 0x3a, 0x01, 0x2a, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x69, 0x6c, 0x79, 0x61, 0x74, 0x6f, 0x73, 0x2f, 0x61, 0x61, 0x63, 0x2d, 0x74,
	0x61, 0x73, 0x6b, 0x2d, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x2f, 0x74, 0x72, 0x61, 0x63,
	0x6b, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_tracker_proto_rawDescOnce sync.Once
	file_tracker_proto_rawDescData = file_tracker_proto_rawDesc
)

func file_tracker_proto_rawDescGZIP() []byte {
	file_tracker_proto_rawDescOnce.Do(func() {
		file_tracker_proto_rawDescData = protoimpl.X.CompressGZIP(file_tracker_proto_rawDescData)
	})
	return file_tracker_proto_rawDescData
}

var file_tracker_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_tracker_proto_goTypes = []interface{}{
	(*CreateTaskRequest)(nil),     // 0: tracker.CreateTaskRequest
	(*CreateTaskResponse)(nil),    // 1: tracker.CreateTaskResponse
	(*ReassignTasksRequest)(nil),  // 2: tracker.ReassignTasksRequest
	(*ReassignTasksResponse)(nil), // 3: tracker.ReassignTasksResponse
	(*CompleteTaskRequest)(nil),   // 4: tracker.CompleteTaskRequest
	(*CompleteTaskResponse)(nil),  // 5: tracker.CompleteTaskResponse
}
var file_tracker_proto_depIdxs = []int32{
	0, // 0: tracker.Tracker.CreateTask:input_type -> tracker.CreateTaskRequest
	2, // 1: tracker.Tracker.ReassignTasks:input_type -> tracker.ReassignTasksRequest
	4, // 2: tracker.Tracker.CompleteTask:input_type -> tracker.CompleteTaskRequest
	1, // 3: tracker.Tracker.CreateTask:output_type -> tracker.CreateTaskResponse
	3, // 4: tracker.Tracker.ReassignTasks:output_type -> tracker.ReassignTasksResponse
	5, // 5: tracker.Tracker.CompleteTask:output_type -> tracker.CompleteTaskResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_tracker_proto_init() }
func file_tracker_proto_init() {
	if File_tracker_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_tracker_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateTaskRequest); i {
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
		file_tracker_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateTaskResponse); i {
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
		file_tracker_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReassignTasksRequest); i {
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
		file_tracker_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReassignTasksResponse); i {
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
		file_tracker_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CompleteTaskRequest); i {
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
		file_tracker_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CompleteTaskResponse); i {
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
			RawDescriptor: file_tracker_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_tracker_proto_goTypes,
		DependencyIndexes: file_tracker_proto_depIdxs,
		MessageInfos:      file_tracker_proto_msgTypes,
	}.Build()
	File_tracker_proto = out.File
	file_tracker_proto_rawDesc = nil
	file_tracker_proto_goTypes = nil
	file_tracker_proto_depIdxs = nil
}
