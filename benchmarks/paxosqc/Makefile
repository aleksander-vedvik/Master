gorum_include := $(shell go list -m -f {{.Dir}} github.com/relab/gorums)
proto_src := proto/paxosqc.proto
proto_go := $(proto_src:%.proto=%.pb.go)

all: proto 

.PHONY: proto
proto: $(proto_go)
	@echo "+ compiling gorums proto files" + $(proto_go)

%.pb.go %_gorums.pb.go : %.proto
	@protoc -I=$(gorum_include):. \
		--go_out=paths=source_relative:. \
		--gorums_out=paths=source_relative:. \
		$<

#proto:
#	protoc -I=$(go list -m -f {{.Dir}} github.com/relab/gorums):. \
#	--go_out=paths=source_relative:. \
#	--gorums_out=paths=source_relative:. \
#	*.proto
#
#protos:
#	protoc --go_out=paths=source_relative:. proto/*.proto

docker:
	docker compose up

docker-rebuild:
	docker compose up --build

srv:
	cd cmd/server; go run main.go --type 0 --id 1
	cd cmd/server; go run main.go --type 0 --id 2

client:
	cd cmd/client; go run main.go

#run: server client
#	@echo "Done!" 

tidy:
	go mod tidy

build:
	go build -o ./bin/

run: build
	cd bin; ./run.sh