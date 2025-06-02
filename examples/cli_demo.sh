#!/bin/bash

echo "RushKV Command Line Client Demo"
echo "=============================="

# Check if server is running
echo "Checking server status..."
if ! pgrep -f "rushkv" > /dev/null; then
    echo "Starting RushKV server..."
    ./rushkv -id=demo -addr=localhost -port=8080 -data=./demo_data &
    SERVER_PID=$!
    sleep 3
    echo "Server started (PID: $SERVER_PID)"
else
    echo "Server is already running"
fi

echo ""
echo "Demonstrating basic operations:"
echo "==============================="

# Basic operations demo
echo "1. Storing data"
./rushkv-cli -server=localhost:8080 -batch -commands="put user:1 {\"name\":\"Alice\",\"age\":30}"

echo ""
echo "2. Retrieving data"
./rushkv-cli -server=localhost:8080 -batch -commands="get user:1"

echo ""
echo "3. Checking if key exists"
./rushkv-cli -server=localhost:8080 -batch -commands="exists user:1"

echo ""
echo "4. Viewing cluster information"
./rushkv-cli -server=localhost:8080 -batch -commands="cluster"

echo ""
echo "5. Deleting data"
./rushkv-cli -server=localhost:8080 -batch -commands="delete user:1"

echo ""
echo "6. Confirming deletion"
./rushkv-cli -server=localhost:8080 -batch -commands="exists user:1"

echo ""
echo "Demo completed!"
echo "To enter interactive mode, run: ./rushkv-cli"

# Cleanup
if [ ! -z "$SERVER_PID" ]; then
    echo "Stopping demo server..."
    kill $SERVER_PID
    rm -rf ./demo_data
fi