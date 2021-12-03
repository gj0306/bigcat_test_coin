// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.18.1
// source: app/blc/contract/contract.proto

package contract

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

type News struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title   string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Content string `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *News) Reset() {
	*x = News{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_blc_contract_contract_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *News) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*News) ProtoMessage() {}

func (x *News) ProtoReflect() protoreflect.Message {
	mi := &file_app_blc_contract_contract_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use News.ProtoReflect.Descriptor instead.
func (*News) Descriptor() ([]byte, []int) {
	return file_app_blc_contract_contract_proto_rawDescGZIP(), []int{0}
}

func (x *News) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *News) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

type Profile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Head      string `protobuf:"bytes,2,opt,name=head,proto3" json:"head,omitempty"`
	Introduce string `protobuf:"bytes,3,opt,name=introduce,proto3" json:"introduce,omitempty"`
}

func (x *Profile) Reset() {
	*x = Profile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_blc_contract_contract_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Profile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Profile) ProtoMessage() {}

func (x *Profile) ProtoReflect() protoreflect.Message {
	mi := &file_app_blc_contract_contract_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Profile.ProtoReflect.Descriptor instead.
func (*Profile) Descriptor() ([]byte, []int) {
	return file_app_blc_contract_contract_proto_rawDescGZIP(), []int{1}
}

func (x *Profile) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Profile) GetHead() string {
	if x != nil {
		return x.Head
	}
	return ""
}

func (x *Profile) GetIntroduce() string {
	if x != nil {
		return x.Introduce
	}
	return ""
}

type Source struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Quote   []byte `protobuf:"bytes,1,opt,name=quote,proto3" json:"quote,omitempty"`
	From    string `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	To      string `protobuf:"bytes,3,opt,name=to,proto3" json:"to,omitempty"`
	Content string `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *Source) Reset() {
	*x = Source{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_blc_contract_contract_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Source) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Source) ProtoMessage() {}

func (x *Source) ProtoReflect() protoreflect.Message {
	mi := &file_app_blc_contract_contract_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Source.ProtoReflect.Descriptor instead.
func (*Source) Descriptor() ([]byte, []int) {
	return file_app_blc_contract_contract_proto_rawDescGZIP(), []int{2}
}

func (x *Source) GetQuote() []byte {
	if x != nil {
		return x.Quote
	}
	return nil
}

func (x *Source) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *Source) GetTo() string {
	if x != nil {
		return x.To
	}
	return ""
}

func (x *Source) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

var File_app_blc_contract_contract_proto protoreflect.FileDescriptor

var file_app_blc_contract_contract_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x61, 0x70, 0x70, 0x2f, 0x62, 0x6c, 0x63, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61,
	0x63, 0x74, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x22, 0x36, 0x0a, 0x04, 0x6e,
	0x65, 0x77, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x22, 0x4f, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x65, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x68, 0x65, 0x61, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x69, 0x6e, 0x74, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x69, 0x6e, 0x74, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x65, 0x22, 0x5c, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x71, 0x75, 0x6f, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x71,
	0x75, 0x6f, 0x74, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x6f, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x74, 0x6f, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x42, 0x4b, 0x0a, 0x16, 0x64, 0x65, 0x76, 0x2e, 0x6b, 0x72, 0x61, 0x74, 0x6f, 0x73,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x6f, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x42, 0x0a, 0x43, 0x6f,
	0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x56, 0x31, 0x50, 0x01, 0x5a, 0x23, 0x6e, 0x65, 0x77, 0x73,
	0x5f, 0x63, 0x6f, 0x69, 0x6e, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x62, 0x6c, 0x63, 0x2f, 0x63, 0x6f,
	0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x3b, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_app_blc_contract_contract_proto_rawDescOnce sync.Once
	file_app_blc_contract_contract_proto_rawDescData = file_app_blc_contract_contract_proto_rawDesc
)

func file_app_blc_contract_contract_proto_rawDescGZIP() []byte {
	file_app_blc_contract_contract_proto_rawDescOnce.Do(func() {
		file_app_blc_contract_contract_proto_rawDescData = protoimpl.X.CompressGZIP(file_app_blc_contract_contract_proto_rawDescData)
	})
	return file_app_blc_contract_contract_proto_rawDescData
}

var file_app_blc_contract_contract_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_app_blc_contract_contract_proto_goTypes = []interface{}{
	(*News)(nil),    // 0: contract.news
	(*Profile)(nil), // 1: contract.profile
	(*Source)(nil),  // 2: contract.source
}
var file_app_blc_contract_contract_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_app_blc_contract_contract_proto_init() }
func file_app_blc_contract_contract_proto_init() {
	if File_app_blc_contract_contract_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_app_blc_contract_contract_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*News); i {
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
		file_app_blc_contract_contract_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Profile); i {
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
		file_app_blc_contract_contract_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Source); i {
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
			RawDescriptor: file_app_blc_contract_contract_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_app_blc_contract_contract_proto_goTypes,
		DependencyIndexes: file_app_blc_contract_contract_proto_depIdxs,
		MessageInfos:      file_app_blc_contract_contract_proto_msgTypes,
	}.Build()
	File_app_blc_contract_contract_proto = out.File
	file_app_blc_contract_contract_proto_rawDesc = nil
	file_app_blc_contract_contract_proto_goTypes = nil
	file_app_blc_contract_contract_proto_depIdxs = nil
}
