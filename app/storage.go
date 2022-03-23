package main

import (
	"fmt"
	"sync"
)

type InMemoryKV struct {
	lock *sync.Mutex
	data map[string]string
}

func NewInMemoryKV() *InMemoryKV {
	return &InMemoryKV{
		lock: &sync.Mutex{},
		data: make(map[string]string),
	}
}

func (kv *InMemoryKV) Set(key, value string) {
	kv.lock.Lock()
	kv.data[key] = value
	kv.lock.Unlock()
}

func (kv *InMemoryKV) Get(key string) (error, string) {
	kv.lock.Lock()
	defer kv.lock.Unlock()
    if val, ok := kv.data[key]; ok {
        return nil, val
    }
    err := fmt.Errorf("Key: %s doesn't exist", key)
    return err, ""
}
