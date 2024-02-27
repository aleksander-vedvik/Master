package client

import (
	"fmt"

	pb "github.com/aleksander-vedvik/Master/protos"
)

type RepySpec struct {
	qsize int
}

func NewReplySpec(qsize int) pb.ReplySpec {
	return &RepySpec{
		qsize: qsize,
	}
}

func (rs *RepySpec) SaveStudent(reqs []*pb.ClientResponse) (*pb.ClientResponse, error) {
	//fmt.Println("CLIENT SERVER:", len(reqs))
	if len(reqs) < rs.qsize {
		return nil, fmt.Errorf("not a quorum")
	}
	return reqs[0], nil
}

func (rs *RepySpec) SaveStudents(reqs []*pb.ClientResponse) (*pb.ClientResponse, error) {
	if len(reqs) < rs.qsize {
		return nil, fmt.Errorf("not a quorum")
	}
	return reqs[0], nil
}
