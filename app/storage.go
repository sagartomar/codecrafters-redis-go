package main

import "sync"

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

func (kv *InMemoryKV) Get(key string) string {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	return kv.data[key]
}
