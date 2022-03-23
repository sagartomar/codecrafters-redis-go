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

    t.Run("Get should return the correct value for provided key and no error", func(t *testing.T) {
        key := "test_key"
        expected := "test_value"
        kv := NewInMemoryKV()
        kv.data[key] = expected

        err, received := kv.Get(key)

        AssertStringEqual(t, received, expected)
        AssertNoError(t, err)
    })

    t.Run("Get should return empty string and error if key doesn't exist", func(t *testing.T) {
        key := "doesn't_exist"
        expected := ""
        kv := NewInMemoryKV()

        err, received := kv.Get(key)

        AssertStringEqual(t, received, expected)
        AssertError(t, err)
    })

}
