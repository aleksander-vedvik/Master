.PHONY: proto
proto: $(proto_go)
	@echo "+ compiling gorums proto files"

%.pb.go %_gorums.pb.go : %.proto
	@protoc -I=$(gorum_include):. \
		--go_out=paths=source_relative:. \
		--gorums_out=paths=source_relative:. \
		$<

srv1:
	go run main.go -server -id=0 -run=$(bench) -log

srv2:
	go run main.go -server -id=1 -run=$(bench) -log

srv3:
	go run main.go -server -id=2 -run=$(bench) -log

client:
	go run main.go -run=$(bench) -local=false -throughput=15000 -steps=1

run: srv1 srv2 srv3