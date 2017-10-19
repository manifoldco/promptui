package promptui

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/chzyer/readline"
)

// SelectedAdd is returned from SelectWithAdd when add is selected.
const SelectedAdd = -1

// TODO allow custom select height
const pagination = 4

// FuncMap defines template helpers for the output. It can be extended as a
// regular map.
var FuncMap = template.FuncMap{
	"red":        red,
	"bold":       bold,
	"faint":      faint,
	"underlined": underlined,
}

// Select represents a list for selecting a single item
type Select struct {
	// Label is the value displayed on the command line prompt. It can be any
	// value one would pass to a text/template Execute(), including just a string.
	Label interface{}

	// LabelTemplate is a text template for label. It can be blank if the label is
	// a string.
	LabelTemplate string

	// Items are the items to use in the list. It can be any slice value one would
	// pass to a text/template execute, including a string slice.
	Items interface{}

	// ItemsTemplate is a text template for each item. It can be blank if items is
	// a string slice.
	ItemsTemplate string

	// FuncMap is a map of helpers for the templates. If nil, the default helpers
	// are used.
	FuncMap template.FuncMap

	// IsVimMode sets whether readline is using Vim mode.
	IsVimMode bool

	label string
	items []string
}

// Run runs the Select list. It returns the index of the selected element,
// and its value.
func (s *Select) Run() (int, string, error) {
	err := s.prepareTemplates()
	if err != nil {
		return 0, "", err
	}
	return s.innerRun(0, ' ')
}

func (s *Select) innerRun(starting int, top rune) (int, string, error) {
	stdin := readline.NewCancelableStdin(os.Stdin)
	c := &readline.Config{}
	err := c.Init()
	if err != nil {
		return 0, "", err
	}

	c.Stdin = stdin

	if s.IsVimMode {
		c.VimMode = true
	}

	c.HistoryLimit = -1
	c.UniqueEditLine = true

	start := 0
	end := 4
	max := len(s.items) - 1

	if len(s.items) <= end {
		end = max
	}

	selected := starting

	rl, err := readline.NewEx(c)
	if err != nil {
		return 0, "", err
	}

	rl.Write([]byte(hideCursor))
	rl.Write([]byte(strings.Repeat("\n", end-start+1)))

	counter := 0

	rl.Operation.ExitVimInsertMode() // Never use insert mode for selects

	c.SetListener(func(line []rune, pos int, key rune) ([]rune, int, bool) {
		if rl.Operation.IsEnableVimMode() {
			rl.Operation.ExitVimInsertMode()
			// Remap j and k for down/up selections immediately after an
			// `i` press
			switch key {
			case 'j':
				key = readline.CharNext
			case 'k':
				key = readline.CharPrev
			}
		}

		switch key {
		case readline.CharEnter:
			return nil, 0, true
		case readline.CharNext:
			switch selected {
			case max:
			case end:
				start++
				end++
				fallthrough
			default:
				selected++
			}
		case readline.CharPrev:
			switch selected {
			case 0:
			case start:
				start--
				end--
				fallthrough
			default:
				selected--
			}
		case 'b':
			start, end, selected = pageup(start, end, selected, max)
		case ' ': // space press
			start, end, selected = pagedown(start, end, selected, max)
		}

		list := make([]string, end-start+1)
		for i := start; i <= end; i++ {
			page := ' '
			selection := " "
			item := s.items[i]

			switch i {
			case 0:
				page = top
			case len(s.items) - 1:
			case start:
				page = '↑'
			case end:
				page = '↓'
			}
			if i == selected {
				selection = "▸"
				item = underlined(item)
			}
			list[i-start] = clearLine + "\r" + string(page) + " " + selection + " " + item
		}

		prefix := ""
		prefix += upLine(uint(len(list))) + "\r" + clearLine
		p := prefix + bold(IconInitial) + " " + bold(s.label) + downLine(1) + strings.Join(list, downLine(1))
		rl.SetPrompt(p)
		rl.Refresh()

		counter++

		return nil, 0, true
	})

	_, err = rl.Readline()
	rl.Close()

	if err != nil {
		switch {
		case err == readline.ErrInterrupt, err.Error() == "Interrupt":
			err = ErrInterrupt
		case err == io.EOF:
			err = ErrEOF
		}

		rl.Write([]byte("\n"))
		rl.Write([]byte(showCursor))
		rl.Refresh()
		return 0, "", err
	}

	rl.Write(bytes.Repeat([]byte(clearLine+upLine(1)), end-start+1))
	rl.Write([]byte("\r"))

	out := s.items[selected]

	rl.Write([]byte(IconGood + " " + s.label + faint(out) + "\n"))

	rl.Write([]byte(showCursor))
	return selected, out, err
}

