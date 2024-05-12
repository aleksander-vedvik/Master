// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v3.12.4
// source: node.proto

package protos

import (
	empty "github.com/golang/protobuf/ptypes/empty"
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

type WriteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	From      string `protobuf:"bytes,2,opt,name=From,proto3" json:"From,omitempty"`
	Message   string `protobuf:"bytes,3,opt,name=Message,proto3" json:"Message,omitempty"`
	Timestamp int64  `protobuf:"varint,4,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
}

func (x *WriteRequest) Reset() {
	*x = WriteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WriteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WriteRequest) ProtoMessage() {}

func (x *WriteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_node_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WriteRequest.ProtoReflect.Descriptor instead.
func (*WriteRequest) Descriptor() ([]byte, []int) {
	return file_node_proto_rawDescGZIP(), []int{0}
}

func (x *WriteRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *WriteRequest) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *WriteRequest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *WriteRequest) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

type PrePrepareRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	View           int32  `protobuf:"varint,2,opt,name=View,proto3" json:"View,omitempty"`
	SequenceNumber int32  `protobuf:"varint,3,opt,name=SequenceNumber,proto3" json:"SequenceNumber,omitempty"`
	Digest         string `protobuf:"bytes,4,opt,name=Digest,proto3" json:"Digest,omitempty"`
	Timestamp      int64  `protobuf:"varint,5,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	Message        string `protobuf:"bytes,6,opt,name=Message,proto3" json:"Message,omitempty"`
}

func (x *PrePrepareRequest) Reset() {
	*x = PrePrepareRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrePrepareRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrePrepareRequest) ProtoMessage() {}

func (x *PrePrepareRequest) ProtoReflect() protoreflect.Message {
	mi := &file_node_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrePrepareRequest.ProtoReflect.Descriptor instead.
func (*PrePrepareRequest) Descriptor() ([]byte, []int) {
	return file_node_proto_rawDescGZIP(), []int{1}
}

func (x *PrePrepareRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *PrePrepareRequest) GetView() int32 {
	if x != nil {
		return x.View
	}
	return 0
}

func (x *PrePrepareRequest) GetSequenceNumber() int32 {
	if x != nil {
		return x.SequenceNumber
	}
	return 0
}

func (x *PrePrepareRequest) GetDigest() string {
	if x != nil {
		return x.Digest
	}
	return ""
}

func (x *PrePrepareRequest) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *PrePrepareRequest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type PrepareRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	View           int32  `protobuf:"varint,2,opt,name=View,proto3" json:"View,omitempty"`
	SequenceNumber int32  `protobuf:"varint,3,opt,name=SequenceNumber,proto3" json:"SequenceNumber,omitempty"`
	Digest         string `protobuf:"bytes,4,opt,name=Digest,proto3" json:"Digest,omitempty"`
	Timestamp      int64  `protobuf:"varint,5,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	Message        string `protobuf:"bytes,6,opt,name=Message,proto3" json:"Message,omitempty"`
}

func (x *PrepareRequest) Reset() {
	*x = PrepareRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrepareRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrepareRequest) ProtoMessage() {}

func (x *PrepareRequest) ProtoReflect() protoreflect.Message {
	mi := &file_node_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrepareRequest.ProtoReflect.Descriptor instead.
func (*PrepareRequest) Descriptor() ([]byte, []int) {
	return file_node_proto_rawDescGZIP(), []int{2}
}

func (x *PrepareRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *PrepareRequest) GetView() int32 {
	if x != nil {
		return x.View
	}
	return 0
}

func (x *PrepareRequest) GetSequenceNumber() int32 {
	if x != nil {
		return x.SequenceNumber
	}
	return 0
}

func (x *PrepareRequest) GetDigest() string {
	if x != nil {
		return x.Digest
	}
	return ""
}

func (x *PrepareRequest) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *PrepareRequest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type CommitRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	View           int32  `protobuf:"varint,2,opt,name=View,proto3" json:"View,omitempty"`
	SequenceNumber int32  `protobuf:"varint,3,opt,name=SequenceNumber,proto3" json:"SequenceNumber,omitempty"`
	Digest         string `protobuf:"bytes,4,opt,name=Digest,proto3" json:"Digest,omitempty"`
	Timestamp      int64  `protobuf:"varint,5,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	Message        string `protobuf:"bytes,6,opt,name=Message,proto3" json:"Message,omitempty"`
}

