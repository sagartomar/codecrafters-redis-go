package main

import "fmt"

const (
	PLUS string = "+"
	CRLF string = "\r\n"
)

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
