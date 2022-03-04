package main

import "testing"

func TestRESPSimpleString(t *testing.T) {
    input := "OK"
    expected := "+OK\r\n"

    output := ConvertToRESPSimpleString(input)

    if output != expected {
        t.Errorf("Expected %s but received %s", expected, output)
    }
}

func TestRESPBulkString(t *testing.T) {
    input := "test"
    expected := "$4\r\ntest\r\n"

    output := ConvertToRESPBulkString(input)

    if output != expected {
        t.Errorf("Expected %s but received %s", expected, output)
    }
}

func TestRemoveCR(t *testing.T) {

    t.Run("Carriage return should be removed if present as the last character", func(t *testing.T) {
        input := "test\r"
        expected := "test"

        output := RemoveCR(input)

        if output != expected {
            t.Errorf("Expected %s but received %s", expected, output)
        }
    })

    t.Run("If carriage return is not last character then input string should be returned", func(t *testing.T) {
        input := "test"
        expected := input

        output := RemoveCR(input)

        if output != expected {
            t.Errorf("Expected %s but received %s", expected, output)
        }
    })

}

func TestRemoveLF(t *testing.T) {

    t.Run("Linefeed should be removed if present as the last character", func(t *testing.T) {
        input := "test\n"
        expected := "test"

        output := RemoveLF(input)

        if output != expected {
            t.Errorf("Expected %s but received %s", expected, output)
        }
    })

    t.Run("If linefeed is not last character then input string should be returned", func(t *testing.T) {
        input := "test"
        expected := input

        output := RemoveLF(input)

        if output != expected {
            t.Errorf("Expected %s but received %s", expected, output)
        }
    })

}
