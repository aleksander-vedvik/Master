module github.com/aleksander-vedvik/Master

go 1.22

require (
	github.com/relab/gorums v0.7.0
	google.golang.org/grpc v1.61.1 // v1.43.0
	google.golang.org/protobuf v1.32.0
)

require github.com/google/uuid v1.6.0

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.18.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240213162025-012b6fc9bca9 // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.3.0 // indirect
)

replace github.com/relab/gorums v0.7.0 => ../../../gorums
