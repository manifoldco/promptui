package promptui

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

// Prompt represents a single line text field input.
type Prompt struct {
	Label string // Label is the value displayed on the command line prompt

	Default string // Default is the initial value to populate in the prompt

	// Validate is optional. If set, this function is used to validate the input
	// after each character entry.
	Validate ValidateFunc

	// If mask is set, this value is displayed instead of the actual input
	// characters.
	Mask rune

	IsConfirm bool
	IsVimMode bool
	Preamble  *string

	// Indent will be placed before the prompt's state symbol
	Indent string

	stdin  io.Reader
	stdout io.Writer
}

// Run runs the prompt, returning the validated input.
func (p *Prompt) Run() (string, error) {
	c := &readline.Config{}
	err := c.Init()
	if err != nil {
		return "", err
	}

	if p.stdin != nil {
		c.Stdin = p.stdin
	}

	if p.stdout != nil {
		c.Stdout = p.stdout
	}

	if p.Mask != 0 {
		c.EnableMask = true
		c.MaskRune = p.Mask
	}

	if p.IsVimMode {
		c.VimMode = true
	}

	if p.Preamble != nil {
		fmt.Println(*p.Preamble)
	}

	suggestedAnswer := ""
	punctuation := ":"
	if p.IsConfirm {
		punctuation = "?"
		answers := "y/N"
		if strings.ToLower(p.Default) == "y" {
			answers = "Y/n"
		}
		suggestedAnswer = " " + faint("["+answers+"]")
		p.Default = ""
	}

	state := IconInitial
	prompt := p.Label + punctuation + suggestedAnswer + " "

	c.Prompt = p.Indent + bold(state) + " " + bold(prompt)
	c.HistoryLimit = -1
	c.UniqueEditLine = true

	firstListen := true
	wroteErr := false
	caughtup := true
	var out string

	if p.Default != "" {
		caughtup = false
		out = p.Default
		c.Stdin = io.MultiReader(bytes.NewBuffer([]byte(out)), os.Stdin)
	}

	rl, err := readline.NewEx(c)
	if err != nil {
		return "", err
	}

	validFn := func(x string) error {
		return nil
	}

	if p.Validate != nil {
		validFn = p.Validate
	}

	c.SetListener(func(line []rune, pos int, key rune) ([]rune, int, bool) {
		if key == readline.CharEnter {
			return nil, 0, false
		}

		if firstListen {
			firstListen = false
			return nil, 0, false
		}

		if !caughtup && out != "" {
			if string(line) == out {
				caughtup = true
			}
			if wroteErr {
				return nil, 0, false
			}
		}

		err := validFn(string(line))
		if err != nil {
			if _, ok := err.(*ValidationError); ok {
				state = IconBad
			} else {
				rl.Close()
				return nil, 0, false
			}
		} else {
			state = IconGood
			if p.IsConfirm {
				state = IconInitial
			}
		}

		rl.SetPrompt(p.Indent + bold(state) + " " + bold(prompt))
		rl.Refresh()
		wroteErr = false

		return nil, 0, false
	})

	for {
		out, err = rl.Readline()

		var msg string
		valid := true
		oerr := validFn(out)
		if oerr != nil {
			if verr, ok := oerr.(*ValidationError); ok {
				msg = verr.msg
				valid = false
				state = IconBad
			} else {
				return "", oerr
			}
		}

		if valid {
			state = IconGood
			break
		}

		if err != nil {
			switch err {
			case readline.ErrInterrupt:
				err = ErrInterrupt
			case io.EOF:
				err = ErrEOF
			}

			break
		}

		caughtup = false

		c.Stdin = io.MultiReader(bytes.NewBuffer([]byte(out)), os.Stdin)
		rl, _ = readline.NewEx(c)

		firstListen = true
		wroteErr = true
		rl.SetPrompt("\n" + red(">> ") + msg + upLine(1) + "\r" + p.Indent + bold(state) + " " + bold(prompt))
		rl.Refresh()
	}

	if wroteErr {
		rl.Write([]byte(downLine(1) + clearLine + upLine(1) + "\r"))
	}

	if err != nil {
		if err.Error() == "Interrupt" {
			err = ErrInterrupt
		}
		rl.Write([]byte("\n"))
		return "", err
	}

	echo := out
	if p.Mask != 0 {
		echo = strings.Repeat(string(p.Mask), len(echo))
	}

	if p.IsConfirm {
		if strings.ToLower(echo) != "y" {
			state = IconBad
			err = ErrAbort
		} else {
			state = IconGood
		}
	}

	rl.Write([]byte(p.Indent + state + " " + prompt + faint(echo) + "\n"))

	return out, err
}
