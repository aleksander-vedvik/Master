package server

import (
	"fmt"
	"sync"

	pb "github.com/aleksander-vedvik/Master/protos"
	"github.com/relab/gorums"
)

type StorageSrv struct {
	mut   sync.Mutex
	State *pb.State
}

func (srv *StorageSrv) Read(ctx gorums.ServerCtx, req *pb.ReadRequest) (resp *pb.State, err error) {
	// any code running before this will be executed in-order
	ctx.Release()
	// after Release() has been called, a new request handler may be started,
	// and thus it is not guaranteed that the replies will be sent back the same order.
	srv.mut.Lock()
	defer srv.mut.Unlock()
	fmt.Println("Got Read()")
	return srv.State, nil
}

func (srv *StorageSrv) Write(_ gorums.ServerCtx, req *pb.State) (resp *pb.WriteResponse, err error) {
	srv.mut.Lock()
	defer srv.mut.Unlock()
	if srv.State.Timestamp < req.Timestamp {
		srv.State = req
		fmt.Println("Got Write(", req.Value, ")")
		return &pb.WriteResponse{New: true}, nil
	}
	return &pb.WriteResponse{New: false}, nil
}
