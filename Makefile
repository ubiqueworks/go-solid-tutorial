.PHONY: build depgraph gogen golint package protos test update update-repos vendor

GO := $(shell which go)
GOLINT := $(shell which golint)
PROTOC := $(shell which protoc)
PROTO_SOURCES := $(shell find . -type f -name '*.proto' -not -path "./vendor/*")
PROTO_FILES := $(patsubst %.proto,%.pb.go,$(PROTO_SOURCES))

gogen:
	@$(GO) generate $(shell go list ./... | grep -v /vendor/)

golint:
	@$(GOLINT) -set_exit_status $(shell go list ./... | grep -v /vendor/)

govet:
	@$(GO) vet $(shell go list ./... | grep -v /vendor/)

depgraph:
	@godepgraph -nostdlib -novendor -o github.com/ubiqueworks/go-solid-tutorial github.com/ubiqueworks/go-solid-tutorial/cmd/server | dot -Tpng -o depgraph.png
	@godepgraph -nostdlib -novendor -o github.com/ubiqueworks/go-solid-tutorial github.com/ubiqueworks/go-solid-tutorial/cmd/server | dot -Tsvg -o depgraph.svg

$(PROTO_FILES): %.pb.go: %.proto

%.pb.go:
	@echo Compiling $<
	@$(PROTOC) -I=${GOPATH}/src -I. --go_out=plugins=grpc:. $<

protos: $(PROTO_FILES)

run:
	@$(GO) run cmd/main.go

test:
	@$(GO) test -v ./...

vendor:
	@echo "Vendoring dependencies..."
	@rm -Rf ./vendor
	@$(GO) mod vendor

#package: protos vendor
#	@echo "Building container image... [version=${IMAGE_VERSION}, build=${IMAGE_BUILD}]"
#	@$(DOCKER) pull gcr.io/distroless/base
#	@$(DOCKER) build --rm \
#		--build-arg SERVICE_BUILD=${IMAGE_BUILD} \
#		--build-arg SERVICE_VERSION=${IMAGE_VERSION} \
#		--tag ${IMAGE_NAME}:${IMAGE_VERSION} .
