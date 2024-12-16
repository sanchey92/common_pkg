LOCAL_BIN:=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	GOBIN=$(LOCAL_BIN) golangci-lint run ./... --config .golangci.config.yaml