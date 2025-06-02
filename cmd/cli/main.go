package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "os"
    "strconv"
    "strings"
    "time"
    
    "rushkv/client"
)

// CLI represents the command line interface client
type CLI struct {
    client *client.RushKVClient
    reader *bufio.Reader
}

// NewCLI creates a new CLI instance
func NewCLI(serverAddr string) (*CLI, error) {
    client, err := client.NewRushKVClient(serverAddr)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to server: %v", err)
    }
    
    return &CLI{
        client: client,
        reader: bufio.NewReader(os.Stdin),
    }, nil
}

// printHelp displays available commands
func (cli *CLI) printHelp() {
    fmt.Println("\nRushKV Command Line Client")
    fmt.Println("Available commands:")
    fmt.Println("  put <key> <value>     - Store a key-value pair")
    fmt.Println("  get <key>             - Retrieve value for a key")
    fmt.Println("  delete <key>          - Delete a key-value pair")
    fmt.Println("  exists <key>          - Check if a key exists")
    fmt.Println("  cluster               - Show cluster information")
    fmt.Println("  stats                 - Show client statistics")
    fmt.Println("  benchmark <n>         - Run performance test (n operations)")
    fmt.Println("  help                  - Show this help message")
    fmt.Println("  exit                  - Exit the client")
    fmt.Println()
}

// handlePut processes the put command
func (cli *CLI) handlePut(args []string) {
    if len(args) < 2 {
        fmt.Println("Error: put command requires key and value arguments")
        fmt.Println("Usage: put <key> <value>")
        return
    }
    
    key := args[0]
    value := strings.Join(args[1:], " ")
    
    start := time.Now()
    err := cli.client.Put(key, []byte(value))
    duration := time.Since(start)
    
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Successfully stored '%s' (took: %v)\n", key, duration)
    }
}

// handleGet processes the get command
func (cli *CLI) handleGet(args []string) {
    if len(args) < 1 {
        fmt.Println("Error: get command requires key argument")
        fmt.Println("Usage: get <key>")
        return
    }
    
    key := args[0]
    
    start := time.Now()
    value, err := cli.client.Get(key)
    duration := time.Since(start)
    
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Value for key '%s': %s (took: %v)\n", key, string(value), duration)
    }
}

// handleDelete processes the delete command
func (cli *CLI) handleDelete(args []string) {
    if len(args) < 1 {
        fmt.Println("Error: delete command requires key argument")
        fmt.Println("Usage: delete <key>")
        return
    }
    
    key := args[0]
    
    start := time.Now()
    err := cli.client.Delete(key)
    duration := time.Since(start)
    
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Successfully deleted '%s' (took: %v)\n", key, duration)
    }
}

// handleExists checks if a key exists
func (cli *CLI) handleExists(args []string) {
    if len(args) < 1 {
        fmt.Println("Error: exists command requires key argument")
        fmt.Println("Usage: exists <key>")
        return
    }
    
    key := args[0]
    
    start := time.Now()
    _, err := cli.client.Get(key)
    duration := time.Since(start)
    
    if err != nil {
        fmt.Printf("Key '%s' does not exist (took: %v)\n", key, duration)
    } else {
        fmt.Printf("Key '%s' exists (took: %v)\n", key, duration)
    }
}

// handleCluster displays cluster information
func (cli *CLI) handleCluster() {
    start := time.Now()
    clusterInfo, err := cli.client.GetClusterInfo()
    duration := time.Since(start)
    
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("\nCluster Information (took: %v):\n", duration)
    fmt.Printf("Leader: %s\n", clusterInfo.Leader)
    fmt.Printf("Node Count: %d\n", len(clusterInfo.Nodes))
    fmt.Println("Node List:")
    
    for _, node := range clusterInfo.Nodes {
        status := "Follower"
        if node.IsLeader {
            status = "Leader"
        }
        fmt.Printf("  - ID: %s, Address: %s:%d, Status: %s\n", 
            node.Id, node.Address, node.Port, status)
    }
    fmt.Println()
}

