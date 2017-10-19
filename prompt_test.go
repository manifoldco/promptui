package promptui

import (
	"bytes"
	"testing"
)

func outputTest(mask rune, input, displayed, output, def string) func(t *testing.T) {
	return func(t *testing.T) {
		in := bytes.Buffer{}
		out := bytes.Buffer{}
		p := Prompt{
			Label:   "test",
			Default: def,
			Mask:    mask,
			stdin:   &in,
			stdout:  &out,
		}

		in.Write([]byte(input + "\n"))
		res, err := p.Run()

		if err != nil {
			t.Errorf("error during prompt: %s", err)
		}

		if res != output {
			t.Errorf("wrong result: %s != %s", res, output)
		}

		expected := "\x1b[1m\x1b[32m✔\x1b[0m \x1b[1mtest\x1b[0m\x1b[1m:\x1b[0m \x1b[2m" + displayed + "\033[0m\n"
		if !bytes.Equal(out.Bytes(), []byte(expected)) {
			t.Errorf("wrong output: %q != %q", out.Bytes(), expected)
		}

	}
}

func TestPrompt(t *testing.T) {
	t.Run("can read input", outputTest(0x0, "hi", "hi", "hi", ""))
	t.Run("displays masked values", outputTest('*', "hi", "**", "hi", ""))
	t.Run("can use a default", outputTest(0x0, "", "hi", "hi", "hi"))
}
