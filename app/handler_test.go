package main

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestPing(t *testing.T) {
	buffer := &bytes.Buffer{}
	mockHandler := Handler{writer: bufio.NewWriter(buffer)}
	expected := "+PONG\r\n"

	err := mockHandler.Ping()

	AssertNoError(t, err)

	received := buffer.String()

	AssertStringEqual(t, received, expected)
}

func TestEcho(t *testing.T) {

    t.Run("Echo should reply with the passed argument", func(t *testing.T) {
        buffer := &bytes.Buffer{}
        mockHandler := Handler{writer: bufio.NewWriter(buffer)}
        input := []string{"ECHO", "testing"}
        expected := "$7\r\ntesting\r\n"

        err := mockHandler.Echo(input)

        AssertNoError(t, err)

        received := buffer.String()

        AssertStringEqual(t, received, expected)
    })

    t.Run("Echo should return an error if argument length is not 2", func(t *testing.T) {
        buffer := &bytes.Buffer{}
        mockHandler := Handler{writer: bufio.NewWriter(buffer)}
        input := []string{"ECHO"}

        err := mockHandler.Echo(input)

        AssertError(t, err)
    })

    t.Run("Echo should return an error if first argument is not 'echo'", func(t *testing.T) {
        buffer := &bytes.Buffer{}
        mockHandler := Handler{writer: bufio.NewWriter(buffer)}
        input := []string{"something", "test string"}

        err := mockHandler.Echo(input)

        AssertError(t, err)
    })
}

func TestReadRESPBulkString(t *testing.T) {

	t.Run("ReadRESPBulkString should process the RESP bulk string and return the correct string", func(t *testing.T) {
		buffer := bytes.NewBufferString("$7\r\ntesting\r\n")
		mockHandler := Handler{reader: bufio.NewReader(buffer)}
		expected := "testing"

		received, err := mockHandler.ReadRESPBulkString()

		AssertNoError(t, err)
		AssertStringEqual(t, received, expected)
	})

	t.Run("ReadRESPBulkString should return an error when first character is not '$'", func(t *testing.T) {
		buffer := bytes.NewBufferString("7\r\ntesting\r\n")
		mockHandler := Handler{reader: bufio.NewReader(buffer)}

		received, err := mockHandler.ReadRESPBulkString()

		AssertError(t, err)
		AssertStringEqual(t, received, "")
	})

	t.Run("ReadRESPBulkString should return an error when string length is not a number", func(t *testing.T) {
		buffer := bytes.NewBufferString("$ab\r\ntesting\r\n")
		mockHandler := Handler{reader: bufio.NewReader(buffer)}

		received, err := mockHandler.ReadRESPBulkString()

		AssertError(t, err)
		AssertStringEqual(t, received, "")

	})

	t.Run("ReadRESPBulkString should return an error when size and string length are not equal", func(t *testing.T) {
		buffer := bytes.NewBufferString("$4\r\ntest string\r\n")
		mockHandler := Handler{reader: bufio.NewReader(buffer)}

		received, err := mockHandler.ReadRESPBulkString()

		AssertError(t, err)
		AssertStringEqual(t, received, "")
	})

}

func TestReadRESPArray(t *testing.T) {

	t.Run("ReadRESPArray returns an error if first character is not '*'", func(t *testing.T) {
		buffer := bytes.NewBufferString("2\r\n$4\r\ntest\r\n$5\r\narray\r\n")
		mockHandler := Handler{reader: bufio.NewReader(buffer)}

		_, err := mockHandler.ReadRESPArray()

		AssertError(t, err)
	})

	t.Run("ReadRESPArray returns an error when the array size is not number", func(t *testing.T) {
		buffer := bytes.NewBufferString("*ab\r\n$4\r\ntest\r\n$5\r\narray\r\n")
		mockHandler := Handler{reader: bufio.NewReader(buffer)}

		_, err := mockHandler.ReadRESPArray()

		AssertError(t, err)
	})

	t.Run("ReadRESPArray returns correct array after processing input", func(t *testing.T) {
		buffer := bytes.NewBufferString("*2\r\n$4\r\ntest\r\n$5\r\narray\r\n")
		mockHandler := Handler{reader: bufio.NewReader(buffer)}
		expected := []string{"test", "array"}

		received, err := mockHandler.ReadRESPArray()

		AssertNoError(t, err)

		if !reflect.DeepEqual(received, expected) {
			t.Errorf("Expected %v but received %v", expected, received)
		}
	})

}

func TestHandler(t *testing.T) {

	t.Run("Handler should close the connection once the client closes it", func(t *testing.T) {
		server, client := net.Pipe()
		handler := NewHandler(server)
		buffer := make([]byte, 5)

		go handler.HandleConnection()
		client.Close()
		// Need to wait so that handler can close the connection
		time.Sleep(50 * time.Millisecond)

		_, err := server.Read(buffer)
		if err != io.ErrClosedPipe {
			t.Errorf("Expected %v but received %v", io.ErrClosedPipe, err)
		}

	})

	tests := []struct {
		description string
		payload     string
		expected    string
	}{
		{
			"PING should reply with PONG",
			"*1\r\n$4\r\nPING\r\n",
			"+PONG\r\n",
		},
		{
			"ECHO should reply with the argument back",
			"*2\r\n$4\r\nECHO\r\n$4\r\ntest\r\n",
			"$4\r\ntest\r\n",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {

			server, client := net.Pipe()
			handler := NewHandler(server)
			clientRW := bufio.NewReadWriter(bufio.NewReader(client), bufio.NewWriter(client))

			go handler.HandleConnection()
			_, err := clientRW.WriteString(test.payload)
			clientRW.Flush()
			AssertNoError(t, err)

			buffer := make([]byte, len(test.expected))
			_, err = clientRW.Read(buffer)
			AssertNoError(t, err)

			received := string(buffer)
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