func (x *CommitRequest) Reset() {
	*x = CommitRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommitRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommitRequest) ProtoMessage() {}

func (x *CommitRequest) ProtoReflect() protoreflect.Message {
	mi := &file_node_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommitRequest.ProtoReflect.Descriptor instead.
func (*CommitRequest) Descriptor() ([]byte, []int) {
	return file_node_proto_rawDescGZIP(), []int{3}
}

func (x *CommitRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CommitRequest) GetView() int32 {
	if x != nil {
		return x.View
	}
	return 0
}

func (x *CommitRequest) GetSequenceNumber() int32 {
	if x != nil {
		return x.SequenceNumber
	}
	return 0
}

func (x *CommitRequest) GetDigest() string {
	if x != nil {
		return x.Digest
	}
	return ""
}

func (x *CommitRequest) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *CommitRequest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type ClientResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	View      int32  `protobuf:"varint,2,opt,name=View,proto3" json:"View,omitempty"`
	Timestamp int64  `protobuf:"varint,3,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	Client    string `protobuf:"bytes,4,opt,name=Client,proto3" json:"Client,omitempty"`
	From      string `protobuf:"bytes,5,opt,name=From,proto3" json:"From,omitempty"`
	Result    string `protobuf:"bytes,6,opt,name=Result,proto3" json:"Result,omitempty"`
}

func (x *ClientResponse) Reset() {
	*x = ClientResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientResponse) ProtoMessage() {}

func (x *ClientResponse) ProtoReflect() protoreflect.Message {
	mi := &file_node_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientResponse.ProtoReflect.Descriptor instead.
func (*ClientResponse) Descriptor() ([]byte, []int) {
	return file_node_proto_rawDescGZIP(), []int{4}
}

func (x *ClientResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ClientResponse) GetView() int32 {
	if x != nil {
		return x.View
	}
	return 0
}

func (x *ClientResponse) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *ClientResponse) GetClient() string {
	if x != nil {
		return x.Client
	}
	return ""
}

func (x *ClientResponse) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *ClientResponse) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

type Heartbeat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint32 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (x *Heartbeat) Reset() {
	*x = Heartbeat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Heartbeat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Heartbeat) ProtoMessage() {}

func (x *Heartbeat) ProtoReflect() protoreflect.Message {
	mi := &file_node_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Heartbeat.ProtoReflect.Descriptor instead.
func (*Heartbeat) Descriptor() ([]byte, []int) {
	return file_node_proto_rawDescGZIP(), []int{5}
}

func (x *Heartbeat) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type Result struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics []*Metric `protobuf:"bytes,1,rep,name=metrics,proto3" json:"metrics,omitempty"`
}

func (x *Result) Reset() {
	*x = Result{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Result) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Result) ProtoMessage() {}

func (x *Result) ProtoReflect() protoreflect.Message {
	mi := &file_node_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Result.ProtoReflect.Descriptor instead.
func (*Result) Descriptor() ([]byte, []int) {
	return file_node_proto_rawDescGZIP(), []int{6}
}

func (x *Result) GetMetrics() []*Metric {
	if x != nil {
		return x.Metrics
	}
	return nil
}

type Metric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TotalNum              uint64                      `protobuf:"varint,1,opt,name=TotalNum,proto3" json:"TotalNum,omitempty"`
	GoroutinesStarted     uint64                      `protobuf:"varint,2,opt,name=GoroutinesStarted,proto3" json:"GoroutinesStarted,omitempty"`
	GoroutinesStopped     uint64                      `protobuf:"varint,3,opt,name=GoroutinesStopped,proto3" json:"GoroutinesStopped,omitempty"`
	Goroutines            map[string]*GoroutineMetric `protobuf:"bytes,4,rep,name=Goroutines,proto3" json:"Goroutines,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	FinishedReqsTotal     uint64                      `protobuf:"varint,5,opt,name=FinishedReqsTotal,proto3" json:"FinishedReqsTotal,omitempty"`
	FinishedReqsSuccesful uint64                      `protobuf:"varint,6,opt,name=FinishedReqsSuccesful,proto3" json:"FinishedReqsSuccesful,omitempty"`
	FinishedReqsFailed    uint64                      `protobuf:"varint,7,opt,name=FinishedReqsFailed,proto3" json:"FinishedReqsFailed,omitempty"`
	Processed             uint64                      `protobuf:"varint,8,opt,name=Processed,proto3" json:"Processed,omitempty"`
	Dropped               uint64                      `protobuf:"varint,9,opt,name=Dropped,proto3" json:"Dropped,omitempty"`
	Invalid               uint64                      `protobuf:"varint,10,opt,name=Invalid,proto3" json:"Invalid,omitempty"`
	AlreadyProcessed      uint64                      `protobuf:"varint,11,opt,name=AlreadyProcessed,proto3" json:"AlreadyProcessed,omitempty"`
	RoundTripLatency      *TimingMetric               `protobuf:"bytes,12,opt,name=RoundTripLatency,proto3" json:"RoundTripLatency,omitempty"`
	ReqLatency            *TimingMetric               `protobuf:"bytes,13,opt,name=ReqLatency,proto3" json:"ReqLatency,omitempty"`
	ShardDistribution     map[uint32]uint64           `protobuf:"bytes,14,rep,name=ShardDistribution,proto3" json:"ShardDistribution,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *Metric) Reset() {
	*x = Metric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_node_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metric.ProtoReflect.Descriptor instead.
func (*Metric) Descriptor() ([]byte, []int) {
	return file_node_proto_rawDescGZIP(), []int{7}
}

func (x *Metric) GetTotalNum() uint64 {
	if x != nil {
		return x.TotalNum
	}
	return 0
}

func (x *Metric) GetGoroutinesStarted() uint64 {
	if x != nil {
		return x.GoroutinesStarted
	}
	return 0
}

func (x *Metric) GetGoroutinesStopped() uint64 {
	if x != nil {
		return x.GoroutinesStopped
	}
	return 0
}

func (x *Metric) GetGoroutines() map[string]*GoroutineMetric {
	if x != nil {
		return x.Goroutines
	}
	return nil
}

func (x *Metric) GetFinishedReqsTotal() uint64 {
	if x != nil {
		return x.FinishedReqsTotal
	}
	return 0
}

func (x *Metric) GetFinishedReqsSuccesful() uint64 {
	if x != nil {
		return x.FinishedReqsSuccesful
	}
	return 0
}

func (x *Metric) GetFinishedReqsFailed() uint64 {
	if x != nil {
		return x.FinishedReqsFailed
	}
	return 0
}

func (x *Metric) GetProcessed() uint64 {
	if x != nil {
		return x.Processed
	}
	return 0
}

func (x *Metric) GetDropped() uint64 {
	if x != nil {
		return x.Dropped
	}
	return 0
}

func (x *Metric) GetInvalid() uint64 {
	if x != nil {
		return x.Invalid
	}
	return 0
}

func (x *Metric) GetAlreadyProcessed() uint64 {
	if x != nil {
		return x.AlreadyProcessed
	}
	return 0
}

func (x *Metric) GetRoundTripLatency() *TimingMetric {
	if x != nil {
		return x.RoundTripLatency
	}
	return nil
}

func (x *Metric) GetReqLatency() *TimingMetric {
	if x != nil {
		return x.ReqLatency
	}
	return nil
}

func (x *Metric) GetShardDistribution() map[uint32]uint64 {
	if x != nil {
		return x.ShardDistribution
	}
	return nil
}

type GoroutineMetric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Start uint64 `protobuf:"varint,1,opt,name=start,proto3" json:"start,omitempty"`
	End   uint64 `protobuf:"varint,2,opt,name=end,proto3" json:"end,omitempty"`
}

func (x *GoroutineMetric) Reset() {
	*x = GoroutineMetric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GoroutineMetric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GoroutineMetric) ProtoMessage() {}

func (x *GoroutineMetric) ProtoReflect() protoreflect.Message {
	mi := &file_node_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GoroutineMetric.ProtoReflect.Descriptor instead.
func (*GoroutineMetric) Descriptor() ([]byte, []int) {
	return file_node_proto_rawDescGZIP(), []int{8}
}

func (x *GoroutineMetric) GetStart() uint64 {
	if x != nil {
		return x.Start
	}
	return 0
}

func (x *GoroutineMetric) GetEnd() uint64 {
	if x != nil {
		return x.End
	}
	return 0
}

type TimingMetric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Avg uint64 `protobuf:"varint,1,opt,name=avg,proto3" json:"avg,omitempty"`
	Min uint64 `protobuf:"varint,2,opt,name=min,proto3" json:"min,omitempty"`
	Max uint64 `protobuf:"varint,3,opt,name=max,proto3" json:"max,omitempty"`
}

func (x *TimingMetric) Reset() {
	*x = TimingMetric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimingMetric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimingMetric) ProtoMessage() {}

func (x *TimingMetric) ProtoReflect() protoreflect.Message {
	mi := &file_node_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimingMetric.ProtoReflect.Descriptor instead.
func (*TimingMetric) Descriptor() ([]byte, []int) {
	return file_node_proto_rawDescGZIP(), []int{9}
}

func (x *TimingMetric) GetAvg() uint64 {
	if x != nil {
		return x.Avg
	}
	return 0
}

func (x *TimingMetric) GetMin() uint64 {
	if x != nil {
		return x.Min
	}
	return 0
}

func (x *TimingMetric) GetMax() uint64 {
	if x != nil {
		return x.Max
	}
	return 0
}

var File_node_proto protoreflect.FileDescriptor

var file_node_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x73, 0x50, 0x42, 0x46, 0x54, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6a, 0x0a, 0x0c, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x22, 0xaf, 0x01, 0x0a, 0x11, 0x50, 0x72, 0x65, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x56, 0x69, 0x65, 0x77, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x56, 0x69, 0x65, 0x77, 0x12, 0x26, 0x0a, 0x0e, 0x53,
	0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0e, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x22, 0xac, 0x01, 0x0a, 0x0e, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x56, 0x69, 0x65, 0x77, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x56, 0x69, 0x65, 0x77, 0x12, 0x26, 0x0a, 0x0e, 0x53, 0x65,
	0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0e, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x4e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x22, 0xab, 0x01, 0x0a, 0x0d, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x56, 0x69, 0x65, 0x77, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x04, 0x56, 0x69, 0x65, 0x77, 0x12, 0x26, 0x0a, 0x0e, 0x53, 0x65, 0x71, 0x75,
	0x65, 0x6e, 0x63, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0e, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x12, 0x16, 0x0a, 0x06, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x22, 0x96, 0x01, 0x0a, 0x0e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x56, 0x69, 0x65, 0x77, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x04, 0x56, 0x69, 0x65, 0x77, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x16, 0x0a, 0x06, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x46, 0x72, 0x6f, 0x6d, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x46, 0x72, 0x6f,
	0x6d, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x1b, 0x0a, 0x09, 0x48, 0x65, 0x61,
	0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x02, 0x49, 0x64, 0x22, 0x36, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x12, 0x2c, 0x0a, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x50, 0x42, 0x46, 0x54, 0x2e, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x22, 0xd1,
	0x06, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x1a, 0x0a, 0x08, 0x54, 0x6f, 0x74,
	0x61, 0x6c, 0x4e, 0x75, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x54, 0x6f, 0x74,
	0x61, 0x6c, 0x4e, 0x75, 0x6d, 0x12, 0x2c, 0x0a, 0x11, 0x47, 0x6f, 0x72, 0x6f, 0x75, 0x74, 0x69,
	0x6e, 0x65, 0x73, 0x53, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x11, 0x47, 0x6f, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x65, 0x73, 0x53, 0x74, 0x61, 0x72,
	0x74, 0x65, 0x64, 0x12, 0x2c, 0x0a, 0x11, 0x47, 0x6f, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x65,
	0x73, 0x53, 0x74, 0x6f, 0x70, 0x70, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x11,
	0x47, 0x6f, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x65, 0x73, 0x53, 0x74, 0x6f, 0x70, 0x70, 0x65,
	0x64, 0x12, 0x42, 0x0a, 0x0a, 0x47, 0x6f, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x65, 0x73, 0x18,
	0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x50, 0x42,
	0x46, 0x54, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x47, 0x6f, 0x72, 0x6f, 0x75, 0x74,
	0x69, 0x6e, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x47, 0x6f, 0x72, 0x6f, 0x75,
	0x74, 0x69, 0x6e, 0x65, 0x73, 0x12, 0x2c, 0x0a, 0x11, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65,
	0x64, 0x52, 0x65, 0x71, 0x73, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x11, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x52, 0x65, 0x71, 0x73, 0x54, 0x6f,
	0x74, 0x61, 0x6c, 0x12, 0x34, 0x0a, 0x15, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x52,
	0x65, 0x71, 0x73, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x66, 0x75, 0x6c, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x15, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x52, 0x65, 0x71, 0x73,
	0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x66, 0x75, 0x6c, 0x12, 0x2e, 0x0a, 0x12, 0x46, 0x69, 0x6e,
	0x69, 0x73, 0x68, 0x65, 0x64, 0x52, 0x65, 0x71, 0x73, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x12, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x52,
	0x65, 0x71, 0x73, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x50, 0x72, 0x6f,
	0x63, 0x65, 0x73, 0x73, 0x65, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x50, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x44, 0x72, 0x6f, 0x70, 0x70,
	0x65, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x44, 0x72, 0x6f, 0x70, 0x70, 0x65,
	0x64, 0x12, 0x18, 0x0a, 0x07, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x07, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x12, 0x2a, 0x0a, 0x10, 0x41,
	0x6c, 0x72, 0x65, 0x61, 0x64, 0x79, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x64, 0x18,
	0x0b, 0x20, 0x01, 0x28, 0x04, 0x52, 0x10, 0x41, 0x6c, 0x72, 0x65, 0x61, 0x64, 0x79, 0x50, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x64, 0x12, 0x44, 0x0a, 0x10, 0x52, 0x6f, 0x75, 0x6e, 0x64,
	0x54, 0x72, 0x69, 0x70, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x0c, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x50, 0x42, 0x46, 0x54, 0x2e, 0x54,
	0x69, 0x6d, 0x69, 0x6e, 0x67, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x10, 0x52, 0x6f, 0x75,
	0x6e, 0x64, 0x54, 0x72, 0x69, 0x70, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x12, 0x38, 0x0a,
	0x0a, 0x52, 0x65, 0x71, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x0d, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x50, 0x42, 0x46, 0x54, 0x2e, 0x54,
	0x69, 0x6d, 0x69, 0x6e, 0x67, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x0a, 0x52, 0x65, 0x71,
	0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x12, 0x57, 0x0a, 0x11, 0x53, 0x68, 0x61, 0x72, 0x64,
	0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0e, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x29, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x50, 0x42, 0x46, 0x54, 0x2e,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x64, 0x44, 0x69, 0x73, 0x74,
	0x72, 0x69, 0x62, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x11, 0x53,
	0x68, 0x61, 0x72, 0x64, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x69, 0x6f, 0x6e,
	0x1a, 0x5a, 0x0a, 0x0f, 0x47, 0x6f, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x65, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x31, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x50, 0x42, 0x46,
	0x54, 0x2e, 0x47, 0x6f, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x44, 0x0a, 0x16,
	0x53, 0x68, 0x61, 0x72, 0x64, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x69, 0x6f,
	0x6e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0x39, 0x0a, 0x0f, 0x47, 0x6f, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x65, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x65,
	0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x65, 0x6e, 0x64, 0x22, 0x44, 0x0a,
	0x0c, 0x54, 0x69, 0x6d, 0x69, 0x6e, 0x67, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x10, 0x0a,
	0x03, 0x61, 0x76, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x61, 0x76, 0x67, 0x12,
	0x10, 0x0a, 0x03, 0x6d, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x6d, 0x69,
	0x6e, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x61, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03,
	0x6d, 0x61, 0x78, 0x32, 0xcd, 0x02, 0x0a, 0x08, 0x50, 0x42, 0x46, 0x54, 0x4e, 0x6f, 0x64, 0x65,
	0x12, 0x3f, 0x0a, 0x05, 0x57, 0x72, 0x69, 0x74, 0x65, 0x12, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x73, 0x50, 0x42, 0x46, 0x54, 0x2e, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x50, 0x42, 0x46, 0x54,
	0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x45, 0x0a, 0x0a, 0x50, 0x72, 0x65, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x12,
	0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x50, 0x42, 0x46, 0x54, 0x2e, 0x50, 0x72, 0x65,
	0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x3f, 0x0a, 0x07, 0x50, 0x72, 0x65, 0x70,
	0x61, 0x72, 0x65, 0x12, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x50, 0x42, 0x46, 0x54,
	0x2e, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x3d, 0x0a, 0x06, 0x43, 0x6f, 0x6d,
	0x6d, 0x69, 0x74, 0x12, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x50, 0x42, 0x46, 0x54,
	0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x39, 0x0a, 0x09, 0x42, 0x65, 0x6e, 0x63,
	0x68, 0x6d, 0x61, 0x72, 0x6b, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x12, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x50, 0x42, 0x46, 0x54, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x22, 0x00, 0x42, 0x0d, 0x5a, 0x0b, 0x70, 0x62, 0x66, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_node_proto_rawDescOnce sync.Once
	file_node_proto_rawDescData = file_node_proto_rawDesc
)

func file_node_proto_rawDescGZIP() []byte {
	file_node_proto_rawDescOnce.Do(func() {
		file_node_proto_rawDescData = protoimpl.X.CompressGZIP(file_node_proto_rawDescData)
	})
	return file_node_proto_rawDescData
}

var file_node_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_node_proto_goTypes = []interface{}{
	(*WriteRequest)(nil),      // 0: protosPBFT.WriteRequest
	(*PrePrepareRequest)(nil), // 1: protosPBFT.PrePrepareRequest
	(*PrepareRequest)(nil),    // 2: protosPBFT.PrepareRequest
	(*CommitRequest)(nil),     // 3: protosPBFT.CommitRequest
	(*ClientResponse)(nil),    // 4: protosPBFT.ClientResponse
	(*Heartbeat)(nil),         // 5: protosPBFT.Heartbeat
	(*Result)(nil),            // 6: protosPBFT.Result
	(*Metric)(nil),            // 7: protosPBFT.Metric
	(*GoroutineMetric)(nil),   // 8: protosPBFT.GoroutineMetric
	(*TimingMetric)(nil),      // 9: protosPBFT.TimingMetric
	nil,                       // 10: protosPBFT.Metric.GoroutinesEntry
	nil,                       // 11: protosPBFT.Metric.ShardDistributionEntry
	(*empty.Empty)(nil),       // 12: google.protobuf.Empty
}
var file_node_proto_depIdxs = []int32{
	7,  // 0: protosPBFT.Result.metrics:type_name -> protosPBFT.Metric
	10, // 1: protosPBFT.Metric.Goroutines:type_name -> protosPBFT.Metric.GoroutinesEntry
	9,  // 2: protosPBFT.Metric.RoundTripLatency:type_name -> protosPBFT.TimingMetric
	9,  // 3: protosPBFT.Metric.ReqLatency:type_name -> protosPBFT.TimingMetric
	11, // 4: protosPBFT.Metric.ShardDistribution:type_name -> protosPBFT.Metric.ShardDistributionEntry
	8,  // 5: protosPBFT.Metric.GoroutinesEntry.value:type_name -> protosPBFT.GoroutineMetric
	0,  // 6: protosPBFT.PBFTNode.Write:input_type -> protosPBFT.WriteRequest
	1,  // 7: protosPBFT.PBFTNode.PrePrepare:input_type -> protosPBFT.PrePrepareRequest
	2,  // 8: protosPBFT.PBFTNode.Prepare:input_type -> protosPBFT.PrepareRequest
	3,  // 9: protosPBFT.PBFTNode.Commit:input_type -> protosPBFT.CommitRequest
	12, // 10: protosPBFT.PBFTNode.Benchmark:input_type -> google.protobuf.Empty
	4,  // 11: protosPBFT.PBFTNode.Write:output_type -> protosPBFT.ClientResponse
	12, // 12: protosPBFT.PBFTNode.PrePrepare:output_type -> google.protobuf.Empty
	12, // 13: protosPBFT.PBFTNode.Prepare:output_type -> google.protobuf.Empty
	12, // 14: protosPBFT.PBFTNode.Commit:output_type -> google.protobuf.Empty
	6,  // 15: protosPBFT.PBFTNode.Benchmark:output_type -> protosPBFT.Result
	11, // [11:16] is the sub-list for method output_type
	6,  // [6:11] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_node_proto_init() }
func file_node_proto_init() {
	if File_node_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_node_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WriteRequest); i {
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
		file_node_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PrePrepareRequest); i {
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
		file_node_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PrepareRequest); i {
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
		file_node_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommitRequest); i {
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
		file_node_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientResponse); i {
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
		file_node_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Heartbeat); i {
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
		file_node_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Result); i {
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
		file_node_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metric); i {
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
		file_node_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GoroutineMetric); i {
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
		file_node_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimingMetric); i {
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
			RawDescriptor: file_node_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_node_proto_goTypes,
		DependencyIndexes: file_node_proto_depIdxs,
		MessageInfos:      file_node_proto_msgTypes,
	}.Build()
	File_node_proto = out.File
	file_node_proto_rawDesc = nil
	file_node_proto_goTypes = nil
	file_node_proto_depIdxs = nil
}
