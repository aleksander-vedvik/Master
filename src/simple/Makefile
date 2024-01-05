proto:
	@protoc \
	--go_out=protos \
	--gorums_out=protos \
	--proto_path=protos \
	storage.proto
	@echo "Proto files generated successfully"

protos:
	@protoc -I=$(go list -m -f {{.Dir}} github.com/relab/gorums):. \
	--proto_path=paths=source_relative/protos/ \
	--go_out=paths=source_relative:. \
	--gorums_out=paths=source_relative:. \
	protos/storage.proto
	@echo "Proto files generated successfully"