# RushKV

一个基于 Go 语言实现的分布式键值存储系统，使用 gRPC 进行通信，BoltDB 作为存储引擎，并采用一致性哈希算法实现数据分片。

## 特性

- 🚀 **高性能**: 基于 gRPC 的高效通信协议
- 🔄 **分布式**: 支持多节点集群部署
- 📊 **一致性哈希**: 智能数据分片和负载均衡
- 💾 **持久化存储**: 使用 BoltDB 确保数据持久性
- 🛠️ **简单易用**: 提供命令行客户端和编程接口
- 🔧 **可扩展**: 支持动态节点加入和离开

## 架构

RushKV 采用分布式架构，主要组件包括：

- **Server**: 核心服务节点，处理数据存储和集群管理
- **Client**: 客户端库，提供简洁的 API 接口
- **Storage Engine**: 基于 BoltDB 的存储引擎
- **Consistent Hash**: 一致性哈希算法实现数据分片
- **CLI**: 命令行客户端工具

## 快速开始

### 环境要求

- Go 1.24.3+
- Protocol Buffers 编译器

### 安装

1. 克隆项目

```bash
git clone <repository-url>
cd RushKV
```

2. 安装依赖

```bash
go mod download
```

3. 生成 protobuf 代码并构建

```bash
make build
```

### 启动单节点

```bash
./rushkv -id=node1 -addr=localhost -port=8080 -data=./data/node1
```

### 启动集群

使用提供的脚本启动 3 节点集群：

```bash
./run_cluster.sh
```

这将启动三个节点：

- node1: localhost:8080
- node2: localhost:8081
- node3: localhost:8082

## 使用方法

### 命令行客户端

```bash
# 存储数据
./rushkv-cli -server=localhost:8080 -batch -commands="put user:1 {\"name\":\"Alice\",\"age\":30}"

# 获取数据
./rushkv-cli -server=localhost:8080 -batch -commands="get user:1"

# 删除数据
./rushkv-cli -server=localhost:8080 -batch -commands="delete user:1"
```

### 编程接口

```go
package main

import (
    "log"
    "rushkv/client"
)

func main() {
    // 创建客户端
    client, err := client.NewRushKVClient("localhost:8080")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // 存储数据
    err = client.Put("key1", []byte("value1"))
    if err != nil {
        log.Fatal(err)
    }

    // 获取数据
    value, err := client.Get("key1")
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Value: %s", value)
}
```

## API 接口

RushKV 提供以下 gRPC 接口：

- `Put(key, value)` - 存储键值对
- `Get(key)` - 获取指定键的值
- `Delete(key)` - 删除指定键
- `Join(nodeInfo)` - 节点加入集群
- `Leave(nodeId)` - 节点离开集群
- `GetClusterInfo()` - 获取集群信息

## 配置选项

| 参数    | 描述       | 默认值    |
| ------- | ---------- | --------- |
| `-id`   | 节点 ID    | node1     |
| `-addr` | 服务器地址 | localhost |
| `-port` | 服务器端口 | 8080      |
| `-data` | 数据目录   | ./data    |

## 开发

### 构建命令

```bash
# 构建所有组件
make build

# 只生成protobuf代码
make proto

# 只构建服务器
make server

# 只构建CLI客户端
make cli

# 运行测试
make test

# 清理构建文件
make clean
```

### 项目结构

```
RushKV/
├── client/          # 客户端库
├── cmd/cli/         # 命令行客户端
├── data/            # 数据目录
├── examples/        # 示例脚本
├── hash/            # 一致性哈希实现
├── proto/           # Protocol Buffers定义
├── server/          # 服务器实现
├── storage/         # 存储引擎
├── main.go          # 服务器入口
├── Makefile         # 构建脚本
└── run_cluster.sh   # 集群启动脚本
```

## 示例

查看 `examples/` 目录获取更多使用示例：

```bash
# 运行CLI演示
./examples/cli_demo.sh
```

## 许可证

本项目采用 [MIT 许可证](LICENSE)。

## 贡献

欢迎提交 Issue 和 Pull Request 来改进项目！

## 技术栈

- **语言**: Go 1.24.3
- **通信**: gRPC + Protocol Buffers
- **存储**: BoltDB
- **算法**: 一致性哈希
- **构建**: Make
