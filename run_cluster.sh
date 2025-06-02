#!/bin/bash

# 启动3个节点的集群
echo "Starting RushKV cluster..."

# 启动节点1（leader）
./rushkv -id=node1 -addr=localhost -port=8080 -data=./data/node1 &
NODE1_PID=$!

sleep 2

# 启动节点2
./rushkv -id=node2 -addr=localhost -port=8081 -data=./data/node2 &
NODE2_PID=$!

sleep 2

# 启动节点3
./rushkv -id=node3 -addr=localhost -port=8082 -data=./data/node3 &
NODE3_PID=$!

echo "Cluster started with PIDs: $NODE1_PID, $NODE2_PID, $NODE3_PID"
echo "Press Ctrl+C to stop the cluster"

# 等待信号
trap "kill $NODE1_PID $NODE2_PID $NODE3_PID; exit" SIGINT SIGTERM
wait