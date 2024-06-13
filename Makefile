.PHONY: test
test: 
	go test ./...

.PHONY: lint
lint: 
	golangci-lint run

.PHONY: all
all: test lint