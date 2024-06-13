.PHONY: proto docker clean eval

network:
	docker network create --gateway 10.0.0.1 --subnet 10.0.0.0/16 MasterLab

docker:
	docker compose up --build

eval:
	docker compose down
	docker compose up

wd := $(shell pwd)
csv_path := $(wd)/csv

histogram:
	@cd benchmarks/util/; go run util.go -num=$(NUM) -bench=$(BENCH) -path=$(csv_path) -t=$(T)
