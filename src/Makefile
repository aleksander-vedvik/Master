.PHONY: proto docker clean eval

gorums_include := $(shell go list -m -f {{.Dir}} github.com/relab/gorums)
paxosBC_src := paxos.bc/proto/multipaxos.pb.go
paxosQC_src := paxosqc/proto/paxosqc.pb.go
paxosQCB_src := paxosqcb/proto/paxosqcb.pb.go
pbftBC_src := pbft.gorums/protos/node.pb.go

proto: clean $(paxosBC_src) $(paxosQC_src) $(paxosQCB_src) $(pbftBC_src) $(pbftBCO_src)
	@echo "compiling gorums proto files" 

%.pb.go %_gorums.pb.go : %.proto
	-@protoc -I=$(gorums_include):. \
		--go_out=paths=source_relative:. \
		--gorums_out=paths=source_relative:. \
		$<

clean:
	-rm $(paxosBC_src) $(paxosQC_src) $(paxosQCB_src) $(pbftBC_src) 

wd := $(shell pwd)
csv_path := $(wd)/csv

binary:
	CGO_ENABLED=0 GOOS=linux go build -o ./app