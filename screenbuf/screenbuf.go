package screenbuf

import (
	"bytes"
	"fmt"
	"io"
)

const esc = "\033["

var (
	clearLine = []byte(esc + "2K\r")
	moveUp    = []byte(esc + "1A")
	moveDown  = []byte(esc + "1B")
)

// ScreenBuf is a convenient way to write to terminal screens. It creates,
// clears and, moves up or down lines as needed to write the output to the
// terminal using ANSI escape codes.
type ScreenBuf struct {
	buf    *bytes.Buffer
	lines  [][]byte
	reset  bool
	flush  bool
	cursor int
}

// New creates and initializes a new ScreenBuf.
func New() *ScreenBuf {
	return &ScreenBuf{buf: &bytes.Buffer{}}
}

// Reset truncates the underlining buffer and marks all its previous lines to be
// cleared during the next Write.
func (s *ScreenBuf) Reset() {
	s.buf.Reset()
	s.reset = true
}

// Write writes a single line to the underlining buffer. If the ScreenBuf was
// previously reset, all previous lines are cleared and the output starts from
// the top. Lines with \r or \n will fail since they can interfere with the
// terminal ability to move between lines.
func (s *ScreenBuf) Write(b []byte) (int, error) {
	if bytes.ContainsAny(b, "\r\n") {
		return 0, fmt.Errorf("%q should not contain either \\r or \\n", b)
	}

	if s.reset {
		for i := 1; i < len(s.lines); i++ {
			_, err := s.buf.Write(clearLine)
			if err != nil {
				return 0, err
			}
			_, err = s.buf.Write(moveUp)
			if err != nil {
				return 0, err
			}
		}
		s.cursor = 0
		s.lines = nil
		s.reset = false
	}

	if len(s.lines) <= s.cursor {
		s.lines = append(s.lines, b)
	} else {
		s.lines[s.cursor] = b
	}

	s.cursor++

	return len(b), nil
}

// WriteString is a convenient function to write a new line passing a string.
// Check ScreenBuf.Write() for a detailed explanation of the function behaviour.
func (s *ScreenBuf) WriteString(str string) (int, error) {
	return s.Write([]byte(str))
}

// WriteTo writes data to w until the buffer is drained or an error occurs. The
// return value n is the number of bytes written; it always fits into an int,
// but it is int64 to match the io.WriterTo interface. Any error encountered
// during the write is also returned.
func (s *ScreenBuf) WriteTo(w io.Writer) (int64, error) {
	if s.flush {
		for range s.lines {
			_, err := s.buf.Write(moveUp)
			if err != nil {
				return 0, err
			}
		}
	}

	for i, line := range s.lines {
		if s.flush {
			_, err := s.buf.Write(clearLine)
			if err != nil {
				return 0, err
			}

			if i < s.cursor {
				_, err = s.buf.Write(line)
				if err != nil {
					return 0, err
				}
			}

			_, err = s.buf.Write(moveDown)
			if err != nil {
				return 0, err
			}
		} else {
			l := append(line, []byte("\n")...)

			_, err := s.buf.Write(l)
			if err != nil {
				return 0, err
			}
		}
	}

	s.flush = true
	s.cursor = 0

	return s.buf.WriteTo(w)
}
