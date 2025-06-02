package storage

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"
    
    "github.com/boltdb/bolt"
)

type StorageEngine struct {
    db       *bolt.DB
    dataPath string
    mutex    sync.RWMutex
}

func NewStorageEngine(dataPath string) (*StorageEngine, error) {
    if err := os.MkdirAll(dataPath, 0755); err != nil {
        return nil, fmt.Errorf("failed to create data directory: %v", err)
    }
    
    dbPath := filepath.Join(dataPath, "rushkv.db")
    db, err := bolt.Open(dbPath, 0600, &bolt.Options{
        Timeout: 1 * time.Second,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %v", err)
    }
    
    // 创建默认bucket
    err = db.Update(func(tx *bolt.Tx) error {
        _, err := tx.CreateBucketIfNotExists([]byte("kv"))
        return err
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create bucket: %v", err)
    }
    
    return &StorageEngine{
        db:       db,
        dataPath: dataPath,
    }, nil
}

func (se *StorageEngine) Put(key string, value []byte) error {
    se.mutex.Lock()
    defer se.mutex.Unlock()
    
    kvPair := &KVPair{
        Key:       key,
        Value:     value,
        Version:   time.Now().UnixNano(),
        Timestamp: time.Now(),
        Deleted:   false,
    }
    
    data, err := json.Marshal(kvPair)
    if err != nil {
        return fmt.Errorf("failed to marshal data: %v", err)
    }
    
    return se.db.Update(func(tx *bolt.Tx) error {
        bucket := tx.Bucket([]byte("kv"))
        return bucket.Put([]byte(key), data)
    })
}

func (se *StorageEngine) Get(key string) ([]byte, error) {
    se.mutex.RLock()
    defer se.mutex.RUnlock()
    
    var result []byte
    err := se.db.View(func(tx *bolt.Tx) error {
        bucket := tx.Bucket([]byte("kv"))
        data := bucket.Get([]byte(key))
        if data == nil {
            return fmt.Errorf("key not found")
        }
        
        var kvPair KVPair
        if err := json.Unmarshal(data, &kvPair); err != nil {
            return fmt.Errorf("failed to unmarshal data: %v", err)
        }
        
        if kvPair.Deleted {
            return fmt.Errorf("key not found")
        }
        
        result = kvPair.Value
        return nil
    })
    
    return result, err
}

func (se *StorageEngine) Delete(key string) error {
    se.mutex.Lock()
    defer se.mutex.Unlock()
    
    return se.db.Update(func(tx *bolt.Tx) error {
        bucket := tx.Bucket([]byte("kv"))
        data := bucket.Get([]byte(key))
        if data == nil {
            return fmt.Errorf("key not found")
        }
        
        var kvPair KVPair
        if err := json.Unmarshal(data, &kvPair); err != nil {
            return fmt.Errorf("failed to unmarshal data: %v", err)
        }
        
        kvPair.Deleted = true
        kvPair.Timestamp = time.Now()
        
        newData, err := json.Marshal(kvPair)
        if err != nil {
            return fmt.Errorf("failed to marshal data: %v", err)
        }
        
        return bucket.Put([]byte(key), newData)
    })
}

func (se *StorageEngine) Close() error {
    return se.db.Close()
}

type KVPair struct {
    Key       string    `json:"key"`
    Value     []byte    `json:"value"`
    Version   int64     `json:"version"`
    Timestamp time.Time `json:"timestamp"`
    Deleted   bool      `json:"deleted"`
}