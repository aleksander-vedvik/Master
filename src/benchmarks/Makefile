.PHONY: proto
proto: $(proto_go)
	@echo "+ compiling gorums proto files"

%.pb.go %_gorums.pb.go : %.proto
	@protoc -I=$(gorum_include):. \
		--go_out=paths=source_relative:. \
		--gorums_out=paths=source_relative:. \
		$<