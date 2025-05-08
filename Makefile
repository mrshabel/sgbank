.PHONY: build run start test lint format clean

build:
	@echo "Compiling source code"
	go build -o bin/sgbank ./cmd/api

run:
	@echo "Running server..."
	go run cmd/api/main.go

start: build
	@echo "Starting server"
	./bin/sgbank

test:
	@echo "Running all tests..."
	go test -v -race -cover ./...

lint:
	@echo "Linting source code..."
	go vet ./...

format:
	@echo "Formatting source code..."
	go fmt ./...

clean:
	@echo "Cleaning previous builds..."
	@rm -rf ./bin


# usage commands
.PHONY: help

help:
	@echo "usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@echo ""
	@echo "Build:		Build source code"
	@echo "Run:		Run the development server"
	@echo "Start:		Start the build"
	@echo "Test:		Run all tests"
	@echo "Lint:		Lint the source code"
	@echo "Format:		Format the source code"
	@echo "Clean:		Clean previous builds"