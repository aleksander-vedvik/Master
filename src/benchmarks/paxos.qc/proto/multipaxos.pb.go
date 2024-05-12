// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v3.12.4
// source: multipaxos.proto

package proto

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

type Value struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientID      string `protobuf:"bytes,1,opt,name=ClientID,proto3" json:"ClientID,omitempty"`
	ClientSeq     uint32 `protobuf:"varint,2,opt,name=ClientSeq,proto3" json:"ClientSeq,omitempty"`
	IsNoop        bool   `protobuf:"varint,3,opt,name=isNoop,proto3" json:"isNoop,omitempty"`
	ClientCommand string `protobuf:"bytes,4,opt,name=ClientCommand,proto3" json:"ClientCommand,omitempty"`
}

func (x *Value) Reset() {
	*x = Value{}
	if protoimpl.UnsafeEnabled {
		mi := &file_multipaxos_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Value) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Value) ProtoMessage() {}

func (x *Value) ProtoReflect() protoreflect.Message {
	mi := &file_multipaxos_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Value.ProtoReflect.Descriptor instead.
func (*Value) Descriptor() ([]byte, []int) {
	return file_multipaxos_proto_rawDescGZIP(), []int{0}
}

func (x *Value) GetClientID() string {
	if x != nil {
		return x.ClientID
	}
	return ""
}

func (x *Value) GetClientSeq() uint32 {
	if x != nil {
		return x.ClientSeq
	}
	return 0
}

func (x *Value) GetIsNoop() bool {
	if x != nil {
		return x.IsNoop
	}
	return false
}

func (x *Value) GetClientCommand() string {
	if x != nil {
		return x.ClientCommand
	}
	return ""
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientID      string `protobuf:"bytes,1,opt,name=ClientID,proto3" json:"ClientID,omitempty"`
	ClientSeq     uint32 `protobuf:"varint,2,opt,name=ClientSeq,proto3" json:"ClientSeq,omitempty"`
	ClientCommand string `protobuf:"bytes,3,opt,name=ClientCommand,proto3" json:"ClientCommand,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_multipaxos_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_multipaxos_proto_msgTypes[1]
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
	return file_multipaxos_proto_rawDescGZIP(), []int{1}
}

func (x *Response) GetClientID() string {
	if x != nil {
		return x.ClientID
	}
	return ""
}

func (x *Response) GetClientSeq() uint32 {
	if x != nil {
		return x.ClientSeq
	}
	return 0
}

func (x *Response) GetClientCommand() string {
	if x != nil {
		return x.ClientCommand
	}
	return ""
}

type PrepareMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Slot uint32 `protobuf:"varint,1,opt,name=Slot,proto3" json:"Slot,omitempty"`
	Crnd int32  `protobuf:"varint,2,opt,name=Crnd,proto3" json:"Crnd,omitempty"`
}

func (x *PrepareMsg) Reset() {
	*x = PrepareMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_multipaxos_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrepareMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrepareMsg) ProtoMessage() {}

func (x *PrepareMsg) ProtoReflect() protoreflect.Message {
	mi := &file_multipaxos_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrepareMsg.ProtoReflect.Descriptor instead.
func (*PrepareMsg) Descriptor() ([]byte, []int) {
	return file_multipaxos_proto_rawDescGZIP(), []int{2}
}

func (x *PrepareMsg) GetSlot() uint32 {
	if x != nil {
		return x.Slot
	}
	return 0
}

func (x *PrepareMsg) GetCrnd() int32 {
	if x != nil {
		return x.Crnd
	}
	return 0
}

// PromiseMsg is the reply from an Acceptor to the Proposer in response to a PrepareMsg.
// The Acceptor will only respond if the PrepareMsg.Rnd > Acceptor.Rnd.
type PromiseMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rnd      int32     `protobuf:"varint,1,opt,name=Rnd,proto3" json:"Rnd,omitempty"`
	Accepted []*PValue `protobuf:"bytes,2,rep,name=Accepted,proto3" json:"Accepted,omitempty"`
}

func (x *PromiseMsg) Reset() {
	*x = PromiseMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_multipaxos_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PromiseMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PromiseMsg) ProtoMessage() {}

func (x *PromiseMsg) ProtoReflect() protoreflect.Message {
	mi := &file_multipaxos_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PromiseMsg.ProtoReflect.Descriptor instead.
func (*PromiseMsg) Descriptor() ([]byte, []int) {
	return file_multipaxos_proto_rawDescGZIP(), []int{3}
}

func (x *PromiseMsg) GetRnd() int32 {
	if x != nil {
		return x.Rnd
	}
	return 0
}

func (x *PromiseMsg) GetAccepted() []*PValue {
	if x != nil {
		return x.Accepted
	}
	return nil
}

// AcceptMsg is sent by the Proposer, asking the Acceptors to lock-in the value, val.
// If AcceptMsg.rnd < Acceptor.rnd, the message will be ignored.
type AcceptMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Slot uint32 `protobuf:"varint,1,opt,name=Slot,proto3" json:"Slot,omitempty"`
	Rnd  int32  `protobuf:"varint,2,opt,name=Rnd,proto3" json:"Rnd,omitempty"`
	Val  *Value `protobuf:"bytes,3,opt,name=Val,proto3" json:"Val,omitempty"`
}

func (x *AcceptMsg) Reset() {
	*x = AcceptMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_multipaxos_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AcceptMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AcceptMsg) ProtoMessage() {}

func (x *AcceptMsg) ProtoReflect() protoreflect.Message {
	mi := &file_multipaxos_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AcceptMsg.ProtoReflect.Descriptor instead.
func (*AcceptMsg) Descriptor() ([]byte, []int) {
	return file_multipaxos_proto_rawDescGZIP(), []int{4}
}

func (x *AcceptMsg) GetSlot() uint32 {
	if x != nil {
		return x.Slot
	}
	return 0
}

func (x *AcceptMsg) GetRnd() int32 {
	if x != nil {
		return x.Rnd
	}
	return 0
}

func (x *AcceptMsg) GetVal() *Value {
	if x != nil {
		return x.Val
	}
	return nil
}

// LearnMsg is sent by an Acceptor to the Proposer, if the Acceptor agreed to lock-in the value, val.
// The LearnMsg is also sent by the Proposer in a Commit.
type LearnMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Slot uint32 `protobuf:"varint,1,opt,name=Slot,proto3" json:"Slot,omitempty"`
	Rnd  int32  `protobuf:"varint,2,opt,name=Rnd,proto3" json:"Rnd,omitempty"`
	Val  *Value `protobuf:"bytes,3,opt,name=Val,proto3" json:"Val,omitempty"`
}

func (x *LearnMsg) Reset() {
	*x = LearnMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_multipaxos_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LearnMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LearnMsg) ProtoMessage() {}

func (x *LearnMsg) ProtoReflect() protoreflect.Message {
	mi := &file_multipaxos_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LearnMsg.ProtoReflect.Descriptor instead.
func (*LearnMsg) Descriptor() ([]byte, []int) {
	return file_multipaxos_proto_rawDescGZIP(), []int{5}
}

func (x *LearnMsg) GetSlot() uint32 {
	if x != nil {
		return x.Slot
	}
	return 0
}

func (x *LearnMsg) GetRnd() int32 {
	if x != nil {
		return x.Rnd
	}
	return 0
}

func (x *LearnMsg) GetVal() *Value {
	if x != nil {
		return x.Val
	}
	return nil
}

type PValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Slot uint32 `protobuf:"varint,1,opt,name=Slot,proto3" json:"Slot,omitempty"`
	Vrnd int32  `protobuf:"varint,2,opt,name=Vrnd,proto3" json:"Vrnd,omitempty"`
	Vval *Value `protobuf:"bytes,3,opt,name=Vval,proto3" json:"Vval,omitempty"`
}

func (x *PValue) Reset() {
	*x = PValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_multipaxos_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PValue) ProtoMessage() {}

func (x *PValue) ProtoReflect() protoreflect.Message {
	mi := &file_multipaxos_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PValue.ProtoReflect.Descriptor instead.
func (*PValue) Descriptor() ([]byte, []int) {
	return file_multipaxos_proto_rawDescGZIP(), []int{6}
}

func (x *PValue) GetSlot() uint32 {
	if x != nil {
		return x.Slot
	}
	return 0
}

func (x *PValue) GetVrnd() int32 {
	if x != nil {
		return x.Vrnd
	}
	return 0
}

func (x *PValue) GetVval() *Value {
	if x != nil {
		return x.Vval
	}
	return nil
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_multipaxos_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_multipaxos_proto_msgTypes[7]
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
	return file_multipaxos_proto_rawDescGZIP(), []int{7}
}

var File_multipaxos_proto protoreflect.FileDescriptor

var file_multipaxos_proto_rawDesc = []byte{
	0x0a, 0x10, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x61, 0x78, 0x6f, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x67, 0x6f, 0x72, 0x75, 0x6d,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7f, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x12, 0x1a, 0x0a, 0x08, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09,
	0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x71, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x09, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x69, 0x73,
	0x4e, 0x6f, 0x6f, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x69, 0x73, 0x4e, 0x6f,
	0x6f, 0x70, 0x12, 0x24, 0x0a, 0x0d, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x43, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x22, 0x6a, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x44,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x44,
	0x12, 0x1c, 0x0a, 0x09, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x71, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x09, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x71, 0x12, 0x24,
	0x0a, 0x0d, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x22, 0x34, 0x0a, 0x0a, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x4d,
	0x73, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x53, 0x6c, 0x6f, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x04, 0x53, 0x6c, 0x6f, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x43, 0x72, 0x6e, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x43, 0x72, 0x6e, 0x64, 0x22, 0x49, 0x0a, 0x0a, 0x50, 0x72,
	0x6f, 0x6d, 0x69, 0x73, 0x65, 0x4d, 0x73, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x52, 0x6e, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x52, 0x6e, 0x64, 0x12, 0x29, 0x0a, 0x08, 0x41, 0x63,
	0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x08, 0x41, 0x63, 0x63,
	0x65, 0x70, 0x74, 0x65, 0x64, 0x22, 0x51, 0x0a, 0x09, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x4d,
	0x73, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x53, 0x6c, 0x6f, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x04, 0x53, 0x6c, 0x6f, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x52, 0x6e, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x03, 0x52, 0x6e, 0x64, 0x12, 0x1e, 0x0a, 0x03, 0x56, 0x61, 0x6c, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x52, 0x03, 0x56, 0x61, 0x6c, 0x22, 0x50, 0x0a, 0x08, 0x4c, 0x65, 0x61, 0x72,
	0x6e, 0x4d, 0x73, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x53, 0x6c, 0x6f, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x04, 0x53, 0x6c, 0x6f, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x52, 0x6e, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x52, 0x6e, 0x64, 0x12, 0x1e, 0x0a, 0x03, 0x56, 0x61,
	0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x03, 0x56, 0x61, 0x6c, 0x22, 0x52, 0x0a, 0x06, 0x50, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x53, 0x6c, 0x6f, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x04, 0x53, 0x6c, 0x6f, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x56, 0x72, 0x6e, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x56, 0x72, 0x6e, 0x64, 0x12, 0x20, 0x0a, 0x04,
	0x56, 0x76, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x04, 0x56, 0x76, 0x61, 0x6c, 0x22, 0x07,
	0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0xda, 0x01, 0x0a, 0x0a, 0x4d, 0x75, 0x6c, 0x74,
	0x69, 0x50, 0x61, 0x78, 0x6f, 0x73, 0x12, 0x35, 0x0a, 0x07, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72,
	0x65, 0x12, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72,
	0x65, 0x4d, 0x73, 0x67, 0x1a, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x72, 0x6f,
	0x6d, 0x69, 0x73, 0x65, 0x4d, 0x73, 0x67, 0x22, 0x04, 0xa0, 0xb5, 0x18, 0x01, 0x12, 0x31, 0x0a,
	0x06, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x12, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x4d, 0x73, 0x67, 0x1a, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x4c, 0x65, 0x61, 0x72, 0x6e, 0x4d, 0x73, 0x67, 0x22, 0x04, 0xa0, 0xb5, 0x18, 0x01,
	0x12, 0x2d, 0x0a, 0x06, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x0f, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x4c, 0x65, 0x61, 0x72, 0x6e, 0x4d, 0x73, 0x67, 0x1a, 0x0c, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x04, 0x98, 0xb5, 0x18, 0x01, 0x12,
	0x33, 0x0a, 0x0c, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x12,
	0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x0f, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x04,
	0xa0, 0xb5, 0x18, 0x01, 0x42, 0x10, 0x5a, 0x0e, 0x70, 0x61, 0x78, 0x6f, 0x73, 0x2e, 0x71, 0x63,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_multipaxos_proto_rawDescOnce sync.Once
	file_multipaxos_proto_rawDescData = file_multipaxos_proto_rawDesc
)

func file_multipaxos_proto_rawDescGZIP() []byte {
	file_multipaxos_proto_rawDescOnce.Do(func() {
		file_multipaxos_proto_rawDescData = protoimpl.X.CompressGZIP(file_multipaxos_proto_rawDescData)
	})
	return file_multipaxos_proto_rawDescData
}

var file_multipaxos_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_multipaxos_proto_goTypes = []interface{}{
	(*Value)(nil),      // 0: proto.Value
	(*Response)(nil),   // 1: proto.Response
	(*PrepareMsg)(nil), // 2: proto.PrepareMsg
	(*PromiseMsg)(nil), // 3: proto.PromiseMsg
	(*AcceptMsg)(nil),  // 4: proto.AcceptMsg
	(*LearnMsg)(nil),   // 5: proto.LearnMsg
	(*PValue)(nil),     // 6: proto.PValue
	(*Empty)(nil),      // 7: proto.Empty
}
var file_multipaxos_proto_depIdxs = []int32{
	6, // 0: proto.PromiseMsg.Accepted:type_name -> proto.PValue
	0, // 1: proto.AcceptMsg.Val:type_name -> proto.Value
	0, // 2: proto.LearnMsg.Val:type_name -> proto.Value
	0, // 3: proto.PValue.Vval:type_name -> proto.Value
	2, // 4: proto.MultiPaxos.Prepare:input_type -> proto.PrepareMsg
	4, // 5: proto.MultiPaxos.Accept:input_type -> proto.AcceptMsg
	5, // 6: proto.MultiPaxos.Commit:input_type -> proto.LearnMsg
	0, // 7: proto.MultiPaxos.ClientHandle:input_type -> proto.Value
	3, // 8: proto.MultiPaxos.Prepare:output_type -> proto.PromiseMsg
	5, // 9: proto.MultiPaxos.Accept:output_type -> proto.LearnMsg
	7, // 10: proto.MultiPaxos.Commit:output_type -> proto.Empty
	1, // 11: proto.MultiPaxos.ClientHandle:output_type -> proto.Response
	8, // [8:12] is the sub-list for method output_type
	4, // [4:8] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_multipaxos_proto_init() }
func file_multipaxos_proto_init() {
	if File_multipaxos_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_multipaxos_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Value); i {
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
		file_multipaxos_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_multipaxos_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PrepareMsg); i {
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
		file_multipaxos_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PromiseMsg); i {
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
		file_multipaxos_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AcceptMsg); i {
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
		file_multipaxos_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LearnMsg); i {
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
		file_multipaxos_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PValue); i {
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
		file_multipaxos_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
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
			RawDescriptor: file_multipaxos_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_multipaxos_proto_goTypes,
		DependencyIndexes: file_multipaxos_proto_depIdxs,
		MessageInfos:      file_multipaxos_proto_msgTypes,
	}.Build()
	File_multipaxos_proto = out.File
	file_multipaxos_proto_rawDesc = nil
	file_multipaxos_proto_goTypes = nil
	file_multipaxos_proto_depIdxs = nil
}
