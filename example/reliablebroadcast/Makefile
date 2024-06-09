.PHONY: proto 

gorums_include := $(shell go list -m -f {{.Dir}} github.com/relab/gorums)
proto_src := reliablebroadcast/proto/rb.pb.go

proto: clean $(proto_src) 
	@echo "compiling gorums proto files" 

%.pb.go %_gorums.pb.go : %.proto
	-@protoc -I=$(gorums_include):. \
		--go_out=paths=source_relative:. \
		--gorums_out=paths=source_relative:. \
		$<

clean:
	-rm $(proto_src) 
