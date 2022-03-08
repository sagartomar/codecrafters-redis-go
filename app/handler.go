package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

const PONG string = "PONG"

type Handler struct {
    conn net.Conn
    reader *bufio.Reader
    writer *bufio.Writer
}

func NewHandler(conn net.Conn) *Handler {
    if conn == nil {
        return nil
    }
    h := Handler {conn: conn}
    h.reader = bufio.NewReader(h.conn)
    h.writer = bufio.NewWriter(h.conn)
    return &h
}

func (h *Handler) HandleConnection() {
    for {
        h.ReadRESPArray()
        h.Ping()
    }
}

func (h *Handler) Ping() error {
    _, err := h.writer.WriteString(ConvertToRESPSimpleString(PONG))
    h.writer.Flush()
    return err
}

func (h *Handler) Echo(message string) error {
    _, err := h.writer.WriteString(ConvertToRESPBulkString(message))
    h.writer.Flush()
    return err
}

func (h *Handler) ReadRESPBulkString() (string, error) {
    input, err := h.reader.ReadString('\n')
    if err != nil {
        return "", err
    }
    input = RemoveLF(input)
    input = RemoveCR(input)
    if input[0] != '$' {
        return "", fmt.Errorf("Expected '$' but received %c as first character", input[0])
    }
    length, err := strconv.Atoi(input[1:])
    if err != nil {
        return "", err
    }
    input, err = h.reader.ReadString('\n')
    if err != nil {
        return "", err
    }
    input = RemoveLF(input)
    input = RemoveCR(input)
    if len(input) != length {
        return "", fmt.Errorf("Expected length %d but received length %d", length, len(input))
    }
    return input, nil
}

func (h *Handler) ReadRESPArray() ([]string, error) {
    input, err := h.reader.ReadString('\n')
    if err != nil {
        return nil, err
    }
    input = RemoveLF(input)
    input = RemoveCR(input)
    if input[0] != '*' {
        return nil, fmt.Errorf("Expected '*' but received %c as first character", input[0])
    }
    length, err := strconv.Atoi(input[1:])
    if err != nil {
        return nil, err
    }
    arr := make([]string, length)
    for i := 0; i < length; i++ {
        element, err := h.ReadRESPBulkString()
        if err != nil {
            return nil, err
        }
        arr[i] = element
    }
    return arr, nil
}
