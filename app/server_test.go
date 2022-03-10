package main

import "testing"

func TestRESPSimpleString(t *testing.T) {
	input := "OK"
	expected := "+OK\r\n"

	AssertStringEqual(t, ConvertToRESPSimpleString(input), expected)
}

func TestRESPBulkString(t *testing.T) {
	input := "test"
	expected := "$4\r\ntest\r\n"

	AssertStringEqual(t, ConvertToRESPBulkString(input), expected)
}

func TestRemoveCR(t *testing.T) {

	t.Run("Carriage return should be removed if present as the last character", func(t *testing.T) {
		input := "test\r"
		expected := "test"

		AssertStringEqual(t, RemoveCR(input), expected)
	})

	t.Run("If carriage return is not last character then input string should be returned", func(t *testing.T) {
		input := "test"
		expected := input

		AssertStringEqual(t, RemoveCR(input), expected)
	})

}

func TestRemoveLF(t *testing.T) {

	t.Run("Linefeed should be removed if present as the last character", func(t *testing.T) {
		input := "test\n"
		expected := "test"

		AssertStringEqual(t, RemoveLF(input), expected)
	})

	t.Run("If linefeed is not last character then input string should be returned", func(t *testing.T) {
		input := "test"
		expected := input

		AssertStringEqual(t, RemoveLF(input), expected)
	})

}

func AssertStringEqual(t testing.TB, received, expected string) {
	t.Helper()

	if received != expected {
		t.Errorf("Expected %s but received %s", expected, received)
	}
}
