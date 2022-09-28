package main

import (
    "net"
    "testing"
    "time"
)

func TestSetupServerShouldHandleMultipleConnections(t *testing.T) {
    go SetupServer(6379)

    // Need to wait for the server to start
    time.Sleep(50 * time.Millisecond)

    // Main goroutine doesn't wait for other goroutines to finish so
    // need to make assesrtions in the main goroutine by sending the
    // errors via the channel
    errs := make(chan error)

    for i := 0; i < 2; i++ {
        go func() {
            conn, err := net.Dial("tcp", "0.0.0.0:6379")
            errs <- err

            err = conn.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
            errs <- err

            _, err = conn.Write([]byte("*1\r\n$4\r\nPING\r\n"))
            errs <- err

            read := make([]byte, 7)
            _, err = conn.Read(read)
            errs <- err
        }()
    }

    for i := 0; i < 8; i++ {
        err := <-errs
        AssertNoError(t, err)
    }
}
