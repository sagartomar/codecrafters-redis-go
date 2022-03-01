package main

import (
	"fmt"
	"net"
	"os"
    "bufio"
    "strconv"
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

    for {
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting connection: ", err.Error())
            os.Exit(1)
        }   

        go HandleConnection(conn)
    }
}

func HandleConnection(conn net.Conn) {
    reader := bufio.NewReader(conn)
    defer conn.Close()

    for {
        input := ProcessInput(reader)
        if input != nil {
            for _, in := range input {
                fmt.Println(in)
            }
        }
        ExecuteCommand(input, conn)
        reply := ConvertToRESPSimpleString(PONG)
        err := WriteToConn(conn, reply)
        if err != nil {
            fmt.Println("Error while writing the data", err.Error())
            break
        }
    }
}

func ProcessInput(reader *bufio.Reader) []string{
    input, err := reader.ReadString('\r')
    if err != nil {
        fmt.Println("Error while reading data:", err.Error())
        return nil
    }
    input = RemoveCR(input)
    input = RemoveLF(input)
    if input[0] != '*' {
        fmt.Println("Invalid character:", input[0])
        return nil
    }
    arrayLen, err := strconv.Atoi(input[1 : ])
    if err != nil {
        fmt.Println("Error while converting array length to int", err.Error())
    }
    arr := make([]string, arrayLen)
    for i := 0; i < arrayLen; i++ {
        data := ProcessElement(reader)
        if data == nil {
            fmt.Println("Error in reading element")
            return nil
        }
        arr[i] = *data
    }
    return arr
}

func ProcessElement(reader *bufio.Reader) *string {
    input, err := reader.ReadString('\r')
    if err != nil {
        fmt.Println("Error while reading element size:", err.Error())
        return nil
    }
    input = RemoveCR(input)
    input = RemoveLF(input)
    if input[0] != '$' {
        fmt.Println("Invalid character:", input[0])
        return nil
    }
    elementLength, err := strconv.Atoi(input[1 : ])
    if err != nil {
        fmt.Println("Error while converting element length to int:", err.Error())
        return nil
    }
    fmt.Println(elementLength)
    input, err = reader.ReadString('\r')
    if err != nil {
        fmt.Println("Error while reading element:", err.Error())
        return nil
    }
    input = RemoveCR(input)
    input = RemoveLF(input)
    if len(input) != elementLength {
        fmt.Printf("Size mismatch: %d != %d\n", elementLength, len(input))
        return nil
    }
    return &input
}

func ExecuteCommand(input []string, conn net.Conn) {
    switch(input[0]) {
    case "PING":
        return
    case "ECHO":
        Echo(conn, input[1 :])
    }
}

func Ping(conn net.Conn) {
}

func Echo(conn net.Conn, arguments []string) {
    if len(arguments) != 1 {
        fmt.Printf("ECHO expects 1 argument but received %d\n", len(arguments))
        return
    }
    err := WriteToConn(conn, arguments[0])
    if err != nil {
        fmt.Println("Error while sending the message:", err.Error())
    }
}

func RemoveCR(input string) string {
    if input[len(input) - 1] == '\r' {
        return input[0 : len(input) - 1]
    }
    return input
}

func RemoveLF(input string) string {
    if input[0] == '\n' {
        return input[1 : ]
    }
    return input
}

func WriteToConn(conn net.Conn, message string) error {
    n, err := conn.Write([]byte(message))
    if err != nil {
        return err
    }
    fmt.Printf("Wrote %d bytes\n", n)
    return nil
}

func ConvertToRESPSimpleString(message string) string {
    return PLUS + message + CRLF
}
