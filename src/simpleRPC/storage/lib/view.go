package lib

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/aleksander-vedvik/Master/protos"
)

type View struct {
	servers []string
}

func NewView(srvAddrs []string) *View {
	return &View{
		servers: srvAddrs,
	}
}

func (q *View) dial(addr string) (pb.StorageClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	return pb.NewStorageClient(conn), conn
}

func (q *View) Write(val string, ids ...int64) bool {
	var id int64
	if len(ids) <= 0 {
		id = int64(rand.Intn(1000))
	} else {
		id = ids[0]
	}
	timestamp := time.Now().Unix()
	respChan := make(chan bool, len(q.servers))
	for _, addr := range q.servers {
		go func(addr string) {
			node, conn := q.dial(addr)
			if node == nil {
				if conn != nil {
					conn.Close()
				}
				respChan <- false
				return
			}
			defer conn.Close()
			_, err := node.Write(context.Background(), &pb.State{
				Id:        id,
				Value:     val,
				Timestamp: timestamp,
			})
			if err != nil {
				fmt.Println(err)
				respChan <- false
				return
			}
			respChan <- true
		}(addr)
	}
	successes := 0
	failures := 0
	for success := range respChan {
		if success {
			successes++
		} else {
			failures++
		}
		if successes >= int(float64(1+len(q.servers)/2)) {
			return true
		}
		if failures+successes >= len(q.servers) {
			return false
		}
	}
	return false
}

func (q *View) Read() string {
	respChan := make(chan string, len(q.servers))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, addr := range q.servers {
		go func(addr string) {
			node, conn := q.dial(addr)
			if node == nil {
				if conn != nil {
					conn.Close()
				}
				respChan <- "ERROR"
				return
			}
			defer conn.Close()
			val, err := node.Read(ctx, &pb.ReadRequest{})
			if err != nil {
				respChan <- "ERROR"
				return
			}
			respChan <- val.Value
		}(addr)
	}
	successes := 0
	failures := 0
	for val := range respChan {
		if val != "ERROR" {
			successes++
		} else {
			failures++
		}
		if successes >= int(float64(1+len(q.servers)/2)) {
			return val
		}
		if failures+successes >= len(q.servers) {
			return ""
		}
	}
	return ""
}
