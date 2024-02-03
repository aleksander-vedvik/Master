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
	./protos/storage.proto
	@echo "Proto files generated successfully"

docker:
	docker compose up

docker-rebuild:
	docker compose up --build

server:
	cd cmd/server; go run main.go

client:
	cd cmd/client; go run main.go

run: server client
	@echo "Done!" 

tidy:
	go mod tidy

dev_path				:= ./protos
storage_proto			:= $(dev_path)/storage.proto
proto_path 				:= $(dev_path):.

.PHONY: dev 

dev: 
	@rm -f $(dev_path)/storage*.pb.go
	@protoc -I=$(proto_path) \
		--go_out=:. \
		--gorums_out=dev=true:. \
		$(storage_proto)

%.pb.go : %.proto
	@protoc -I=$(proto_path) --go_out=paths=source_relative:. $^

%_grpc.pb.go : %.proto
	@protoc -I=$(proto_path) --go-grpc_out=paths=source_relative:. $^
