package main

import (
	"fmt"
	"net"
	"os"
)

const (
	PLUS string = "+"
	CRLF string = "\r\n"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

    kv := NewInMemoryKV()

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

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

func RemoveCR(input string) string {
	if input[len(input)-1] == '\r' {
		return input[0 : len(input)-1]
	}
	return input
}

func RemoveLF(input string) string {
	if input[len(input)-1] == '\n' {
		return input[0 : len(input)-1]
	}
	return input
}

func ConvertToRESPSimpleString(message string) string {
	return PLUS + message + CRLF
}

func ConvertToRESPBulkString(message string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(message), message)
}
