// Test to benchmark quorum functions with and without the request parameter.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v3.12.4
// source: broadcast/broadcast.proto

package broadcast

import (
	_ "github.com/relab/gorums"
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

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value int64  `protobuf:"varint,1,opt,name=Value,proto3" json:"Value,omitempty"`
	From  string `protobuf:"bytes,2,opt,name=From,proto3" json:"From,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_broadcast_broadcast_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_broadcast_broadcast_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_broadcast_broadcast_proto_rawDescGZIP(), []int{0}
}

func (x *Request) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *Request) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result  int64  `protobuf:"varint,1,opt,name=Result,proto3" json:"Result,omitempty"`
	Replies int64  `protobuf:"varint,2,opt,name=Replies,proto3" json:"Replies,omitempty"`
	From    string `protobuf:"bytes,3,opt,name=From,proto3" json:"From,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_broadcast_broadcast_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_broadcast_broadcast_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_broadcast_broadcast_proto_rawDescGZIP(), []int{1}
}

func (x *Response) GetResult() int64 {
	if x != nil {
		return x.Result
	}
	return 0
}

func (x *Response) GetReplies() int64 {
	if x != nil {
		return x.Replies
	}
	return 0
}

func (x *Response) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_broadcast_broadcast_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_broadcast_broadcast_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_broadcast_broadcast_proto_rawDescGZIP(), []int{2}
}

var File_broadcast_broadcast_proto protoreflect.FileDescriptor

