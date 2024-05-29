package server

import (
	"strconv"
	"sync"

	pb "github.com/aleksander-vedvik/benchmark/pbft/protos"
)

type MessageLog struct {
	mut         sync.Mutex
	data        []any
	preprepares map[string]*pb.PrePrepareRequest
	prepares    map[string][]*pb.PrepareRequest
	commits     map[string][]*pb.CommitRequest
}

func newMessageLog() *MessageLog {
	buffer := 1000
	return &MessageLog{
		data:        make([]any, buffer),
		preprepares: make(map[string]*pb.PrePrepareRequest, buffer),
		prepares:    make(map[string][]*pb.PrepareRequest, buffer),
		commits:     make(map[string][]*pb.CommitRequest, buffer),
	}
}

func (m *MessageLog) add(elem any, view, n int32) {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.data = append(m.data, elem)
	id := "view:" + strconv.Itoa(int(view)) + ",n:" + strconv.Itoa(int(n))
	switch req := elem.(type) {
	case *pb.PrePrepareRequest:
		m.preprepares[id] = req
	case *pb.PrepareRequest:
		m.prepares[id] = append(m.prepares[id], req)
	case *pb.CommitRequest:
		m.commits[id] = append(m.commits[id], req)
	}
}

func (m *MessageLog) getPrePrepareReq(n, view int32) (*pb.PrePrepareRequest, bool) {
	m.mut.Lock()
	defer m.mut.Unlock()
	id := "view:" + strconv.Itoa(int(view)) + ",n:" + strconv.Itoa(int(n))
	req, ok := m.preprepares[id]
	if ok && req.View == view && req.SequenceNumber == n {
		return req, true
	}
	return nil, false
}

func (m *MessageLog) getPrepareReqs(digest string, n, view int32) ([]*pb.PrepareRequest, bool) {
	m.mut.Lock()
	defer m.mut.Unlock()
	id := "view:" + strconv.Itoa(int(view)) + ",n:" + strconv.Itoa(int(n))
	res := make([]*pb.PrepareRequest, 0)
	for _, req := range m.prepares[id] {
		if req.Digest == digest && req.View == view && req.SequenceNumber == n {
			res = append(res, req)
		}
	}
	return res, len(res) > 0
}

func (m *MessageLog) getCommitReqs(digest string, n, view int32) ([]*pb.CommitRequest, bool) {
	m.mut.Lock()
	defer m.mut.Unlock()
	id := "view:" + strconv.Itoa(int(view)) + ",n:" + strconv.Itoa(int(n))
	res := make([]*pb.CommitRequest, 0)
	for _, req := range m.commits[id] {
		if req.Digest == digest && req.View == view && req.SequenceNumber == n {
			res = append(res, req)
		}
	}
	return res, len(res) > 0
}

func (m *MessageLog) Clear() {
	m.mut.Lock()
	defer m.mut.Unlock()
	buffer := 1000
	m.data = make([]any, buffer)
	m.preprepares = make(map[string]*pb.PrePrepareRequest, buffer)
	m.prepares = make(map[string][]*pb.PrepareRequest, buffer)
	m.commits = make(map[string][]*pb.CommitRequest, buffer)
}
