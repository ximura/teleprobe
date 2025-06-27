MAIN_PACKAGE_PATH := ./cmd/sensor
BINARY_NAME := sensor

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...

## grpcui: start grpcui on port 50051
.PHONY: grpcui
grpcui:
	grpcui -plaintext :50051

## gogenerate:      run go codegen
.PHONY: gogenerate
gogenerate:: 
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/*.proto
	go generate ./...

## build-sensor: build the sensor
.PHONY: build-sensor
build-sensor:
	go build -o=./bin/sensor ./cmd/sensor
	chmod +x ./bin/sensor

## build-sink: build the sink
.PHONY: build-sink
build-sink:
	go build -o=./bin/sink ./cmd/sink
	chmod +x ./bin/sink

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=./bin/coverage.out ./...
	go tool cover -html=./bin/coverage.out

docker-build-sensor:
	docker build -f Dockerfile.sensor -t teleprobe/sensor .

docker-build-sink:
	docker build -f Dockerfile.sink -t teleprobe/sink .