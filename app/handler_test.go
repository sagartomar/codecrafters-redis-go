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

    AssertNoError(t, err)

    received := buffer.String()

    AssertStringEqual(t, received, expected)
}

func TestEcho(t *testing.T) {
    buffer := &bytes.Buffer{}
    mockHandler := Handler {writer: buffer}
    input := "testing"
    expected := "$7\r\ntesting\r\n"

    err := mockHandler.Echo(input)

    AssertNoError(t, err)

    received := buffer.String()

    AssertStringEqual(t, received, expected)
}

func AssertNoError(t testing.TB, err error) {
    t.Helper()
    if err != nil {
        t.Fatalf("Didn't expect an error but received error %v", err)
    }
}
