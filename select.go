package promptui

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"text/template"

	"github.com/chzyer/readline"
	"github.com/juju/ansiterm"
	"github.com/manifoldco/promptui/screenbuf"
)

// SelectedAdd is returned from SelectWithAdd when add is selected.
const SelectedAdd = -1

// Select represents a list for selecting a single item
type Select struct {
	// Label is the value displayed on the command line prompt. It can be any
	// value one would pass to a text/template Execute(), including just a string.
	Label interface{}

	// Items are the items to use in the list. It can be any slice type one would
	// pass to a text/template execute, including a string slice.
	Items interface{}

	// Size is the number of items that should appear on the select before
	// scrolling. If it is 0, defaults to 5.
	Size int

	// IsVimMode sets whether readline is using Vim mode.
	IsVimMode bool

	// Templates can be used to customize the select output. If nil is passed, the
	// default templates are used.
	Templates *SelectTemplates

	label string
	items []interface{}
}

// SelectTemplates allow a select prompt to be customized following stdlib
// text/template syntax. If any field is blank a default template is used.
type SelectTemplates struct {
	// Active is a text/template for the label.
	Label string

	// Active is a text/template for when an item is current active.
	Active string

	// Inactive is a text/template for when an item is not current active.
	Inactive string

	// Selected is a text/template for when an item was successfully selected.
	Selected string

	// Details is a text/template for when an item current active to show
	// additional information. It can have multiple lines.
	Details string

	// FuncMap is a map of helpers for the templates. If nil, the default helpers
	// are used.
	FuncMap template.FuncMap

	label    *template.Template
	active   *template.Template
	inactive *template.Template
	selected *template.Template
	details  *template.Template
}

// Run runs the Select list. It returns the index of the selected element,
// and its value.
func (s *Select) Run() (int, string, error) {
	if s.Size == 0 {
		s.Size = 5
	}

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
	end := s.listHeight()
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
	sb := screenbuf.New()

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
			start, end, selected = s.pageup(start, end, selected, max)
		case ' ': // space press
			start, end, selected = s.pagedown(start, end, selected, max)
		}

		label := renderBytes(s.Templates.label, s.Label)
		sb.Write(label)

		for i := start; i <= end; i++ {
			page := " "
			item := s.items[i]

			switch i {
			case 0:
				page = string(top)
			case max:
			case start:
				page = "↑"
			case end:
				page = "↓"
			}

			output := []byte(page + " ")

			if i == selected {
				output = append(output, renderBytes(s.Templates.active, item)...)
			} else {
				output = append(output, renderBytes(s.Templates.inactive, item)...)
			}

			sb.Write(output)
		}

		details := s.detailsOutput(selected)
		for _, d := range details {
			sb.Write(d)
		}

		sb.WriteTo(rl)
		rl.Refresh()

		return nil, 0, true
	})

	_, err = rl.Readline()

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

	item := s.items[selected]

	output := renderBytes(s.Templates.selected, item)

	sb.Reset()
	sb.Write(output)
	sb.WriteTo(rl)

	rl.Write([]byte(showCursor))
	rl.Close()

	return selected, fmt.Sprintf("%v", item), err
}

func (s *Select) prepareTemplates() error {
	if s.Items == nil || reflect.TypeOf(s.Items).Kind() != reflect.Slice {
		return fmt.Errorf("Items %v is not a slice", s.Items)
	}

	tpls := s.Templates
	if tpls == nil {
		tpls = &SelectTemplates{}
	}

	if tpls.FuncMap == nil {
		tpls.FuncMap = FuncMap
	}

	if tpls.Label == "" {
		tpls.Label = fmt.Sprintf("%s {{.}}: ", IconInitial)
	}

	tpl, err := template.New("").Funcs(tpls.FuncMap).Parse(tpls.Label)
	if err != nil {
		return err
	}

	tpls.label = tpl

	if tpls.Active == "" {
		tpls.Active = fmt.Sprintf("%s {{ . | underline }}", IconSelect)
	}

	tpl, err = template.New("").Funcs(tpls.FuncMap).Parse(tpls.Active)
	if err != nil {
		return err
	}

	tpls.active = tpl

	if tpls.Inactive == "" {
		tpls.Inactive = "  {{.}}"
	}

	tpl, err = template.New("").Funcs(tpls.FuncMap).Parse(tpls.Inactive)
	if err != nil {
		return err
	}

	tpls.inactive = tpl

	if tpls.Selected == "" {
		tpls.Selected = fmt.Sprintf(`{{ "%s" | green }} {{ . | faint }}`, IconGood)
	}

	tpl, err = template.New("").Funcs(tpls.FuncMap).Parse(tpls.Selected)
	if err != nil {
		return err
	}
	tpls.selected = tpl

	if tpls.Details != "" {
		tpl, err = template.New("").Funcs(tpls.FuncMap).Parse(tpls.Details)
		if err != nil {
			return err
		}

		tpls.details = tpl
	}

	list := reflect.ValueOf(s.Items)
	for i := 0; i < list.Len(); i++ {
		s.items = append(s.items, list.Index(i))
	}

	s.Templates = tpls

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
			Size:      5,
		}

		err := s.prepareTemplates()
		if err != nil {
			return 0, "", err
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

func (s *Select) pagedown(start, end, selected, max int) (newStart, newEnd, newSelected int) {
	newEnd = end + s.listHeight()

	if newEnd > max {
		newEnd = max
	}

	newStart = newEnd - s.listHeight()

	if newStart < 0 {
		newStart = 0
	}

	newSelected = newStart

	if newEnd == end {
		newSelected = newEnd
	} else if newSelected < selected {
		newSelected = selected
	}

	return newStart, newEnd, newSelected
}

func (s *Select) pageup(start, end, selected, max int) (newStart, newEnd, newSelected int) {
	newStart = start - s.listHeight()

	if newStart < 0 {
		newStart = 0
	}

	newEnd = newStart + s.listHeight()

	if newEnd > max {
		newEnd = max
	}

	newSelected = newStart

	if newSelected > selected {
		newSelected = selected
	}

	return newStart, newEnd, newSelected
}

func (s *Select) detailsOutput(idx int) [][]byte {
	if s.Templates.details == nil {
		return nil
	}

	var buf bytes.Buffer
	w := ansiterm.NewTabWriter(&buf, 0, 0, 8, ' ', 0)

	item := s.items[idx]
	err := s.Templates.details.Execute(w, item)
	if err != nil {
		fmt.Fprintf(w, "%v", item)
	}

	w.Flush()

	output := buf.Bytes()

	return bytes.Split(output, []byte("\n"))
}

func (s *Select) listHeight() int {
	if s.Size <= 0 {
		return 1
	}

	return s.Size - 1
}

func renderBytes(tpl *template.Template, data interface{}) []byte {
	var buf bytes.Buffer
	err := tpl.Execute(&buf, data)
	if err != nil {
		return []byte(fmt.Sprintf("%v", data))
	}
	return buf.Bytes()
}
