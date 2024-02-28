package client

import (
	pb "github.com/aleksander-vedvik/Master/protos"
)

type ReplySpec struct {
	qsize int
}

func NewReplySpec(qsize int) pb.ReplySpec {
	return &ReplySpec{
		qsize: qsize,
	}
}

func (rs *ReplySpec) SaveStudent(reqs []*pb.ClientResponse) (*pb.ClientResponse, bool) {
	//fmt.Println("CLIENT SERVER:", len(reqs))
	if len(reqs) < rs.qsize {
		return nil, false
	}
	return reqs[0], true
}

func (rs *ReplySpec) SaveStudents(reqs []*pb.ClientResponse) (*pb.ClientResponse, bool) {
	if len(reqs) < rs.qsize {
		return nil, false
	}
	return reqs[0], true
}