var file_broadcast_broadcast_proto_rawDesc = []byte{
	0x0a, 0x19, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2f, 0x62, 0x72, 0x6f, 0x61,
	0x64, 0x63, 0x61, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x62, 0x72, 0x6f,
	0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x1a, 0x0c, 0x67, 0x6f, 0x72, 0x75, 0x6d, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x33, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x22, 0x50, 0x0a, 0x08, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x18, 0x0a,
	0x07, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07,
	0x52, 0x65, 0x70, 0x6c, 0x69, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x22, 0x07, 0x0a, 0x05, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x32, 0x82, 0x09, 0x0a, 0x10, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61,
	0x73, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3b, 0x0a, 0x0a, 0x51, 0x75, 0x6f,
	0x72, 0x75, 0x6d, 0x43, 0x61, 0x6c, 0x6c, 0x12, 0x12, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63,
	0x61, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x62, 0x72,
	0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x04, 0xa0, 0xb5, 0x18, 0x01, 0x12, 0x4c, 0x0a, 0x17, 0x51, 0x75, 0x6f, 0x72, 0x75, 0x6d,
	0x43, 0x61, 0x6c, 0x6c, 0x57, 0x69, 0x74, 0x68, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73,
	0x74, 0x12, 0x12, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73,
	0x74, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x08, 0xa0, 0xb5, 0x18, 0x01,
	0x90, 0xb8, 0x18, 0x01, 0x12, 0x48, 0x0a, 0x17, 0x51, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x43, 0x61,
	0x6c, 0x6c, 0x57, 0x69, 0x74, 0x68, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x63, 0x61, 0x73, 0x74, 0x12,
	0x12, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x04, 0xa0, 0xb5, 0x18, 0x01, 0x12, 0x37,
	0x0a, 0x09, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x63, 0x61, 0x73, 0x74, 0x12, 0x12, 0x2e, 0x62, 0x72,
	0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x10, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x04, 0x98, 0xb5, 0x18, 0x01, 0x12, 0x43, 0x0a, 0x15, 0x4d, 0x75, 0x6c, 0x74, 0x69,
	0x63, 0x61, 0x73, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x74, 0x65,
	0x12, 0x12, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x04, 0x98, 0xb5, 0x18, 0x01, 0x12, 0x3e, 0x0a, 0x0d,
	0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x43, 0x61, 0x6c, 0x6c, 0x12, 0x12, 0x2e,
	0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x13, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x04, 0xb0, 0xb5, 0x18, 0x01, 0x12, 0x43, 0x0a, 0x15,
	0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6d, 0x65,
	0x64, 0x69, 0x61, 0x74, 0x65, 0x12, 0x12, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73,
	0x74, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x62, 0x72, 0x6f, 0x61,
	0x64, 0x63, 0x61, 0x73, 0x74, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x04, 0x90, 0xb8, 0x18,
	0x01, 0x12, 0x37, 0x0a, 0x09, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x12, 0x12,
	0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x10, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x04, 0x90, 0xb8, 0x18, 0x01, 0x12, 0x45, 0x0a, 0x14, 0x42, 0x72,
	0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x43, 0x61, 0x6c, 0x6c, 0x46, 0x6f, 0x72, 0x77, 0x61,
	0x72, 0x64, 0x12, 0x12, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61,
	0x73, 0x74, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x04, 0xb0, 0xb5, 0x18,
	0x01, 0x12, 0x40, 0x0a, 0x0f, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x43, 0x61,
	0x6c, 0x6c, 0x54, 0x6f, 0x12, 0x12, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74,
	0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64,
	0x63, 0x61, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x04, 0xb0,
	0xb5, 0x18, 0x01, 0x12, 0x41, 0x0a, 0x13, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74,
	0x54, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x2e, 0x62, 0x72, 0x6f,
	0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10,
	0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x04, 0x90, 0xb8, 0x18, 0x01, 0x12, 0x37, 0x0a, 0x06, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x12, 0x12, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74,
	0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x04, 0xb0, 0xb5, 0x18, 0x01, 0x12,
	0x40, 0x0a, 0x0f, 0x4c, 0x6f, 0x6e, 0x67, 0x52, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x54, 0x61,
	0x73, 0x6b, 0x12, 0x12, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61,
	0x73, 0x74, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x04, 0xb0, 0xb5, 0x18,
	0x01, 0x12, 0x37, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x56, 0x61, 0x6c, 0x12, 0x12, 0x2e, 0x62, 0x72,
	0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x13, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x04, 0xb0, 0xb5, 0x18, 0x01, 0x12, 0x36, 0x0a, 0x05, 0x4f, 0x72,
	0x64, 0x65, 0x72, 0x12, 0x12, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63,
	0x61, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x04, 0xb0, 0xb5,
	0x18, 0x01, 0x12, 0x38, 0x0a, 0x0a, 0x50, 0x72, 0x65, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65,
	0x12, 0x12, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x04, 0x90, 0xb8, 0x18, 0x01, 0x12, 0x35, 0x0a, 0x07,
	0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x12, 0x12, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63,
	0x61, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x62, 0x72,
	0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x04, 0x90,
	0xb8, 0x18, 0x01, 0x12, 0x34, 0x0a, 0x06, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x12, 0x2e,
	0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x10, 0x2e, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x04, 0x90, 0xb8, 0x18, 0x01, 0x42, 0x29, 0x5a, 0x27, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x65, 0x6c, 0x61, 0x62, 0x2f, 0x67, 0x6f,
	0x72, 0x75, 0x6d, 0x73, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x73, 0x2f, 0x62, 0x72, 0x6f, 0x61, 0x64,
	0x63, 0x61, 0x73, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_broadcast_broadcast_proto_rawDescOnce sync.Once
	file_broadcast_broadcast_proto_rawDescData = file_broadcast_broadcast_proto_rawDesc
)

func file_broadcast_broadcast_proto_rawDescGZIP() []byte {
	file_broadcast_broadcast_proto_rawDescOnce.Do(func() {
		file_broadcast_broadcast_proto_rawDescData = protoimpl.X.CompressGZIP(file_broadcast_broadcast_proto_rawDescData)
	})
	return file_broadcast_broadcast_proto_rawDescData
}

