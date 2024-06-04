GO_FILES := $(shell find . -type f -name "*.go" ! -name "*_test.go")

all: build

build:
	@echo "Building..."
	@go build -o bin/main main.go

run:
	@echo "Running..."
	@go run $(GO_FILES) $(args)

test:
	@echo "Testing..."
	@go test ./tests -v

clean:
	@echo "Cleaning..."
	@rm -f bin/*

.PHONY: all build run test clean
