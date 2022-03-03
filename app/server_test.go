package main

import "testing"

func TestRESPSimpleString(t *testing.T) {
    input := "OK"
    expected := "+OK\r\n"

    output := ConvertToRESPSimpleString(input)

    if output != expected {
        t.Errorf("Expected %s but received %s", expected, output)
    }
}

func TestRESPBulkString(t *testing.T) {
    input := "test"
    expected := "$4\r\ntest\r\n"

    output := ConvertToRESPBulkString(input)

    if output != expected {
        t.Errorf("Expected %s but received %s", expected, output)
    }
}
