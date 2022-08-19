package main

import (
	"fmt"
	"sync"
	"time"
)

type clock interface {
	GetTime() time.Time
}

type TimeWrapper struct {
}

func (tm *TimeWrapper) GetTime() time.Time {
	return time.Now()
}

type tuple struct {
	value  string
	expiry time.Time
}

type InMemoryKV struct {
	lock *sync.RWMutex
	data map[string]tuple
	clk  clock
}

func NewInMemoryKV(clk clock) *InMemoryKV {
	return &InMemoryKV{
		lock: &sync.RWMutex{},
		data: make(map[string]tuple),
		clk:  clk,
	}
}

func (kv *InMemoryKV) Set(key, value string) {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	kv.data[key] = tuple{value: value}
}

func (kv *InMemoryKV) SetWithExpiry(key, value string, dur time.Duration) {
	tpl := tuple{
		value:  value,
		expiry: kv.clk.GetTime().Add(dur),
	}
	kv.lock.Lock()
	defer kv.lock.Unlock()
	kv.data[key] = tpl
}

func (kv *InMemoryKV) Get(key string) (error, string) {
	kv.lock.RLock()
	isLocked := true
	if val, ok := kv.data[key]; ok {
		if val.expiry.IsZero() || val.expiry.After(kv.clk.GetTime()) {
			kv.lock.RUnlock()
			return nil, val.value
		}
		kv.lock.RUnlock()
		kv.lock.Lock()
		delete(kv.data, key)
		kv.lock.Unlock()
		isLocked = false
	}
	if isLocked {
		kv.lock.RUnlock()
	}
	err := fmt.Errorf("Key: %s doesn't exist", key)
	return err, ""
}