func (s *Select) prepareTemplates() error {
	if reflect.TypeOf(s.Items).Kind() != reflect.Slice {
		fmt.Errorf("Items is not a slice")
	}

	if s.LabelTemplate == "" {
		s.LabelTemplate = "{{.}}:"
	}

	if s.ItemsTemplate == "" {
		s.ItemsTemplate = "{{.}}"
	}

	if s.FuncMap == nil {
		s.FuncMap = FuncMap
	}

	tpl, err := template.New("label").Funcs(s.FuncMap).Parse(s.LabelTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, s.Label)
	if err != nil {
		return err
	}

	s.label = buf.String()

	tpl, err = template.New("items").Funcs(s.FuncMap).Parse(s.ItemsTemplate)
	if err != nil {
		return err
	}

	list := reflect.ValueOf(s.Items)
	for i := 0; i < list.Len(); i++ {
		var buf bytes.Buffer
		err := tpl.Execute(&buf, list.Index(i))
		if err != nil {
			return err
		}
		s.items = append(s.items, buf.String())
	}

	return nil
}

// SelectWithAdd represents a list for selecting a single item, or selecting
// a newly created item.
type SelectWithAdd struct {
	Label string   // Label is the value displayed on the command line prompt.
	Items []string // Items are the items to use in the list.

	AddLabel string // The label used in the item list for creating a new item.

	// Validate is optional. If set, this function is used to validate the input
	// after each character entry.
	Validate ValidateFunc

	IsVimMode bool // Whether readline is using Vim mode.
}

// Run runs the Select list. It returns the index of the selected element,
// and its value. If a new element is created, -1 is returned as the index.
func (sa *SelectWithAdd) Run() (int, string, error) {
	if len(sa.Items) > 0 {
		newItems := append([]string{sa.AddLabel}, sa.Items...)

		s := Select{
			Label:     sa.Label,
			Items:     newItems,
			IsVimMode: sa.IsVimMode,
		}

		selected, value, err := s.innerRun(1, '+')
		if err != nil || selected != 0 {
			return selected - 1, value, err
		}

		// XXX run through terminal for windows
		os.Stdout.Write([]byte(upLine(1) + "\r" + clearLine))
	}

	p := Prompt{
		Label:     sa.AddLabel,
		Validate:  sa.Validate,
		IsVimMode: sa.IsVimMode,
	}
	value, err := p.Run()
	return SelectedAdd, value, err
}

func pagedown(start, end, selected, max int) (newStart, newEnd, newSelected int) {
	newEnd = end + pagination

	if newEnd > max {
		newEnd = max
	}

	newStart = newEnd - pagination

	if newStart < 0 {
		newStart = 0
	}

	newSelected = newStart

	if newSelected < selected {
		newSelected = selected
	}

	return newStart, newEnd, newSelected
}

func pageup(start, end, selected, max int) (newStart, newEnd, newSelected int) {
	newStart = start - pagination

	if newStart < 0 {
		newStart = 0
	}

	newEnd = newStart + pagination

	if newEnd > max {
		newEnd = max
	}

	newSelected = newStart

	if newSelected > selected {
		newSelected = selected
	}

	return newStart, newEnd, newSelected
}
