package promptui

import (
	"errors"
	"io"
	"strings"
	"testing"
)

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }

func TestPromptRace(t *testing.T) {
	for i := 0; i < 200; i++ {
		p := Prompt{
			Label:  "Test",
			Stdin:  io.NopCloser(strings.NewReader("hello\r\n")),
			Stdout: nopWriteCloser{io.Discard},
		}
		res, err := p.Run()
		if err != nil && err != io.EOF {
			t.Fatalf("unexpected error: %v", err)
		}
		if res != "hello" {
			t.Fatalf("expected hello, got %q", res)
		}
	}
}

func TestPromptRun(t *testing.T) {
	t.Run("default value on enter", func(t *testing.T) {
		p := Prompt{
			Label:   "Username",
			Default: "admin",
			Stdin:   io.NopCloser(strings.NewReader("\r\n")),
			Stdout:  nopWriteCloser{io.Discard},
		}
		res, err := p.Run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res != "admin" {
			t.Fatalf("expected admin, got %q", res)
		}
	})

	t.Run("validation retry", func(t *testing.T) {
		attempts := 0
		p := Prompt{
			Label: "Number",
			Validate: func(input string) error {
				attempts++
				if input == "bad" {
					return errors.New("invalid input")
				}
				return nil
			},
			Stdin:  io.NopCloser(strings.NewReader("bad\r\n\b\b\bgood\r\n")),
			Stdout: nopWriteCloser{io.Discard},
		}
		res, err := p.Run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res != "good" {
			t.Fatalf("expected good, got %q", res)
		}
	})

	t.Run("confirmation prompt", func(t *testing.T) {
		p := Prompt{
			Label:     "Continue",
			IsConfirm: true,
			Stdin:     io.NopCloser(strings.NewReader("y\r\n")),
			Stdout:    nopWriteCloser{io.Discard},
		}
		res, err := p.Run()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res != "y" {
			t.Fatalf("expected y, got %q", res)
		}
	})
}
