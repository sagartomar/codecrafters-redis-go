package main

import (
	"testing"
	"time"
)

type mockClock struct {
	tm time.Time
}

func (c *mockClock) GetTime() time.Time {
	return c.tm
}

func TestSet(t *testing.T) {

	t.Run("Set should store key, value and zero expiry", func(t *testing.T) {
		key := "test_key"
		expected := "test_value"
		kv := NewInMemoryKV(nil)

		kv.Set(key, expected)

		AssertStringEqual(t, kv.data[key].value, expected)

		AssertZeroTime(t, kv.data[key].expiry)
	})

}

func TestSetWithExpiry(t *testing.T) {

	t.Run("SetWithExpiry should store key, value and correct exipiration time", func(t *testing.T) {
		tm := time.Now()
		clk := mockClock{tm}
		key := "test key"
		expectedValue := "test value with expiry"
		expectedTime := tm.Add(100 * time.Millisecond)
		kv := NewInMemoryKV(&clk)

		kv.SetWithExpiry(key, expectedValue, 100*time.Millisecond)

		AssertStringEqual(t, kv.data[key].value, expectedValue)

		AssertTimeEqual(t, kv.data[key].expiry, expectedTime)
	})

}

func TestGet(t *testing.T) {

	t.Run("Get should return the correct value for provided key and no error", func(t *testing.T) {
		key := "test_key"
		expected := "test_value"
		kv := NewInMemoryKV(nil)
		kv.data[key] = tuple{value: expected}

		err, received := kv.Get(key)

		AssertStringEqual(t, received, expected)
		AssertNoError(t, err)
	})

	t.Run("Get should return empty string and error if key doesn't exist", func(t *testing.T) {
		key := "doesn't_exist"
		expected := ""
		kv := NewInMemoryKV(nil)

		err, received := kv.Get(key)

		AssertStringEqual(t, received, expected)
		AssertError(t, err)
	})

	tm := time.Now()
	clk := mockClock{tm: tm}
	kv := NewInMemoryKV(&clk)
	key := "key_expiry"
	expected := "value_expiry"

	kv.data[key] = tuple{
		value:  expected,
		expiry: tm.Add(200 * time.Millisecond),
	}

	t.Run("Get should return the correct value if not yet expired", func(t *testing.T) {
		clk.tm = tm.Add(100 * time.Millisecond)

		err, received := kv.Get(key)

		AssertNoError(t, err)
		AssertStringEqual(t, received, expected)
	})

	t.Run("Get should return empty string and error if key has expired", func(t *testing.T) {
		clk.tm = tm.Add(250 * time.Millisecond)

		err, received := kv.Get(key)

		AssertError(t, err)
		AssertStringEqual(t, received, "")
	})

}

func AssertTimeEqual(t testing.TB, received, expected time.Time) {
	t.Helper()

	if !received.Equal(expected) {
		t.Errorf("Expected time %s but received time %s", expected, received)
	}
}

func AssertZeroTime(t testing.TB, received time.Time) {
	t.Helper()

	if !received.IsZero() {
		t.Errorf("Expected zero time but received %s", received)
	}
}
