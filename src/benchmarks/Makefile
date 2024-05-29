.PHONY: proto srv1 srv2 srv3 client run
srv1:
	go run main.go -server -id=0 -run=$(bench) -log

srv2:
	go run main.go -server -id=1 -run=$(bench) -log

srv3:
	go run main.go -server -id=2 -run=$(bench) -log

client:
	go run main.go -run=$(bench) -local=false -throughput=15000 -steps=1

run: srv1 srv2 srv3

all: pbft/protos/node.pb.go pbft/protos/node_gorums.pb.go

%.pb.go %_gorums.pb.go : %.proto
	@protoc -I=$(go list -m -f {{.Dir}} github.com/relab/gorums):. \
		--go_out=paths=source_relative:. \
		--gorums_out=paths=source_relative:. \
		$<
