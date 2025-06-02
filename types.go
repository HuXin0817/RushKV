package main

import (
    "time"
)

// KVPair 键值对结构
type KVPair struct {
    Key       string    `json:"key"`
    Value     []byte    `json:"value"`
    Version   int64     `json:"version"`
    Timestamp time.Time `json:"timestamp"`
    Deleted   bool      `json:"deleted"`
}

// Node 节点信息
type Node struct {
    ID       string `json:"id"`
    Address  string `json:"address"`
    Port     int    `json:"port"`
    IsLeader bool   `json:"is_leader"`
}

// Cluster 集群信息
type Cluster struct {
    Nodes   map[string]*Node `json:"nodes"`
    Leader  string           `json:"leader"`
    Version int64            `json:"version"`
}

// Request 请求结构
type Request struct {
    Type      string            `json:"type"`
    Key       string            `json:"key"`
    Value     []byte            `json:"value"`
    Metadata  map[string]string `json:"metadata"`
    Timestamp time.Time         `json:"timestamp"`
}

// Response 响应结构
type Response struct {
    Success   bool              `json:"success"`
    Data      []byte            `json:"data"`
    Error     string            `json:"error"`
    Metadata  map[string]string `json:"metadata"`
    Timestamp time.Time         `json:"timestamp"`
}