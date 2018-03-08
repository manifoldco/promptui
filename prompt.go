package promptui

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/chzyer/readline"
	"github.com/manifoldco/promptui/screenbuf"
)

const cursor = "\u2588"

// Prompt represents a single line text field input.
type Prompt struct {
	// Label is the value displayed on the command line prompt. It can be any
	// value one would pass to a text/template Execute(), including just a string.
	Label interface{}

	Default string // Default is the initial value to populate in the prompt

	// AllowEdit lets the user edit the default value. If false, any key press
	// other than <Enter> automatically clears the default value.
	AllowEdit bool

	// Validate is optional. If set, this function is used to validate the input
	// after each character entry.
	Validate ValidateFunc

	// If mask is set, this value is displayed instead of the actual input
	// characters.
	Mask rune

	// Templates can be used to customize the prompt output. If nil is passed, the
	// default templates are used.
	Templates *PromptTemplates

	// IsConfirm sets the prompt to be a [y/N] question.
	IsConfirm bool

	// IsVimMode enables vi-like movements (hjkl) and editing.
	IsVimMode bool

	stdin  io.ReadCloser
	stdout io.WriteCloser
}

// PromptTemplates allow a prompt to be customized following stdlib
// text/template syntax. If any field is blank a default template is used.
type PromptTemplates struct {
	// Prompt is a text/template for the initial prompt question.
	Prompt string

	// Prompt is a text/template if the prompt is a confirmation.
	Confirm string

	// Valid is a text/template for when the current input is valid.
	Valid string

	// Invalid is a text/template for when the current input is invalid.
	Invalid string

	// Success is a text/template for the successful result.
	Success string

	// Prompt is a text/template when there is a validation error.
	ValidationError string

	// FuncMap is a map of helpers for the templates. If nil, the default helpers
	// are used.
	FuncMap template.FuncMap

	prompt     *template.Template
	valid      *template.Template
	invalid    *template.Template
	validation *template.Template
	success    *template.Template
}

// Run runs the prompt, returning the validated input.
func (p *Prompt) Run() (string, error) {
	c := &readline.Config{}
	err := c.Init()
	if err != nil {
		return "", err
	}

	err = p.prepareTemplates()
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

	c.HistoryLimit = -1
	c.UniqueEditLine = true

	rl, err := readline.NewEx(c)
	if err != nil {
		return "", err
	}

	rl.Write([]byte(hideCursor))
	sb := screenbuf.New(rl)

	validFn := func(x string) error {
		return nil
	}

	if p.Validate != nil {
		validFn = p.Validate
	}

	var inputErr error
	input := p.Default
	if p.IsConfirm {
		input = ""
	}
	eraseDefault := input != "" && !p.AllowEdit

	c.SetListener(func(line []rune, pos int, key rune) ([]rune, int, bool) {
		if line != nil {
			input += string(line)
		}

		switch key {
		case 0: // empty
		case KeyEnter:
			return nil, 0, false
		case KeyBackspace:
			if eraseDefault {
				eraseDefault = false
				input = ""
			}
			if len(input) > 0 {
				r := []rune(input)
				input = string(r[:len(r)-1])
			}
		default:
			if eraseDefault {
				eraseDefault = false
				input = string(line)
			}
		}

		err := validFn(input)
		var prompt []byte

		if err != nil {
			prompt = render(p.Templates.invalid, p.Label)
		} else {
			prompt = render(p.Templates.valid, p.Label)
			if p.IsConfirm {
				prompt = render(p.Templates.prompt, p.Label)
			}
		}

		echo := input
		if p.Mask != 0 {
			echo = strings.Repeat(string(p.Mask), len(echo))
		}

		prompt = append(prompt, []byte(echo+cursor)...)

		sb.Reset()
		sb.Write(prompt)

		if inputErr != nil {
			validation := render(p.Templates.validation, inputErr)
			sb.Write(validation)
			inputErr = nil
		}

		sb.Flush()

		return nil, 0, true
	})

	for {
		_, err = rl.Readline()

		inputErr = validFn(input)
		if inputErr == nil {
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
	}

	if err != nil {
		if err.Error() == "Interrupt" {
			err = ErrInterrupt
		}
		sb.Reset()
		sb.WriteString("")
		sb.Flush()
		rl.Write([]byte(showCursor))
		rl.Close()
		return "", err
	}

	echo := input
	if p.Mask != 0 {
		echo = strings.Repeat(string(p.Mask), len(echo))
	}

	prompt := render(p.Templates.valid, p.Label)
	prompt = append(prompt, []byte(echo)...)

	if p.IsConfirm {
		lowerDefault := strings.ToLower(p.Default)
		if strings.ToLower(echo) != "y" && (lowerDefault != "y" || (lowerDefault == "y" && echo != "")) {
			prompt = render(p.Templates.invalid, p.Label)
			err = ErrAbort
		}
	}

	sb.Reset()
	sb.Write(prompt)
	sb.Flush()
	rl.Write([]byte(showCursor))
	rl.Close()

	return input, err
}

func (p *Prompt) prepareTemplates() error {
	tpls := p.Templates
	if tpls == nil {
		tpls = &PromptTemplates{}
	}

	if tpls.FuncMap == nil {
		tpls.FuncMap = FuncMap
	}

	bold := Styler(FGBold)

	if p.IsConfirm {
		if tpls.Confirm == "" {
			confirm := "y/N"
			if strings.ToLower(p.Default) == "y" {
				confirm = "Y/n"
			}
			tpls.Confirm = fmt.Sprintf(`{{ "%s" | bold }} {{ . | bold }}? {{ "[%s]" | faint }} `, IconInitial, confirm)
		}

		tpl, err := template.New("").Funcs(tpls.FuncMap).Parse(tpls.Confirm)
		if err != nil {
			return err
		}

		tpls.prompt = tpl
	} else {
		if tpls.Prompt == "" {
			tpls.Prompt = fmt.Sprintf("%s {{ . | bold }}%s ", bold(IconInitial), bold(":"))
		}

		tpl, err := template.New("").Funcs(tpls.FuncMap).Parse(tpls.Prompt)
		if err != nil {
			return err
		}

		tpls.prompt = tpl
	}

	if tpls.Valid == "" {
		tpls.Valid = fmt.Sprintf("%s {{ . | bold }}%s ", bold(IconGood), bold(":"))
	}

	tpl, err := template.New("").Funcs(tpls.FuncMap).Parse(tpls.Valid)
	if err != nil {
		return err
	}

	tpls.valid = tpl

	if tpls.Invalid == "" {
		tpls.Invalid = fmt.Sprintf("%s {{ . | bold }}%s ", bold(IconBad), bold(":"))
	}

	tpl, err = template.New("").Funcs(tpls.FuncMap).Parse(tpls.Invalid)
	if err != nil {
		return err
	}

	tpls.invalid = tpl

	if tpls.ValidationError == "" {
		tpls.ValidationError = `{{ ">>" | red }} {{ . | red }}`
	}

	tpl, err = template.New("").Funcs(tpls.FuncMap).Parse(tpls.ValidationError)
	if err != nil {
		return err
	}

	tpls.validation = tpl

	if tpls.Success == "" {
		tpls.Success = `{{ . | faint }}`
	}

	tpl, err = template.New("").Funcs(tpls.FuncMap).Parse(tpls.Success)
	if err != nil {
		return err
	}

	tpls.success = tpl

	p.Templates = tpls

	return nil
}
