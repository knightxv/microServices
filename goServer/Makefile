
GOPATH:=$(shell go env GOPATH)

.PHONY: proto test docker


proto:
	protoc --proto_path=${GOPATH}/src:. --micro_out=. --go_out=. proto/goServer.proto

build: proto

	go build -o goServer main.go

test:
	go test -v ./... -cover

docker:
	docker build . -t goServer:latest
