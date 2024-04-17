proto:
	protoc -I=$(go list -m -f {{.Dir}} github.com/relab/gorums):. \
	--go_out=paths=source_relative:. \
	--gorums_out=paths=source_relative:. \
	*.proto

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