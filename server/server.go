package server

import (
    "context"
    "fmt"
    "log"
    "net"
    "sync"
    
    "google.golang.org/grpc"
    "rushkv/hash"
    "rushkv/proto"
    "rushkv/storage"
)

type RushKVServer struct {
    proto.UnimplementedRushKVServer
    
    nodeID      string
    address     string
    port        int
    storage     *storage.StorageEngine
    hash        *hash.ConsistentHash
    nodes       map[string]*proto.NodeInfo
    isLeader    bool
    mutex       sync.RWMutex
    grpcServer  *grpc.Server
}

func NewRushKVServer(nodeID, address string, port int, dataPath string) (*RushKVServer, error) {
    storageEngine, err := storage.NewStorageEngine(dataPath)
    if err != nil {
        return nil, fmt.Errorf("failed to create storage engine: %v", err)
    }
    
    return &RushKVServer{
        nodeID:   nodeID,
        address:  address,
        port:     port,
        storage:  storageEngine,
        hash:     hash.NewConsistentHash(3), // 每个节点3个虚拟节点
        nodes:    make(map[string]*proto.NodeInfo),
        isLeader: true, // 简化实现，第一个节点为leader
    }, nil
}

func (s *RushKVServer) Put(ctx context.Context, req *proto.PutRequest) (*proto.PutResponse, error) {
    // 检查key应该存储在哪个节点
    targetNode := s.hash.GetNode(req.Key)
    if targetNode != s.nodeID {
        // 转发到正确的节点（简化实现，这里直接返回错误）
        return &proto.PutResponse{
            Success: false,
            Error:   fmt.Sprintf("key should be stored on node %s", targetNode),
        }, nil
    }
    
    err := s.storage.Put(req.Key, req.Value)
    if err != nil {
        return &proto.PutResponse{
            Success: false,
            Error:   err.Error(),
        }, nil
    }
    
    return &proto.PutResponse{
        Success: true,
    }, nil
}

func (s *RushKVServer) Get(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
    targetNode := s.hash.GetNode(req.Key)
    if targetNode != s.nodeID {
        return &proto.GetResponse{
            Success: false,
            Error:   fmt.Sprintf("key should be retrieved from node %s", targetNode),
        }, nil
    }
    
    value, err := s.storage.Get(req.Key)
    if err != nil {
        return &proto.GetResponse{
            Success: false,
            Error:   err.Error(),
        }, nil
    }
    
    return &proto.GetResponse{
        Success: true,
        Value:   value,
    }, nil
}

func (s *RushKVServer) Delete(ctx context.Context, req *proto.DeleteRequest) (*proto.DeleteResponse, error) {
    targetNode := s.hash.GetNode(req.Key)
    if targetNode != s.nodeID {
        return &proto.DeleteResponse{
            Success: false,
            Error:   fmt.Sprintf("key should be deleted from node %s", targetNode),
        }, nil
    }
    
    err := s.storage.Delete(req.Key)
    if err != nil {
        return &proto.DeleteResponse{
            Success: false,
            Error:   err.Error(),
        }, nil
    }
    
    return &proto.DeleteResponse{
        Success: true,
    }, nil
}

func (s *RushKVServer) Join(ctx context.Context, req *proto.JoinRequest) (*proto.JoinResponse, error) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    nodeInfo := &proto.NodeInfo{
        Id:       req.NodeId,
        Address:  req.Address,
        Port:     req.Port,
        IsLeader: false,
    }
    
    s.nodes[req.NodeId] = nodeInfo
    s.hash.AddNode(req.NodeId)
    
    log.Printf("Node %s joined the cluster", req.NodeId)
    
    return &proto.JoinResponse{
        Success: true,
    }, nil
}

func (s *RushKVServer) Leave(ctx context.Context, req *proto.LeaveRequest) (*proto.LeaveResponse, error) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    delete(s.nodes, req.NodeId)
    s.hash.RemoveNode(req.NodeId)
    
    log.Printf("Node %s left the cluster", req.NodeId)
    
    return &proto.LeaveResponse{
        Success: true,
    }, nil
}

func (s *RushKVServer) GetClusterInfo(ctx context.Context, req *proto.ClusterInfoRequest) (*proto.ClusterInfoResponse, error) {
    s.mutex.RLock()
    defer s.mutex.RUnlock()
    
    nodes := make([]*proto.NodeInfo, 0, len(s.nodes))
    var leader string
    
    for _, node := range s.nodes {
        nodes = append(nodes, node)
        if node.IsLeader {
            leader = node.Id
        }
    }
    
    return &proto.ClusterInfoResponse{
        Nodes:  nodes,
        Leader: leader,
    }, nil
}

func (s *RushKVServer) Start() error {
    lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.address, s.port))
    if err != nil {
        return fmt.Errorf("failed to listen: %v", err)
    }
    
    s.grpcServer = grpc.NewServer()
    proto.RegisterRushKVServer(s.grpcServer, s)
    
    // 将自己添加到集群
    s.hash.AddNode(s.nodeID)
    s.nodes[s.nodeID] = &proto.NodeInfo{
        Id:       s.nodeID,
        Address:  s.address,
        Port:     int32(s.port),
        IsLeader: s.isLeader,
    }
    
    log.Printf("RushKV server %s starting on %s:%d", s.nodeID, s.address, s.port)
    return s.grpcServer.Serve(lis)
}

func (s *RushKVServer) Stop() {
    if s.grpcServer != nil {
        s.grpcServer.GracefulStop()
    }
    if s.storage != nil {
        s.storage.Close()
    }
}