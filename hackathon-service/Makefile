GOPATH:=$(shell go env GOPATH)
GOOS:=$(shell go env GOOS)
GOARCH:=$(shell go env GOARCH)

all: build

.PHONY: build
build: deps
build: clean
	CGO_ENABLED=0 go build -o hackathon-service -v

.PHONY: test
test:
	go test -v ./... -cover

deps:
	go get  github.com/golang/glog
	go get  github.com/micro/go-micro
	go get  github.com/ssgo/s
	go get  github.com/micro/go-micro/metadata
	go get  github.com/micro/go-grpc
	go get  github.com/micro/go-plugins/registry/etcdv3
	go get  github.com/prometheus/client_golang/prometheus/promhttp
	go get  github.com/spf13/viper

clean:
	go clean -x github.com/YouDail/golang_micro/hackathon-service

.PHONY: docker
docker: build
	docker build . -t hackathon-service:latest
