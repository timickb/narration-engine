// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: worker.proto

package stateflow

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

// HandleRequest - структура с входными данными для обработчика.
type HandleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// instance_id - идентификатор экземпляра.
	InstanceId string `protobuf:"bytes,1,opt,name=instance_id,json=instanceId,proto3" json:"instance_id,omitempty"`
	// context - данные экземпляра (json)
	Context string `protobuf:"bytes,2,opt,name=context,proto3" json:"context,omitempty"`
	// event_name = событие, по которому произошел переход в состояние;
	EventName string `protobuf:"bytes,3,opt,name=event_name,json=eventName,proto3" json:"event_name,omitempty"`
	// event_params - данные, переданные с событием
	EventParams string `protobuf:"bytes,4,opt,name=event_params,json=eventParams,proto3" json:"event_params,omitempty"`
	// state - название соответствующего состояния.
	State string `protobuf:"bytes,5,opt,name=state,proto3" json:"state,omitempty"`
	// handler - название обработчика.
	Handler string `protobuf:"bytes,6,opt,name=handler,proto3" json:"handler,omitempty"`
}

func (x *HandleRequest) Reset() {
	*x = HandleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_worker_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HandleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HandleRequest) ProtoMessage() {}

func (x *HandleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_worker_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HandleRequest.ProtoReflect.Descriptor instead.
func (*HandleRequest) Descriptor() ([]byte, []int) {
	return file_worker_proto_rawDescGZIP(), []int{0}
}

func (x *HandleRequest) GetInstanceId() string {
	if x != nil {
		return x.InstanceId
	}
	return ""
}

func (x *HandleRequest) GetContext() string {
	if x != nil {
		return x.Context
	}
	return ""
}

func (x *HandleRequest) GetEventName() string {
	if x != nil {
		return x.EventName
	}
	return ""
}

func (x *HandleRequest) GetEventParams() string {
	if x != nil {
		return x.EventParams
	}
	return ""
}

func (x *HandleRequest) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *HandleRequest) GetHandler() string {
	if x != nil {
		return x.Handler
	}
	return ""
}

// HandleResponse - структура с результатом работы обработчика.
type HandleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// status - статус ответа.
	Status *Status `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	// next_event - следующее событие
	NextEvent string `protobuf:"bytes,2,opt,name=next_event,json=nextEvent,proto3" json:"next_event,omitempty"`
	// next_event_payload - данные для следующего события (json).
	NextEventPayload string `protobuf:"bytes,3,opt,name=next_event_payload,json=nextEventPayload,proto3" json:"next_event_payload,omitempty"`
	// data_to_context - данные для сохранения в контексте экземпляра.
	DataToContext string `protobuf:"bytes,4,opt,name=data_to_context,json=dataToContext,proto3" json:"data_to_context,omitempty"`
}

func (x *HandleResponse) Reset() {
	*x = HandleResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_worker_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HandleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HandleResponse) ProtoMessage() {}

func (x *HandleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_worker_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HandleResponse.ProtoReflect.Descriptor instead.
func (*HandleResponse) Descriptor() ([]byte, []int) {
	return file_worker_proto_rawDescGZIP(), []int{1}
}

func (x *HandleResponse) GetStatus() *Status {
	if x != nil {
		return x.Status
	}
	return nil
}

func (x *HandleResponse) GetNextEvent() string {
	if x != nil {
		return x.NextEvent
	}
	return ""
}

func (x *HandleResponse) GetNextEventPayload() string {
	if x != nil {
		return x.NextEventPayload
	}
	return ""
}

func (x *HandleResponse) GetDataToContext() string {
	if x != nil {
		return x.DataToContext
	}
	return ""
}

var File_worker_proto protoreflect.FileDescriptor

var file_worker_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09,
	0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e, 0x76, 0x31, 0x1a, 0x0b, 0x74, 0x79, 0x70, 0x65, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xbc, 0x01, 0x0a, 0x0d, 0x48, 0x61, 0x6e, 0x64, 0x6c,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x6e, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x69,
	0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x78, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x70, 0x61, 0x72, 0x61,
	0x6d, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x68,
	0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x68, 0x61,
	0x6e, 0x64, 0x6c, 0x65, 0x72, 0x22, 0xb0, 0x01, 0x0a, 0x0e, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x73, 0x63, 0x68, 0x65, 0x6d,
	0x61, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x65, 0x78, 0x74, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x12, 0x2c, 0x0a, 0x12, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x5f, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10,
	0x6e, 0x65, 0x78, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64,
	0x12, 0x26, 0x0a, 0x0f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x74, 0x6f, 0x5f, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x78, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x64, 0x61, 0x74, 0x61, 0x54,
	0x6f, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x32, 0x4e, 0x0a, 0x0d, 0x57, 0x6f, 0x72, 0x6b,
	0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3d, 0x0a, 0x06, 0x48, 0x61, 0x6e,
	0x64, 0x6c, 0x65, 0x12, 0x18, 0x2e, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e, 0x76, 0x31, 0x2e,
	0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e,
	0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x43, 0x5a, 0x41, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x69, 0x6d, 0x69, 0x63, 0x6b, 0x62, 0x2f, 0x6e,
	0x61, 0x72, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2d, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x66,
	0x6c, 0x6f, 0x77, 0x3b, 0x73, 0x74, 0x61, 0x74, 0x65, 0x66, 0x6c, 0x6f, 0x77, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_worker_proto_rawDescOnce sync.Once
	file_worker_proto_rawDescData = file_worker_proto_rawDesc
)

func file_worker_proto_rawDescGZIP() []byte {
	file_worker_proto_rawDescOnce.Do(func() {
		file_worker_proto_rawDescData = protoimpl.X.CompressGZIP(file_worker_proto_rawDescData)
	})
	return file_worker_proto_rawDescData
}

var file_worker_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_worker_proto_goTypes = []interface{}{
	(*HandleRequest)(nil),  // 0: schema.v1.HandleRequest
	(*HandleResponse)(nil), // 1: schema.v1.HandleResponse
	(*Status)(nil),         // 2: schema.v1.Status
}
var file_worker_proto_depIdxs = []int32{
	2, // 0: schema.v1.HandleResponse.status:type_name -> schema.v1.Status
	0, // 1: schema.v1.WorkerService.Handle:input_type -> schema.v1.HandleRequest
	1, // 2: schema.v1.WorkerService.Handle:output_type -> schema.v1.HandleResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_worker_proto_init() }
func file_worker_proto_init() {
	if File_worker_proto != nil {
		return
	}
	file_types_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_worker_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HandleRequest); i {
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
		file_worker_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HandleResponse); i {
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
			RawDescriptor: file_worker_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_worker_proto_goTypes,
		DependencyIndexes: file_worker_proto_depIdxs,
		MessageInfos:      file_worker_proto_msgTypes,
	}.Build()
	File_worker_proto = out.File
	file_worker_proto_rawDesc = nil
	file_worker_proto_goTypes = nil
	file_worker_proto_depIdxs = nil
}
