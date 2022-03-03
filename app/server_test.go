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
