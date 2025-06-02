#!/bin/bash

echo "Building RushKV Distributed Storage System..."

# Generate protobuf code
echo "Generating protobuf code..."
protoc --go_out=. --go-grpc_out=. proto/rushkv.proto

# Download dependencies
echo "Downloading dependencies..."
go mod tidy

# Build server
echo "Building server..."
go build -o rushkv main.go

# Build command line client
echo "Building command line client..."
go build -o rushkv-cli cmd/cli/main.go

# Set executable permissions
chmod +x rushkv
chmod +x rushkv-cli
chmod +x run_cluster.sh

echo "Build completed!"
echo "Executables:"
echo "  - rushkv        (server)"
echo "  - rushkv-cli    (command line client)"
echo "  - run_cluster.sh (cluster startup script)"