// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v3.12.4
// source: proto/send_mail.proto

package send_mail

import (
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type MailRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FromDate     *timestamp.Timestamp `protobuf:"bytes,1,opt,name=fromDate,proto3" json:"fromDate,omitempty"`
	ToDate       *timestamp.Timestamp `protobuf:"bytes,2,opt,name=toDate,proto3" json:"toDate,omitempty"`
	MailReceiver string               `protobuf:"bytes,3,opt,name=mailReceiver,proto3" json:"mailReceiver,omitempty"`
}

func (x *MailRequest) Reset() {
	*x = MailRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_send_mail_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MailRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MailRequest) ProtoMessage() {}

func (x *MailRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_send_mail_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MailRequest.ProtoReflect.Descriptor instead.
func (*MailRequest) Descriptor() ([]byte, []int) {
	return file_proto_send_mail_proto_rawDescGZIP(), []int{0}
}

func (x *MailRequest) GetFromDate() *timestamp.Timestamp {
	if x != nil {
		return x.FromDate
	}
	return nil
}

func (x *MailRequest) GetToDate() *timestamp.Timestamp {
	if x != nil {
		return x.ToDate
	}
	return nil
}

func (x *MailRequest) GetMailReceiver() string {
	if x != nil {
		return x.MailReceiver
	}
	return ""
}

type MailResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsSuccess bool   `protobuf:"varint,1,opt,name=isSuccess,proto3" json:"isSuccess,omitempty"`
	FilePath  string `protobuf:"bytes,2,opt,name=filePath,proto3" json:"filePath,omitempty"`
}

func (x *MailResponse) Reset() {
	*x = MailResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_send_mail_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MailResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MailResponse) ProtoMessage() {}

func (x *MailResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_send_mail_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MailResponse.ProtoReflect.Descriptor instead.
func (*MailResponse) Descriptor() ([]byte, []int) {
	return file_proto_send_mail_proto_rawDescGZIP(), []int{1}
}

func (x *MailResponse) GetIsSuccess() bool {
	if x != nil {
		return x.IsSuccess
	}
	return false
}

func (x *MailResponse) GetFilePath() string {
	if x != nil {
		return x.FilePath
	}
	return ""
}

var File_proto_send_mail_proto protoreflect.FileDescriptor

var file_proto_send_mail_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x6e, 0x64, 0x5f, 0x6d, 0x61, 0x69,
	0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x73, 0x65, 0x6e, 0x64, 0x5f, 0x6d, 0x61,
	0x69, 0x6c, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x9d, 0x01, 0x0a, 0x0b, 0x4d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x36, 0x0a, 0x08, 0x66, 0x72, 0x6f, 0x6d, 0x44, 0x61, 0x74, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x08, 0x66, 0x72, 0x6f, 0x6d, 0x44, 0x61, 0x74, 0x65, 0x12, 0x32, 0x0a, 0x06, 0x74,
	0x6f, 0x44, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x06, 0x74, 0x6f, 0x44, 0x61, 0x74, 0x65, 0x12,
	0x22, 0x0a, 0x0c, 0x6d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x6d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x63, 0x65, 0x69,
	0x76, 0x65, 0x72, 0x22, 0x48, 0x0a, 0x0c, 0x4d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x69, 0x73, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x69, 0x73, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x32, 0x49, 0x0a,
	0x08, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x61, 0x69, 0x6c, 0x12, 0x3d, 0x0a, 0x0a, 0x44, 0x6f, 0x53,
	0x65, 0x6e, 0x64, 0x4d, 0x61, 0x69, 0x6c, 0x12, 0x16, 0x2e, 0x73, 0x65, 0x6e, 0x64, 0x5f, 0x6d,
	0x61, 0x69, 0x6c, 0x2e, 0x4d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x17, 0x2e, 0x73, 0x65, 0x6e, 0x64, 0x5f, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x4d, 0x61, 0x69, 0x6c,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x13, 0x5a, 0x11, 0x2e, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x6e, 0x64, 0x5f, 0x6d, 0x61, 0x69, 0x6c, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_send_mail_proto_rawDescOnce sync.Once
	file_proto_send_mail_proto_rawDescData = file_proto_send_mail_proto_rawDesc
)

func file_proto_send_mail_proto_rawDescGZIP() []byte {
	file_proto_send_mail_proto_rawDescOnce.Do(func() {
		file_proto_send_mail_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_send_mail_proto_rawDescData)
	})
	return file_proto_send_mail_proto_rawDescData
}

var file_proto_send_mail_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_send_mail_proto_goTypes = []interface{}{
	(*MailRequest)(nil),         // 0: send_mail.MailRequest
	(*MailResponse)(nil),        // 1: send_mail.MailResponse
	(*timestamp.Timestamp)(nil), // 2: google.protobuf.Timestamp
}
var file_proto_send_mail_proto_depIdxs = []int32{
	2, // 0: send_mail.MailRequest.fromDate:type_name -> google.protobuf.Timestamp
	2, // 1: send_mail.MailRequest.toDate:type_name -> google.protobuf.Timestamp
	0, // 2: send_mail.SendMail.DoSendMail:input_type -> send_mail.MailRequest
	1, // 3: send_mail.SendMail.DoSendMail:output_type -> send_mail.MailResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_send_mail_proto_init() }
func file_proto_send_mail_proto_init() {
	if File_proto_send_mail_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_send_mail_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MailRequest); i {
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
		file_proto_send_mail_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MailResponse); i {
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
			RawDescriptor: file_proto_send_mail_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_send_mail_proto_goTypes,
		DependencyIndexes: file_proto_send_mail_proto_depIdxs,
		MessageInfos:      file_proto_send_mail_proto_msgTypes,
	}.Build()
	File_proto_send_mail_proto = out.File
	file_proto_send_mail_proto_rawDesc = nil
	file_proto_send_mail_proto_goTypes = nil
	file_proto_send_mail_proto_depIdxs = nil
}