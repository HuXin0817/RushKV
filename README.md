# RushKV

A distributed key-value storage system implemented in Go, using gRPC for communication, BoltDB as the storage engine, and consistent hashing algorithm for data sharding.

## Features

- 🚀 **High Performance**: Efficient communication protocol based on gRPC
- 🔄 **Distributed**: Support for multi-node cluster deployment
- 📊 **Consistent Hashing**: Intelligent data sharding and load balancing
- 💾 **Persistent Storage**: Data persistence guaranteed by BoltDB
- 🛠️ **Easy to Use**: Command-line client and programming interface provided
- 🔧 **Scalable**: Support for dynamic node joining and leaving

## Architecture

RushKV adopts a distributed architecture with the following main components:

- **Server**: Core service node that handles data storage and cluster management
- **Client**: Client library providing clean API interface
- **Storage Engine**: Storage engine based on BoltDB
- **Consistent Hash**: Consistent hashing algorithm for data sharding
- **CLI**: Command-line client tool

## Quick Start

### Requirements

- Go 1.24.3+
- Protocol Buffers compiler

### Installation

1. Clone the project

```bash
git clone https://github.com/HuXin0817/RushKV
cd RushKV
```

2. Install dependencies

```bash
go mod download
```

3. Generate protobuf code and build

```bash
make build
```

### Start Single Node

```bash
./rushkv -id=node1 -addr=localhost -port=8080 -data=./data/node1
```

### Start Cluster

Use the provided script to start a 3-node cluster:

```bash
./run_cluster.sh
```

This will start three nodes:

- node1: localhost:8080
- node2: localhost:8081
- node3: localhost:8082

## Usage

### Command Line Client

```bash
# Store data
./rushkv-cli -server=localhost:8080 -batch -commands=\"put user:1 {\\\"name\\\":\\\"Alice\\\",\\\"age\\\":30}\"

# Get data
./rushkv-cli -server=localhost:8080 -batch -commands=\"get user:1\"

# Delete data
./rushkv-cli -server=localhost:8080 -batch -commands=\"delete user:1\"
```

### Programming Interface

```go
package main

import (
    \"log\"
    \"rushkv/client\"
)

func main() {
    // Create client
    client, err := client.NewRushKVClient(\"localhost:8080\")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Store data
    err = client.Put(\"key1\", []byte(\"value1\"))
    if err != nil {
        log.Fatal(err)
    }

    // Get data
    value, err := client.Get(\"key1\")
    if err != nil {
        log.Fatal(err)
    }

    log.Printf(\"Value: %s\", value)
}
```

## API Reference

RushKV provides the following gRPC interfaces:

- `Put(key, value)` - Store key-value pair
- `Get(key)` - Get value for specified key
- `Delete(key)` - Delete specified key
- `Join(nodeInfo)` - Node joins cluster
- `Leave(nodeId)` - Node leaves cluster
- `GetClusterInfo()` - Get cluster information

## Configuration Options

| Parameter | Description    | Default   |
| --------- | -------------- | --------- |
| `-id`     | Node ID        | node1     |
| `-addr`   | Server address | localhost |
| `-port`   | Server port    | 8080      |
| `-data`   | Data directory | ./data    |

## Development

### Build Commands

```bash
# Build all components
make build

# Generate protobuf code only
make proto

# Build server only
make server

# Build CLI client only
make cli

# Run tests
make test

# Clean build files
make clean
```

### Project Structure

```
RushKV/
├── client/          # Client library
├── cmd/cli/         # Command-line client
├── data/            # Data directory
├── examples/        # Example scripts
├── hash/            # Consistent hashing implementation
├── proto/           # Protocol Buffers definitions
├── server/          # Server implementation
├── storage/         # Storage engine
├── main.go          # Server entry point
├── Makefile         # Build script
└── run_cluster.sh   # Cluster startup script
```

## Examples

Check the `examples/` directory for more usage examples:

```bash
# Run CLI demo
./examples/cli_demo.sh
```

## License

This project is licensed under the [MIT License](LICENSE).

## Contributing

Welcome to submit Issues and Pull Requests to improve the project!

## Tech Stack

- **Language**: Go 1.24.3
- **Communication**: gRPC + Protocol Buffers
- **Storage**: BoltDB
- **Algorithm**: Consistent Hashing
- **Build**: Make
