.PHONY: depgraph gogen golint protos run-server run-client test

GO := $(shell which go)
GOLINT := $(shell which golint)
PROTOC := $(shell which protoc)
PROTO_SOURCES := $(shell find . -type f -name '*.proto' -not -path "./vendor/*")
PROTO_FILES := $(patsubst %.proto,%.pb.go,$(PROTO_SOURCES))

depgraph:
	@godepgraph -nostdlib -novendor -o github.com/ubiqueworks/go-solid-tutorial github.com/ubiqueworks/go-solid-tutorial/cmd/server | dot -Tpng -o depgraph.png
	@godepgraph -nostdlib -novendor -o github.com/ubiqueworks/go-solid-tutorial github.com/ubiqueworks/go-solid-tutorial/cmd/server | dot -Tsvg -o depgraph.svg

gogen:
	@$(GO) generate $(shell go list ./... | grep -v /vendor/)

golint:
	@$(GOLINT) -set_exit_status $(shell go list ./... | grep -v /vendor/)

$(PROTO_FILES): %.pb.go: %.proto

%.pb.go:
	@echo Compiling $<
	@$(PROTOC) -I=${GOPATH}/src -I. --go_out=plugins=grpc:. $<

protos: $(PROTO_FILES)

run-client:
	@$(GO) run cmd/client/main.go create --name george
	@$(GO) run cmd/client/main.go create --name matteo
	@$(GO) run cmd/client/main.go create --name john

run-server:
	@$(GO) run cmd/server/main.go

test:
	@$(GO) test -v ./...
