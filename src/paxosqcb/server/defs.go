package server

import (
	"context"
	"slices"

	pb "github.com/aleksander-vedvik/benchmark/paxosqcb/proto"
	"github.com/relab/gorums"
)

// myIndex returns the index of myID in the sorted list of IDs from the node map.
func myIndex(myID int, nodeMap map[string]uint32) int {
	ids := Values(nodeMap)
	slices.Sort(ids)
	return slices.Index(ids, uint32(myID))
}

func Keys[K comparable, V any](m map[K]V) []K {
	ks := make([]K, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}

func Values[K comparable, V any](m map[K]V) []V {
	vs := make([]V, 0, len(m))
	for _, v := range m {
		vs = append(vs, v)
	}
	return vs
}

type (
	Round = int32
	Slot  = uint32
)

const (
	NoRound int32  = -1
	NoSlot  uint32 = 0
)

// MultiPaxosConfig defines the RPC calls on the configuration.
// This interface is used for mocking the configuration in unit tests.
type MultiPaxosConfig interface {
	Prepare(ctx context.Context, request *pb.PrepareMsg) (response *pb.PromiseMsg, err error)
	Accept(ctx context.Context, request *pb.AcceptMsg) (response *pb.LearnMsg, err error)
	ClientHandle(ctx context.Context, request *pb.Value) (*pb.Response, error)
}

// Mock configuration used for testing.
// It stores the received value and allows you to specify the returned value for all RPC calls.
type MockConfiguration struct {
	pb.Configuration

	ErrOut error

	PrpIn  *pb.PrepareMsg
	PrmOut *pb.PromiseMsg

	AccIn  *pb.AcceptMsg
	LrnOut *pb.LearnMsg

	LrnIn  *pb.LearnMsg
	EmpOut *pb.Empty

	ValIn   *pb.Value
	RespOut *pb.Response
}

func (mc *MockConfiguration) Prepare(ctx context.Context, request *pb.PrepareMsg) (response *pb.PromiseMsg, err error) {
	mc.PrpIn = request
	return mc.PrmOut, mc.ErrOut
}

func (mc *MockConfiguration) Accept(ctx context.Context, request *pb.AcceptMsg) (response *pb.LearnMsg, err error) {
	mc.AccIn = request
	return mc.LrnOut, mc.ErrOut
}

func (mc *MockConfiguration) Commit(ctx context.Context, request *pb.LearnMsg, opts ...gorums.CallOption) {
	mc.LrnIn = request
}

func (mc *MockConfiguration) ClientHandle(ctx context.Context, request *pb.Value) (response *pb.Response, err error) {
	mc.ValIn = request
	return mc.RespOut, mc.ErrOut
}
