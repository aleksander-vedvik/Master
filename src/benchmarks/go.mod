module github.com/aleksander-vedvik/benchmark

go 1.22.2

require (
	github.com/golang/protobuf v1.5.4
	github.com/relab/gorums v0.7.0
	google.golang.org/grpc v1.63.2
	google.golang.org/protobuf v1.33.0
)

require (
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.19.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.3.0 // indirect
)

replace github.com/relab/gorums v0.7.0 => ../../../gorums
