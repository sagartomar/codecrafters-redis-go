package main

import (
    "fmt"
    "net"
    "os"
)

func main() {
    SetupServer(6379)
}

func SetupServer(port uint32) {
    kv := NewInMemoryKV(&TimeWrapper{})

    l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
    if err != nil {
        fmt.Println("Failed to bind to port 6379")
        os.Exit(1)
    }
    fmt.Printf("Listening on port: %d", port)

    for {
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting connection: ", err.Error())
            os.Exit(1)
        }

        handler := NewHandler(conn, kv)
        go handler.HandleConnection()
    }
}
