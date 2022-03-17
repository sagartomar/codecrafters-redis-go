package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

const PONG string = "PONG"

type Handler struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewHandler(conn net.Conn) *Handler {
	if conn == nil {
		return nil
	}
	h := Handler{conn: conn}
	h.reader = bufio.NewReader(h.conn)
	h.writer = bufio.NewWriter(h.conn)
	return &h
}

func (h *Handler) HandleConnection() {
	for {
		array, err := h.ReadRESPArray()
		if err == io.EOF {
			h.conn.Close()
			break
		}
		if array != nil {
			switch strings.ToUpper(array[0]) {
			case "PING":
				h.Ping()
			case "ECHO":
				h.Echo(array)
			}
		}
	}
}

func (h *Handler) Ping() error {
	_, err := h.writer.WriteString(ConvertToRESPSimpleString(PONG))
	h.writer.Flush()
	return err
}

func (h *Handler) Echo(arguments []string) error {
    if len(arguments) != 2 {
        return fmt.Errorf("Expected 2 arguments but received %d", len(arguments))
    }
    if strings.ToUpper(arguments[0]) != "ECHO" {
        return fmt.Errorf("Expected 'ECHO' but received %s", arguments[0])
    }
	_, err := h.writer.WriteString(ConvertToRESPBulkString(arguments[1]))
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
