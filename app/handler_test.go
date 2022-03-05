package main

import (
	"bytes"
	"testing"
)

func TestPing(t *testing.T) {
    buffer := &bytes.Buffer{}
    mockHandler := Handler {writer: buffer}
    expected := "+PONG\r\n"

    err := mockHandler.Ping()

    if err != nil {
        t.Fatalf("Didn't expect error but received error %v", err)
    }

    received := buffer.String()

    AssertStringEqual(t, received, expected)
}
