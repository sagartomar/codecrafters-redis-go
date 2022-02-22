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

    stringRead, err := ReadFromConn(conn)
    if err != nil {
        fmt.Println("Error while reading the data", err)
        os.Exit(1)
    }
    fmt.Println("Read:", *stringRead)

    reply := ConvertToRESPSimpleString(PONG)
    err = WriteToConn(conn, reply)
    if err != nil {
        fmt.Println("Error while writing the data", err)
        os.Exit(1)
    }
}

func ReadFromConn(conn net.Conn) (*string, error) {
    var bytesRead []byte
    n, err := conn.Read(bytesRead)
    if err != nil {
        return nil, err
    }
    fmt.Printf("Read %d bytes", n)
    stringRead := string(bytesRead)
    return &stringRead, nil
}

func WriteToConn(conn net.Conn, message string) error {
    n, err := conn.Write([]byte(message))
    if err != nil {
        return err
    }
    fmt.Printf("Wrote %d bytes", n)
    return nil
}

func ConvertToRESPSimpleString(message string) string {
    return PLUS + message + CRLF
}