var file_broadcast_broadcast_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_broadcast_broadcast_proto_goTypes = []interface{}{
	(*Request)(nil),  // 0: broadcast.Request
	(*Response)(nil), // 1: broadcast.Response
	(*Empty)(nil),    // 2: broadcast.Empty
}
var file_broadcast_broadcast_proto_depIdxs = []int32{
	0,  // 0: broadcast.BroadcastService.QuorumCall:input_type -> broadcast.Request
	0,  // 1: broadcast.BroadcastService.QuorumCallWithBroadcast:input_type -> broadcast.Request
	0,  // 2: broadcast.BroadcastService.QuorumCallWithMulticast:input_type -> broadcast.Request
	0,  // 3: broadcast.BroadcastService.Multicast:input_type -> broadcast.Request
	0,  // 4: broadcast.BroadcastService.MulticastIntermediate:input_type -> broadcast.Request
	0,  // 5: broadcast.BroadcastService.BroadcastCall:input_type -> broadcast.Request
	0,  // 6: broadcast.BroadcastService.BroadcastIntermediate:input_type -> broadcast.Request
	0,  // 7: broadcast.BroadcastService.Broadcast:input_type -> broadcast.Request
	0,  // 8: broadcast.BroadcastService.BroadcastCallForward:input_type -> broadcast.Request
	0,  // 9: broadcast.BroadcastService.BroadcastCallTo:input_type -> broadcast.Request
	0,  // 10: broadcast.BroadcastService.BroadcastToResponse:input_type -> broadcast.Request
	0,  // 11: broadcast.BroadcastService.Search:input_type -> broadcast.Request
	0,  // 12: broadcast.BroadcastService.LongRunningTask:input_type -> broadcast.Request
	0,  // 13: broadcast.BroadcastService.GetVal:input_type -> broadcast.Request
	0,  // 14: broadcast.BroadcastService.Order:input_type -> broadcast.Request
	0,  // 15: broadcast.BroadcastService.PrePrepare:input_type -> broadcast.Request
	0,  // 16: broadcast.BroadcastService.Prepare:input_type -> broadcast.Request
	0,  // 17: broadcast.BroadcastService.Commit:input_type -> broadcast.Request
	1,  // 18: broadcast.BroadcastService.QuorumCall:output_type -> broadcast.Response
	1,  // 19: broadcast.BroadcastService.QuorumCallWithBroadcast:output_type -> broadcast.Response
	1,  // 20: broadcast.BroadcastService.QuorumCallWithMulticast:output_type -> broadcast.Response
	2,  // 21: broadcast.BroadcastService.Multicast:output_type -> broadcast.Empty
	2,  // 22: broadcast.BroadcastService.MulticastIntermediate:output_type -> broadcast.Empty
	1,  // 23: broadcast.BroadcastService.BroadcastCall:output_type -> broadcast.Response
	2,  // 24: broadcast.BroadcastService.BroadcastIntermediate:output_type -> broadcast.Empty
	2,  // 25: broadcast.BroadcastService.Broadcast:output_type -> broadcast.Empty
	1,  // 26: broadcast.BroadcastService.BroadcastCallForward:output_type -> broadcast.Response
	1,  // 27: broadcast.BroadcastService.BroadcastCallTo:output_type -> broadcast.Response
	2,  // 28: broadcast.BroadcastService.BroadcastToResponse:output_type -> broadcast.Empty
	1,  // 29: broadcast.BroadcastService.Search:output_type -> broadcast.Response
	1,  // 30: broadcast.BroadcastService.LongRunningTask:output_type -> broadcast.Response
	1,  // 31: broadcast.BroadcastService.GetVal:output_type -> broadcast.Response
	1,  // 32: broadcast.BroadcastService.Order:output_type -> broadcast.Response
	2,  // 33: broadcast.BroadcastService.PrePrepare:output_type -> broadcast.Empty
	2,  // 34: broadcast.BroadcastService.Prepare:output_type -> broadcast.Empty
	2,  // 35: broadcast.BroadcastService.Commit:output_type -> broadcast.Empty
	18, // [18:36] is the sub-list for method output_type
	0,  // [0:18] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_broadcast_broadcast_proto_init() }
func file_broadcast_broadcast_proto_init() {
	if File_broadcast_broadcast_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_broadcast_broadcast_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
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
		file_broadcast_broadcast_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
		file_broadcast_broadcast_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
			RawDescriptor: file_broadcast_broadcast_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_broadcast_broadcast_proto_goTypes,
		DependencyIndexes: file_broadcast_broadcast_proto_depIdxs,
		MessageInfos:      file_broadcast_broadcast_proto_msgTypes,
	}.Build()
	File_broadcast_broadcast_proto = out.File
	file_broadcast_broadcast_proto_rawDesc = nil
	file_broadcast_broadcast_proto_goTypes = nil
	file_broadcast_broadcast_proto_depIdxs = nil
}
