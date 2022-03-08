package main

import (
	"bufio"
	"bytes"
	"net"
	"reflect"
	"testing"
)

func TestPing(t *testing.T) {
    buffer := &bytes.Buffer{}
    mockHandler := Handler {writer: bufio.NewWriter(buffer)}
    expected := "+PONG\r\n"

    err := mockHandler.Ping()

    AssertNoError(t, err)

    received := buffer.String()

    AssertStringEqual(t, received, expected)
}

func TestEcho(t *testing.T) {
    buffer := &bytes.Buffer{}
    mockHandler := Handler {writer: bufio.NewWriter(buffer)}
    input := "testing"
    expected := "$7\r\ntesting\r\n"

    err := mockHandler.Echo(input)

    AssertNoError(t, err)

    received := buffer.String()

    AssertStringEqual(t, received, expected)
}

func TestReadRESPBulkString(t *testing.T) {

    t.Run("ReadRESPBulkString should process the RESP bulk string and return the correct string", func(t *testing.T) {
        buffer := bytes.NewBufferString("$7\r\ntesting\r\n")
        mockHandler := Handler {reader: bufio.NewReader(buffer)}
        expected := "testing"

        received, err := mockHandler.ReadRESPBulkString()

        AssertNoError(t, err)
        AssertStringEqual(t, received, expected)
    })

    t.Run("ReadRESPBulkString should return an error when first character is not '$'", func(t *testing.T) {
        buffer := bytes.NewBufferString("7\r\ntesting\r\n")
        mockHandler := Handler {reader: bufio.NewReader(buffer)}
        
        received, err := mockHandler.ReadRESPBulkString()

        AssertError(t, err)
        AssertStringEqual(t, received, "")
    })

    t.Run("ReadRESPBulkString should return an error when size and string length are not equal", func(t *testing.T) {
        buffer := bytes.NewBufferString("$4\r\ntest string\r\n")
        mockHandler := Handler {reader: bufio.NewReader(buffer)}

        received, err := mockHandler.ReadRESPBulkString()
        
        AssertError(t, err)
        AssertStringEqual(t, received, "")
    })

}

func TestReadRESPArray(t *testing.T) {

    t.Run("ReadRESPArray returns an error if first character is not '*'", func(t *testing.T) {
        buffer := bytes.NewBufferString("2\r\n$4\r\ntest\r\n$5\r\narray\r\n")
        mockHandler := Handler {reader: bufio.NewReader(buffer)}

        _, err := mockHandler.ReadRESPArray()

        AssertError(t, err)
    })

    t.Run("ReadRESPArray returns correct array after processing input", func(t *testing.T) {
        buffer := bytes.NewBufferString("*2\r\n$4\r\ntest\r\n$5\r\narray\r\n")
        mockHandler := Handler {reader: bufio.NewReader(buffer)}
        expected := []string {"test", "array"}

        received, err := mockHandler.ReadRESPArray()

        AssertNoError(t, err)

        if !reflect.DeepEqual(received, expected) {
            t.Errorf("Expected %v but received %v", expected, received)
        }
    })

}

func TestHandler(t *testing.T) {

    tests := []struct {
        description string
        payload string
        expected string
    }{
        {
            "PING should reply with PONG",
            "*1\r\n$4\r\nPING\r\n",
            "+PONG\r\n",
        },
    }

    for _, test := range tests {
        t.Run(test.description, func(t *testing.T) {

            server, client := net.Pipe()
            handler := NewHandler(server)
            clientRW := bufio.NewReadWriter(bufio.NewReader(client), bufio.NewWriter(client))
            reader := bufio.NewReader(client)

            go handler.HandleConnection()
            _, err := clientRW.WriteString(test.payload)
            clientRW.Flush()
            AssertNoError(t, err)

            received, err := reader.ReadString('\n')
            AssertNoError(t, err)

            AssertStringEqual(t, received, test.expected)

        })
    }

}

func AssertError(t testing.TB, err error) {
    t.Helper()
    if err == nil {
        t.Error("Expected an error but didn't receive any")
    }
}

func AssertNoError(t testing.TB, err error) {
    t.Helper()
    if err != nil {
        t.Fatalf("Didn't expect an error but received error %v", err)
    }
}
