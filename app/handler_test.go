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

    t.Run("Ping should reply with PONG", func(t *testing.T) {
        buffer := &bytes.Buffer{}
        mockHandler := Handler{writer: bufio.NewWriter(buffer)}
        input := []string{"PING"}
        expected := "+PONG\r\n"

        err := mockHandler.Ping(input)

        AssertNoError(t, err)

        received := buffer.String()

        AssertStringEqual(t, received, expected)
    })

    t.Run("Ping should return an error if argument length is not 1", func(t *testing.T) {
        buffer := &bytes.Buffer{}
        mockHandler := Handler{writer: bufio.NewWriter(buffer)}
        input := []string{"PING", "some", "data"}

        err := mockHandler.Ping(input)

        AssertError(t, err)
    })

    t.Run("Ping should return an error if first argument is not 'ping'", func(t *testing.T) {
        buffer := &bytes.Buffer{}
        mockHandler := Handler{writer: bufio.NewWriter(buffer)}
        input := []string{"SomeCommand"}

        err := mockHandler.Ping(input)

        AssertError(t, err)
    })
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

func TestSetValue(t *testing.T) {

    t.Run("Set should save the provided key value pair in the storage and reply with OK", func(t *testing.T) {
        buffer := &bytes.Buffer{}
        kv := NewInMemoryKV(nil)
        mockHandler := Handler{
            writer: bufio.NewWriter(buffer),
            store:  kv,
        }
        key := "some_key"
        value := "some_value"
        input := []string{"SET", key, value}
        expected := "+OK\r\n"

        mockHandler.Set(input)

        received_tpl := kv.data[key]
        AssertStringEqual(t, received_tpl.value, value)
        AssertZeroTime(t, received_tpl.expiry)

        received := buffer.String()

        AssertStringEqual(t, received, expected)
    })

    t.Run("Set should save the provided key value pair along with exipiration time in the storage and reply with OK", func(t *testing.T) {
        buffer := &bytes.Buffer{}
        tm := time.Now()
        clk := mockClock{tm: tm}
        kv := NewInMemoryKV(&clk)
        mockHandler := Handler{
            writer: bufio.NewWriter(buffer),
            store:  kv,
        }
        key := "key_expiry"
        value := "value expiry"
        input := []string{"SET", key, value, "PX", "200"}
        expected := "+OK\r\n"

        mockHandler.Set(input)

        received_tpl := kv.data[key]
        AssertStringEqual(t, received_tpl.value, value)
        AssertTimeEqual(t, received_tpl.expiry, tm.Add(200*time.Millisecond))

        received := buffer.String()

        AssertStringEqual(t, received, expected)
    })

}

func TestGetValue(t *testing.T) {

    t.Run("Get should return the value for the provided key", func(t *testing.T) {
        kv := NewInMemoryKV(nil)
        key := "test_key"
        value := "test_value"
        kv.Set(key, value)
        buffer := &bytes.Buffer{}
        mockHandler := Handler{
            writer: bufio.NewWriter(buffer),
            store:  kv,
        }
        input := []string{"GET", key}
        expected := "$10\r\ntest_value\r\n"

        mockHandler.Get(input)

        received := buffer.String()

        AssertStringEqual(t, received, expected)
    })

    t.Run("Get should return null bulk string if key doesn't exist", func(t *testing.T) {
        kv := NewInMemoryKV(nil)
        key := "doesnt_exist"
        buffer := &bytes.Buffer{}
        mockHandler := Handler{
            writer: bufio.NewWriter(buffer),
            store:  kv,
        }
        input := []string{"GET", key}
        expected := "$-1\r\n"

        mockHandler.Get(input)

        received := buffer.String()

        AssertStringEqual(t, received, expected)
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

    server, client := net.Pipe()
    kv := NewInMemoryKV(&TimeWrapper{})
    handler := NewHandler(server, kv)
    clientRW := bufio.NewReadWriter(bufio.NewReader(client), bufio.NewWriter(client))

    go handler.HandleConnection()

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
        {
            "SET should reply with OK",
            "*3\r\n$3\r\nSET\r\n$6\r\ntest_k\r\n$6\r\ntest_v\r\n",
            "+OK\r\n",
        },
        {
            "GET should reply with correct value for the key",
            "*2\r\n$3\r\nGET\r\n$6\r\ntest_k\r\n",
            "$6\r\ntest_v\r\n",
        },
        {
            "GET should reply with null bulk string for non-existing key",
            "*2\r\n$3\r\nGET\r\n$8\r\nkey_miss\r\n",
            "$-1\r\n",
        },
        {
            "SET with PX argument should reply with OK",
            "*5\r\n$3\r\nSET\r\n$5\r\nkey_e\r\n$5\r\nval_e\r\n$2\r\nPX\r\n$3\r\n500\r\n",
            "+OK\r\n",
        },
        {
            "GET should reply with correct value before expiry",
            "*2\r\n$3\r\nGET\r\n$5\r\nkey_e\r\n",
            "$5\r\nval_e\r\n",
        },
        {
            "GET should reply with null bulk string for expired key",
            "*2\r\n$3\r\nGET\r\n$5\r\nkey_e\r\n",
            "$-1\r\n",
        },
    }

    RunHandlerTest := func(description, payload, expected string) {
        t.Run(description, func(t *testing.T) {
            _, err := clientRW.WriteString(payload)
            clientRW.Flush()
            AssertNoError(t, err)

            buffer := make([]byte, len(expected))
            _, err = clientRW.Read(buffer)
            AssertNoError(t, err)

            received := string(buffer)
            AssertStringEqual(t, received, expected)
        })
    }

    for _, test := range tests[:7] {
        RunHandlerTest(test.description, test.payload, test.expected)
    }

    time.Sleep(500 * time.Millisecond)

    RunHandlerTest(tests[7].description, tests[7].payload, tests[7].expected)

    t.Run("Handler should close the connection once the client closes it", func(t *testing.T) {
        buffer := make([]byte, 5)

        client.Close()
        // Need to wait so that handler can close the connection
        time.Sleep(50 * time.Millisecond)

        _, err := server.Read(buffer)
        if err != io.ErrClosedPipe {
            t.Errorf("Expected %v but received %v", io.ErrClosedPipe, err)
        }

    })

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
