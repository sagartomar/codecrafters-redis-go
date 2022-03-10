package main

import "sync"

type KV struct {
	lock *sync.Mutex
	data map[string]string
}

func NewKV() *KV {
	return &KV{
		lock: &sync.Mutex{},
		data: make(map[string]string),
	}
}

func (kv *KV) Set(key, value string) {
	kv.lock.Lock()
	kv.data[key] = value
	kv.lock.Unlock()
}

func (kv *KV) Get(key string) string {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	return kv.data[key]
}
