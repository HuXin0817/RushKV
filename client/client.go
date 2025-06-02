package client

import (
    "context"
    "fmt"
    "time"
    
    "google.golang.org/grpc"
    "rushkv/proto"
)

type RushKVClient struct {
    conn   *grpc.ClientConn
    client proto.RushKVClient
}

func NewRushKVClient(address string) (*RushKVClient, error) {
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        return nil, fmt.Errorf("failed to connect: %v", err)
    }
    
    client := proto.NewRushKVClient(conn)
    
    return &RushKVClient{
        conn:   conn,
        client: client,
    }, nil
}

func (c *RushKVClient) Put(key string, value []byte) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    resp, err := c.client.Put(ctx, &proto.PutRequest{
        Key:   key,
        Value: value,
    })
    if err != nil {
        return fmt.Errorf("put failed: %v", err)
    }
    
    if !resp.Success {
        return fmt.Errorf("put failed: %s", resp.Error)
    }
    
    return nil
}

func (c *RushKVClient) Get(key string) ([]byte, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    resp, err := c.client.Get(ctx, &proto.GetRequest{
        Key: key,
    })
    if err != nil {
        return nil, fmt.Errorf("get failed: %v", err)
    }
    
    if !resp.Success {
        return nil, fmt.Errorf("get failed: %s", resp.Error)
    }
    
    return resp.Value, nil
}

func (c *RushKVClient) Delete(key string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    resp, err := c.client.Delete(ctx, &proto.DeleteRequest{
        Key: key,
    })
    if err != nil {
        return fmt.Errorf("delete failed: %v", err)
    }
    
    if !resp.Success {
        return fmt.Errorf("delete failed: %s", resp.Error)
    }
    
    return nil
}

func (c *RushKVClient) GetClusterInfo() (*proto.ClusterInfoResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    return c.client.GetClusterInfo(ctx, &proto.ClusterInfoRequest{})
}

func (c *RushKVClient) Close() error {
    return c.conn.Close()
}