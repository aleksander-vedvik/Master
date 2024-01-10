proto:
	protoc -I=$(go list -m -f {{.Dir}} github.com/relab/gorums):. \
	--go_out=paths=source_relative:. \
	--gorums_out=paths=source_relative:. \
	storage.proto

protos:
	@protoc -I=$(go list -m -f {{.Dir}} github.com/relab/gorums):. \
	--proto_path=paths=source_relative/protos/ \
	--go_out=paths=source_relative:. \
	--gorums_out=paths=source_relative:. \
	protos/storage.proto
	@echo "Proto files generated successfully"

docker:
	docker compose up

docker-rebuild:
	docker compose up --build