module pbft

go 1.20

require (
	github.com/relab/gorums v0.7.0
	google.golang.org/grpc v1.60.1 // v1.43.0
	google.golang.org/protobuf v1.32.0
)

require github.com/google/uuid v1.3.1

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.16.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	golang.org/x/tools v0.6.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231002182017-d307bd883b97 // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.2.0 // indirect
)

replace github.com/relab/gorums v0.7.0 => ../../../gorums
