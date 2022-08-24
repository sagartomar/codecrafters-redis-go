package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	PONG             string = "PONG"
	OK               string = "OK"
	NULL_BULK_STRING string = "$-1\r\n"
)

type Handler struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	store  *InMemoryKV
}

func NewHandler(conn net.Conn, store *InMemoryKV) *Handler {
	if conn == nil {
		return nil
	}
	h := Handler{conn: conn, store: store}
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
				h.Ping(array)
			case "ECHO":
				h.Echo(array)
			case "SET":
				h.Set(array)
			case "GET":
				h.Get(array)
			}
		}
	}
}

func (h *Handler) Ping(arguments []string) error {
	if len(arguments) != 1 {
		return fmt.Errorf("Expected 1 argument but received %d", len(arguments))
	}
	if strings.ToUpper(arguments[0]) != "PING" {
		return fmt.Errorf("Expected 'PING' but received %s", arguments[0])
	}
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

func (h *Handler) Set(arguments []string) {
	switch len(arguments) {
	case 3:
		h.store.Set(arguments[1], arguments[2])
	case 5:
		dur_arg, _ := strconv.Atoi(arguments[4])
		duration := time.Duration(dur_arg) * time.Millisecond
		h.store.SetWithExpiry(arguments[1], arguments[2], duration)
	}
	h.writer.WriteString(ConvertToRESPSimpleString(OK))
	h.writer.Flush()
}

func (h *Handler) Get(arguments []string) {
	err, value := h.store.Get(arguments[1])
	if err != nil {
		h.writer.WriteString(NULL_BULK_STRING)
		h.writer.Flush()
		return
	}
	h.writer.WriteString(ConvertToRESPBulkString(value))
	h.writer.Flush()
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
