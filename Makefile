APPNAME=ckp

.PHONY: build help check-lint clean build lint test
.DEFAULT_GOAL := help

## test: run tests on cmd and pkg files.
test: vet fmt
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

## build: build application binary.
build:
	go build -o $(APPNAME)

check-lint:
ifeq (, $(shell which golangci-lint))
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.23.8
endif

lint: check-lint
	golangci-lint run ./... --timeout 15m0s

## clean: remove releases
clean:
	rm -rf $(APPNAME)

all: help
help: Makefile
	@echo " Choose a command..."
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
