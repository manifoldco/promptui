package screenbuf

import (
	"bytes"
	"testing"
)

func TestScreen(t *testing.T) {
	// overwrite regular movement codes for easier visualization
	clearLine = []byte("\\c")
	moveUp = []byte("\\u")
	moveDown = []byte("\\d")

	s := New()

	t.Run("initial flush", func(t *testing.T) {
		s.Write([]byte("Hello Darkness,"))
		s.Write([]byte("My Old Friend"))
		s.Write([]byte("I've come to talk with you again"))

		var buf bytes.Buffer

		_, err := s.WriteTo(&buf)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		expect := `Hello Darkness,
My Old Friend
I've come to talk with you again
`

		got := buf.String()

		if expect != got {
			t.Errorf("expected %s, got %s", expect, got)
		}
	})

	t.Run("flush with same amount of lines", func(t *testing.T) {
		s.Write([]byte("Because a vision softly creeping"))
		s.Write([]byte("Left it's seeds while I was sleeping"))
		s.Write([]byte("And the vision that was planted"))

		var buf bytes.Buffer
		_, err := s.WriteTo(&buf)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		expect := `\u\u\u\cBecause a vision softly creeping\d\cLeft it's seeds while I was sleeping\d\cAnd the vision that was planted\d`

		got := buf.String()

		if expect != got {
			t.Errorf("expected:\n%s\ngot:\n%s", expect, got)
		}
	})

	t.Run("flush with less lines", func(t *testing.T) {
		s.Write([]byte("In my brain still remains"))
		s.Write([]byte("Within the sound of silence"))

		var buf bytes.Buffer
		_, err := s.WriteTo(&buf)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		expect := `\u\u\cIn my brain still remains\d\cWithin the sound of silence\d\c\d`

		got := buf.String()

		if expect != got {
			t.Errorf("expected:\n%s\ngot:\n%s", expect, got)
		}
	})

	t.Run("flush with more lines", func(t *testing.T) {
		s.Write([]byte("In restless dreams I walked alone"))
		s.Write([]byte("Narrow streets of cobblestone"))
		s.Write([]byte("'Neath the halo of a street lamp"))
		s.Write([]byte("I turned my collar to the cold and damp"))
		s.Write([]byte("When my eyes were stabbed by the flash of"))
		s.Write([]byte("A neon light that split the night"))
		s.Write([]byte("And touched the sound of silence"))

		var buf bytes.Buffer
		_, err := s.WriteTo(&buf)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		expect := `\u\u\u\u\u\u\u\cIn restless dreams I walked alone\d\cNarrow streets of cobblestone\d\c'Neath the halo of a street lamp\d\cI turned my collar to the cold and damp\d\cWhen my eyes were stabbed by the flash of\d\cA neon light that split the night\d\cAnd touched the sound of silence\d`

		got := buf.String()

		if expect != got {
			t.Errorf("expected:\n%s\ngot:\n%s", expect, got)
		}
	})

	t.Run("reset with fewer lines", func(t *testing.T) {
		s.Reset()
		s.Write([]byte("The Sound of Silence"))

		var buf bytes.Buffer
		_, err := s.WriteTo(&buf)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		expect := `\c\u\c\u\c\u\c\u\c\u\c\u\u\cThe Sound of Silence\d`

		got := buf.String()

		if expect != got {
			t.Errorf("expected:\n%s\ngot:\n%s", expect, got)
		}
	})
}
