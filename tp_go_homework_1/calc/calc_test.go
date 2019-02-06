package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

// сюда писать тесты

func TestCalcCorrect(t *testing.T) {
	var cases = []struct {
		expected string
		input    io.Reader
	} {
		{
			expected: "Result = 21",
			input: strings.NewReader("1 2 + 3 4 + * ="),
		},
		{
			expected: "Result = 2",
			input: strings.NewReader("4 2 / ="),
		},
		{
			expected: "Result = 15",
			input: strings.NewReader("1 2 3 4 + * + ="),
		},
	}

	for key, value := range cases {
		out := new(bytes.Buffer)
		if err := calc(value.input, out); err == nil {
			result := out.String()
			if result != value.expected {
				t.Errorf("test %d failed, wrong result: %s != %s", key, result, value.expected)
			}
		} else {
			t.Errorf("test %d failed, error: %s", key, err)
		}
	}
}

func TestCalcIncorrect(t *testing.T) {
	var cases = []struct {
		expected string
		input    io.Reader
	} {
		{
			expected: "stack is empty",
			input: strings.NewReader("="),
		},
		{
			expected: "stack is empty",
			input: strings.NewReader("1 +"),
		},
		{
			expected: "stack is empty",
			input: strings.NewReader("1 2 + + + + + +"),
		},
	}

	for key, value := range cases {
		out := new(bytes.Buffer)
		if err := calc(value.input, out); err == nil {
			t.Errorf("test %d failed, no error but expected", key)
		} else {
			if err.Error() != value.expected {
				t.Errorf("test %d failed, expect %s, got %s", key, value.expected, err.Error())
			}
		}
	}
}