// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: model/profile.proto

package model

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

type SmallProfile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	SecondName  string `protobuf:"bytes,2,opt,name=second_name,json=secondName,proto3" json:"second_name,omitempty"`
	Avatar      string `protobuf:"bytes,3,opt,name=avatar,proto3" json:"avatar,omitempty"`
	Description string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Birthday    string `protobuf:"bytes,5,opt,name=birthday,proto3" json:"birthday,omitempty"`
	Gender      int32  `protobuf:"varint,6,opt,name=gender,proto3" json:"gender,omitempty"`
	Email       string `protobuf:"bytes,7,opt,name=email,proto3" json:"email,omitempty"`
}

func (x *SmallProfile) Reset() {
	*x = SmallProfile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_model_profile_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SmallProfile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SmallProfile) ProtoMessage() {}

func (x *SmallProfile) ProtoReflect() protoreflect.Message {
	mi := &file_model_profile_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SmallProfile.ProtoReflect.Descriptor instead.
func (*SmallProfile) Descriptor() ([]byte, []int) {
	return file_model_profile_proto_rawDescGZIP(), []int{0}
}

func (x *SmallProfile) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SmallProfile) GetSecondName() string {
	if x != nil {
		return x.SecondName
	}
	return ""
}

func (x *SmallProfile) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

func (x *SmallProfile) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *SmallProfile) GetBirthday() string {
	if x != nil {
		return x.Birthday
	}
	return ""
}

func (x *SmallProfile) GetGender() int32 {
	if x != nil {
		return x.Gender
	}
	return 0
}

func (x *SmallProfile) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type UpdateSmallProfile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	SecondName  string `protobuf:"bytes,2,opt,name=second_name,json=secondName,proto3" json:"second_name,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Birthday    string `protobuf:"bytes,4,opt,name=birthday,proto3" json:"birthday,omitempty"`
	Gender      int32  `protobuf:"varint,5,opt,name=gender,proto3" json:"gender,omitempty"`
	Email       string `protobuf:"bytes,6,opt,name=email,proto3" json:"email,omitempty"`
}

func (x *UpdateSmallProfile) Reset() {
	*x = UpdateSmallProfile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_model_profile_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateSmallProfile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateSmallProfile) ProtoMessage() {}

func (x *UpdateSmallProfile) ProtoReflect() protoreflect.Message {
	mi := &file_model_profile_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateSmallProfile.ProtoReflect.Descriptor instead.
func (*UpdateSmallProfile) Descriptor() ([]byte, []int) {
	return file_model_profile_proto_rawDescGZIP(), []int{1}
}

func (x *UpdateSmallProfile) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateSmallProfile) GetSecondName() string {
	if x != nil {
		return x.SecondName
	}
	return ""
}

func (x *UpdateSmallProfile) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *UpdateSmallProfile) GetBirthday() string {
	if x != nil {
		return x.Birthday
	}
	return ""
}

func (x *UpdateSmallProfile) GetGender() int32 {
	if x != nil {
		return x.Gender
	}
	return 0
}

func (x *UpdateSmallProfile) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

var File_model_profile_proto protoreflect.FileDescriptor

var file_model_profile_proto_rawDesc = []byte{
	0x0a, 0x13, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x22, 0xc7,
	0x01, 0x0a, 0x0c, 0x53, 0x6d, 0x61, 0x6c, 0x6c, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x12, 0x20, 0x0a, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a,
	0x0a, 0x08, 0x62, 0x69, 0x72, 0x74, 0x68, 0x64, 0x61, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x62, 0x69, 0x72, 0x74, 0x68, 0x64, 0x61, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x67, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x67, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0xb5, 0x01, 0x0a, 0x12, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x53, 0x6d, 0x61, 0x6c, 0x6c, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x69, 0x72, 0x74, 0x68, 0x64,
	0x61, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x62, 0x69, 0x72, 0x74, 0x68, 0x64,
	0x61, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x06, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d,
	0x61, 0x69, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c,
	0x42, 0x3a, 0x5a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x45,
	0x75, 0x67, 0x75, 0x6e, 0x65, 0x2d, 0x55, 0x73, 0x61, 0x63, 0x68, 0x65, 0x76, 0x2f, 0x73, 0x6f,
	0x63, 0x69, 0x61, 0x6c, 0x2d, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x3b, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_model_profile_proto_rawDescOnce sync.Once
	file_model_profile_proto_rawDescData = file_model_profile_proto_rawDesc
)

func file_model_profile_proto_rawDescGZIP() []byte {
	file_model_profile_proto_rawDescOnce.Do(func() {
		file_model_profile_proto_rawDescData = protoimpl.X.CompressGZIP(file_model_profile_proto_rawDescData)
	})
	return file_model_profile_proto_rawDescData
}

var file_model_profile_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_model_profile_proto_goTypes = []interface{}{
	(*SmallProfile)(nil),       // 0: profile.SmallProfile
	(*UpdateSmallProfile)(nil), // 1: profile.UpdateSmallProfile
}
var file_model_profile_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_model_profile_proto_init() }
func file_model_profile_proto_init() {
	if File_model_profile_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_model_profile_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SmallProfile); i {
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
		file_model_profile_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateSmallProfile); i {
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
			RawDescriptor: file_model_profile_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_model_profile_proto_goTypes,
		DependencyIndexes: file_model_profile_proto_depIdxs,
		MessageInfos:      file_model_profile_proto_msgTypes,
	}.Build()
	File_model_profile_proto = out.File
	file_model_profile_proto_rawDesc = nil
	file_model_profile_proto_goTypes = nil
	file_model_profile_proto_depIdxs = nil
}
