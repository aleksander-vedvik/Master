.PHONY: proto docker clean eval

network:
	docker network create --gateway 10.0.0.1 --subnet 10.0.0.0/16 MasterLab

docker:
	docker compose up --build

eval:
	docker compose down
	docker compose up

gorums_include := $(shell go list -m -f {{.Dir}} github.com/relab/gorums)
paxosBC_src := paxos.b/proto/multipaxos.pb.go
paxosQC_src := paxosqc/proto/paxosqc.pb.go
paxosQCB_src := paxosqcb/proto/paxosqcb.pb.go
pbftBC_src := pbft/protos/node.pb.go
pbftBCO_src := pbft.o/protos/nodeo.pb.go

proto: clean $(paxosBC_src) $(paxosQC_src) $(paxosQCB_src) $(pbftBC_src) $(pbftBCO_src)
	@echo "compiling gorums proto files" 

%.pb.go %_gorums.pb.go : %.proto
	-@protoc -I=$(gorums_include):. \
		--go_out=paths=source_relative:. \
		--gorums_out=paths=source_relative:. \
		$<

clean:
	-rm $(paxosBC_src) $(paxosQC_src) $(paxosQCB_src) $(pbftBC_src) $(pbftBCO_src)

wd := $(shell pwd)
csv_path := $(wd)/csv

histogram:
	@cd Master/benchmarks/util/; go run util.go -num=$(NUM) -bench=$(BENCH) -path=$(csv_path)