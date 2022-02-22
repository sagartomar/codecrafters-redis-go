package main

import (
	"fmt"
	"net"
	"os"
)

const (
    PONG string = "PONG"
    PLUS string = "+"
    CRLF string = "\r\n"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
    }
    conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
    var bytesRead []byte
    _, err = conn.Read(bytesRead)
    if err != nil {
        fmt.Println("Error while reading the data", err)
        os.Exit(1)
    }
    stringRead := string(bytesRead)
    fmt.Println("Read:", stringRead)
    reply := ConvertToRESPSimpleString(PONG)
    n, err := conn.Write([]byte(reply))
    if err != nil {
        fmt.Println("Error while writinh the data", err)
        os.Exit(1)
    }
    fmt.Printf("Wrote %d bytes", n)
}

func ConvertToRESPSimpleString(message string) string {
    return PLUS + message + CRLF
}
