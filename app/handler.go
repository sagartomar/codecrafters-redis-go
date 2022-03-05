package main

import (
	"io"
)

const PONG string = "PONG"

type Handler struct {
    writer io.Writer
}

func (h *Handler) Ping() error {
    _, err := h.writer.Write([]byte(ConvertToRESPSimpleString(PONG)))
    return err
}

func (h *Handler) Echo(message string) error {
    _, err := h.writer.Write([]byte(ConvertToRESPBulkString(message)))
    return err
}
