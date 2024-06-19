# Copyright 2024 Hewlett Packard Enterprise Development LP

.PHONY: test
test: 
	go test ./...

.PHONY: lint
lint: 
	golangci-lint run

.PHONY: all
all: test lint