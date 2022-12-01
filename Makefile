APPNAME=ckp

.DEFAULT_GOAL := help

## test: run tests on cmd and pkg files.
.PHONY: test
test: vet fmt
	go test ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

## build: build application binary.
.PHONY: build
build:
	go build -o $(APPNAME)

.PHONY: check-lint
check-lint:
ifeq (, $(shell which golangci-lint))
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.23.8
endif
ifeq (, $(shell which errcheck))
	go get -u github.com/kisielk/errcheck
endif

## lint: run linters over the entire code base
.PHONY: lint
lint: check-lint
	golangci-lint run ./... --timeout 15m0s
	errcheck -exclude ./.golangci-errcheck-exclude.txt ./...

## install-hooks: install hooks
.PHONY: install-hooks
install-hooks:
	ln -s $(PWD)/githooks/pre-push .git/hooks/pre-push

## mockgen: generate mocks
.PHONY: mockgen
mockgen:
	mockgen -source internal/exec/exec.go -destination mocks/IExec.go -package=mocks
	mockgen -source internal/printers/printers.go -destination mocks/IPrinters.go -package=mocks


## clean: remove releases
.PHONY: clean
clean:
	rm -rf $(APPNAME)

.PHONY: all
all: help

.PHONY: help
help: Makefile
	@echo " Choose a command..."
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
