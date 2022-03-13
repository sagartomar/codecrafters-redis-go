package main

import "testing"

func TestSet(t *testing.T) {
	key := "test_key"
	expected := "test_value"
	kv := NewInMemoryKV()

	kv.Set(key, expected)

	AssertStringEqual(t, kv.data[key], expected)
}

func TestGet(t *testing.T) {
	key := "test_key"
	expected := "test_value"
	kv := NewInMemoryKV()
	kv.data[key] = expected

	received := kv.Get(key)

	AssertStringEqual(t, received, expected)
}