// handleStats displays client statistics
func (cli *CLI) handleStats() {
    fmt.Println("\nClient Statistics:")
    fmt.Println("  Connection Status: Connected")
    fmt.Printf("  Client Start Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
    fmt.Println()
}

// handleBenchmark runs performance tests
func (cli *CLI) handleBenchmark(args []string) {
    n := 1000 // Default 1000 operations
    if len(args) > 0 {
        if num, err := strconv.Atoi(args[0]); err == nil && num > 0 {
            n = num
        }
    }
    
    fmt.Printf("\nStarting benchmark test (%d operations)...\n", n)
    
    // PUT performance test
    fmt.Println("Testing PUT operations...")
    start := time.Now()
    putErrors := 0
    
    for i := 0; i < n; i++ {
        key := fmt.Sprintf("bench_key_%d", i)
        value := fmt.Sprintf("bench_value_%d_%d", i, time.Now().UnixNano())
        
        if err := cli.client.Put(key, []byte(value)); err != nil {
            putErrors++
        }
        
        if (i+1)%100 == 0 {
            fmt.Printf("  Completed %d/%d PUT operations\n", i+1, n)
        }
    }
    
    putDuration := time.Since(start)
    putThroughput := float64(n-putErrors) / putDuration.Seconds()
    
    // GET performance test
    fmt.Println("Testing GET operations...")
    start = time.Now()
    getErrors := 0
    
    for i := 0; i < n; i++ {
        key := fmt.Sprintf("bench_key_%d", i)
        
        if _, err := cli.client.Get(key); err != nil {
            getErrors++
        }
        
        if (i+1)%100 == 0 {
            fmt.Printf("  Completed %d/%d GET operations\n", i+1, n)
        }
    }
    
    getDuration := time.Since(start)
    getThroughput := float64(n-getErrors) / getDuration.Seconds()
    
    // Cleanup test data
    fmt.Println("Cleaning up test data...")
    for i := 0; i < n; i++ {
        key := fmt.Sprintf("bench_key_%d", i)
        cli.client.Delete(key)
    }
    
    // Display results
    fmt.Printf("\nBenchmark Results:\n")
    fmt.Printf("PUT Operations:\n")
    fmt.Printf("  Total Time: %v\n", putDuration)
    fmt.Printf("  Success Count: %d/%d\n", n-putErrors, n)
    fmt.Printf("  Throughput: %.2f ops/sec\n", putThroughput)
    fmt.Printf("  Average Latency: %.2f ms\n", putDuration.Seconds()*1000/float64(n))
    
    fmt.Printf("GET Operations:\n")
    fmt.Printf("  Total Time: %v\n", getDuration)
    fmt.Printf("  Success Count: %d/%d\n", n-getErrors, n)
    fmt.Printf("  Throughput: %.2f ops/sec\n", getThroughput)
    fmt.Printf("  Average Latency: %.2f ms\n", getDuration.Seconds()*1000/float64(n))
    fmt.Println()
}

// processCommand processes a single command
func (cli *CLI) processCommand(input string) bool {
    input = strings.TrimSpace(input)
    if input == "" {
        return true
    }
    
    parts := strings.Fields(input)
    command := strings.ToLower(parts[0])
    args := parts[1:]
    
    switch command {
    case "put":
        cli.handlePut(args)
    case "get":
        cli.handleGet(args)
    case "delete", "del":
        cli.handleDelete(args)
    case "exists":
        cli.handleExists(args)
    case "cluster":
        cli.handleCluster()
    case "stats":
        cli.handleStats()
    case "benchmark", "bench":
        cli.handleBenchmark(args)
    case "help", "h":
        cli.printHelp()
    case "exit", "quit", "q":
        fmt.Println("Goodbye!")
        return false
    default:
        fmt.Printf("Unknown command: %s\n", command)
        fmt.Println("Type 'help' to see available commands")
    }
    
    return true
}

// runInteractive starts the interactive mode
func (cli *CLI) runInteractive() {
    fmt.Println("Welcome to RushKV Command Line Client!")
    fmt.Println("Type 'help' to see available commands")
    
    for {
        fmt.Print("rushkv> ")
        input, err := cli.reader.ReadString('\n')
        if err != nil {
            fmt.Printf("Error reading input: %v\n", err)
            break
        }
        
        if !cli.processCommand(input) {
            break
        }
    }
}

// runBatch executes commands in batch mode
func (cli *CLI) runBatch(commands []string) {
    for _, cmd := range commands {
        fmt.Printf("Executing: %s\n", cmd)
        if !cli.processCommand(cmd) {
            break
        }
        fmt.Println()
    }
}

// Close closes the CLI client
func (cli *CLI) Close() {
    if cli.client != nil {
        cli.client.Close()
    }
}

func main() {
    var (
        serverAddr = flag.String("server", "localhost:8080", "RushKV server address")
        batchMode  = flag.Bool("batch", false, "Batch mode")
        commands   = flag.String("commands", "", "Batch commands separated by semicolon")
    )
    flag.Parse()
    
    cli, err := NewCLI(*serverAddr)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    defer cli.Close()
    
    if *batchMode && *commands != "" {
        cmdList := strings.Split(*commands, ";")
        cli.runBatch(cmdList)
    } else {
        cli.runInteractive()
    }
}