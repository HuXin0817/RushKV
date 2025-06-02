package hash

import (
	"crypto/sha1"
	"sort"
	"strconv"
	"sync"
)

type ConsistentHash struct {
	replicas int
	keys     []int
	hashMap  map[int]string
	mutex    sync.RWMutex
}

func NewConsistentHash(replicas int) *ConsistentHash {
	return &ConsistentHash{
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
}

func (ch *ConsistentHash) hash(key string) int {
	h := sha1.New()
	h.Write([]byte(key))
	hashBytes := h.Sum(nil)

	// 取前4个字节转换为int
	hash := int(hashBytes[0])<<24 + int(hashBytes[1])<<16 + int(hashBytes[2])<<8 + int(hashBytes[3])
	if hash < 0 {
		hash = -hash
	}
	return hash
}

func (ch *ConsistentHash) AddNode(node string) {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()

	for i := 0; i < ch.replicas; i++ {
		key := ch.hash(node + strconv.Itoa(i))
		ch.keys = append(ch.keys, key)
		ch.hashMap[key] = node
	}
	sort.Ints(ch.keys)
}

func (ch *ConsistentHash) RemoveNode(node string) {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()

	for i := 0; i < ch.replicas; i++ {
		key := ch.hash(node + strconv.Itoa(i))
		delete(ch.hashMap, key)

		// 从keys中移除
		for j, k := range ch.keys {
			if k == key {
				ch.keys = append(ch.keys[:j], ch.keys[j+1:]...)
				break
			}
		}
	}
}

func (ch *ConsistentHash) GetNode(key string) string {
	ch.mutex.RLock()
	defer ch.mutex.RUnlock()

	if len(ch.keys) == 0 {
		return ""
	}

	hash := ch.hash(key)

	// 二分查找第一个大于等于hash的节点
	idx := sort.Search(len(ch.keys), func(i int) bool {
		return ch.keys[i] >= hash
	})

	// 如果没找到，则使用第一个节点（环形）
	if idx == len(ch.keys) {
		idx = 0
	}

	return ch.hashMap[ch.keys[idx]]
}

func (ch *ConsistentHash) GetNodes() []string {
	ch.mutex.RLock()
	defer ch.mutex.RUnlock()

	nodeSet := make(map[string]bool)
	for _, node := range ch.hashMap {
		nodeSet[node] = true
	}

	nodes := make([]string, 0, len(nodeSet))
	for node := range nodeSet {
		nodes = append(nodes, node)
	}

	return nodes
}
