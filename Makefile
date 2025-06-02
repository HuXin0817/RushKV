.PHONY: all build clean test proto server cli demo

# Default target
all: build

# Build all components
build: proto server cli

# Generate protobuf code
proto:
	@echo "Generating protobuf code..."
	protoc --go_out=. --go-grpc_out=. proto/rushkv.proto

# Build server
server:
	@echo "Building server..."
	go build -o rushkv main.go

# Build command line client
cli:
	@echo "Building command line client..."
	go build -o rushkv-cli cmd/cli/main.go

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Clean build files
clean:
	@echo "Cleaning build files..."
	rm -f rushkv rushkv-cli
	rm -rf data/
	rm -rf demo_data/

# Run demo
demo: build
	@echo "Running demo..."
	chmod +x examples/cli_demo.sh
	./examples/cli_demo.sh

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run single node server
run-server:
	@echo "Starting single node server..."
	./rushkv -id=node1 -addr=localhost -port=8080 -data=./data/node1

# Run cluster
run-cluster:
	@echo "Starting cluster..."
	chmod +x run_cluster.sh
	./run_cluster.sh

# Run interactive client
run-cli:
	@echo "Starting command line client..."
	./rushkv-cli