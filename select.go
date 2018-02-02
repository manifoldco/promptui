package promptui

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/chzyer/readline"
	"github.com/juju/ansiterm"
	"github.com/manifoldco/promptui/list"
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

	// Keys can be used to change movement and search keys.
	Keys *SelectKeys

	// Searcher can be implemented to teach the select how to search for items.
	Searcher list.Searcher

	// Starts the prompt in search mode.
	StartInSearchMode bool

	label string

	list *list.List
}

// SelectKeys defines which keys can be used for movement and search.
type SelectKeys struct {
	Next     Key // Next defaults to down arrow key
	Prev     Key // Prev defaults to up arrow key
	PageUp   Key // PageUp defaults to left arrow key
	PageDown Key // PageDown defaults to right arrow key
	Search   Key // Search defaults to '/' key
}

// Key defines a keyboard code and a display representation for the help
// Check https://github.com/chzyer/readline for a list of codes
type Key struct {
	Code    rune
	Display string
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

	// Help is a text/template for displaying instructions at the top. By default
	// it shows keys for movement and search.
	Help string

	// FuncMap is a map of helpers for the templates. If nil, the default helpers
	// are used.
	FuncMap template.FuncMap

	label    *template.Template
	active   *template.Template
	inactive *template.Template
	selected *template.Template
	details  *template.Template
	help     *template.Template
}

// Run runs the Select list. It returns the index of the selected element,
// and its value.
func (s *Select) Run() (int, string, error) {
	if s.Size == 0 {
		s.Size = 5
	}

	l, err := list.New(s.Items, s.Size)
	if err != nil {
		return 0, "", err
	}
	l.Searcher = s.Searcher

	s.list = l

	s.setKeys()

	err = s.prepareTemplates()
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

	rl, err := readline.NewEx(c)
	if err != nil {
		return 0, "", err
	}

	rl.Write([]byte(hideCursor))
	sb := screenbuf.New(rl)

	var searchInput []rune
	canSearch := s.Searcher != nil
	searchMode := s.StartInSearchMode

	c.SetListener(func(line []rune, pos int, key rune) ([]rune, int, bool) {
		switch {
		case key == KeyEnter:
			return nil, 0, true
		case key == s.Keys.Next.Code || (key == 'j' && !searchMode):
			s.list.Next()
		case key == s.Keys.Prev.Code || (key == 'k' && !searchMode):
			s.list.Prev()
		case key == s.Keys.Search.Code:
			if !canSearch {
				break
			}

			if searchMode {
				searchMode = false
				searchInput = nil
				s.list.CancelSearch()
			} else {
				searchMode = true
			}
		case key == KeyBackspace:
			if !canSearch || !searchMode {
				break
			}

			if len(searchInput) > 1 {
				searchInput = searchInput[:len(searchInput)-1]
				s.list.Search(string(searchInput))
			} else {
				searchInput = nil
				s.list.CancelSearch()
			}
		case key == s.Keys.PageUp.Code || (key == 'h' && !searchMode):
			s.list.PageUp()
		case key == s.Keys.PageDown.Code || (key == 'l' && !searchMode):
			s.list.PageDown()
		default:
			if canSearch && searchMode {
				searchInput = append(searchInput, line...)
				s.list.Search(string(searchInput))
			}
		}

		if searchMode {
			header := fmt.Sprintf("Search: %s%s", string(searchInput), cursor)
			sb.WriteString(header)
		} else {
			help := s.renderHelp(canSearch)
			sb.Write(help)
		}

		label := render(s.Templates.label, s.Label)
		sb.Write(label)

		items, idx := s.list.Items()
		last := len(items) - 1

		for i, item := range items {
			page := " "

			switch i {
			case 0:
				if s.list.CanPageUp() {
					page = "↑"
				} else {
					page = string(top)
				}
			case last:
				if s.list.CanPageDown() {
					page = "↓"
				}
			}

			output := []byte(page + " ")

			if i == idx {
				output = append(output, render(s.Templates.active, item)...)
			} else {
				output = append(output, render(s.Templates.inactive, item)...)
			}

			sb.Write(output)
		}

		if idx == list.NotFound {
			sb.WriteString("")
			sb.WriteString("No results")
		} else {
			active := items[idx]

			details := s.renderDetails(active)
			for _, d := range details {
				sb.Write(d)
			}
		}

		sb.Flush()

		return nil, 0, true
	})

	for {
		_, err = rl.Readline()

		if err != nil {
			switch {
			case err == readline.ErrInterrupt, err.Error() == "Interrupt":
				err = ErrInterrupt
			case err == io.EOF:
				err = ErrEOF
			}
			break
		}

		_, idx := s.list.Items()
		if idx != list.NotFound {
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
		return 0, "", err
	}

	items, idx := s.list.Items()
	item := items[idx]

	output := render(s.Templates.selected, item)

	sb.Reset()
	sb.Write(output)
	sb.Flush()
	rl.Write([]byte(showCursor))
	rl.Close()

	return s.list.Index(), fmt.Sprintf("%v", item), err
}

func (s *Select) prepareTemplates() error {
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

	if tpls.Help == "" {
		tpls.Help = fmt.Sprintf(`{{ "Use the arrow keys to navigate:" | faint }} {{ .NextKey | faint }} ` +
			`{{ .PrevKey | faint }} {{ .PageDownKey | faint }} {{ .PageUpKey | faint }} ` +
			`{{ if .Search }} {{ "and" | faint }} {{ .SearchKey | faint }} {{ "toggles search" | faint }}{{ end }}`)
	}

	tpl, err = template.New("").Funcs(tpls.FuncMap).Parse(tpls.Help)
	if err != nil {
		return err
	}

	tpls.help = tpl

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

		list, err := list.New(newItems, 5)
		if err != nil {
			return 0, "", err
		}

		s := Select{
			Label:     sa.Label,
			Items:     newItems,
			IsVimMode: sa.IsVimMode,
			Size:      5,
			list:      list,
		}
		s.setKeys()

		err = s.prepareTemplates()
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

func (s *Select) setKeys() {
	if s.Keys != nil {
		return
	}
	s.Keys = &SelectKeys{
		Prev:     Key{Code: KeyPrev, Display: "↑"},
		Next:     Key{Code: KeyNext, Display: "↓"},
		PageUp:   Key{Code: KeyBackward, Display: "←"},
		PageDown: Key{Code: KeyForward, Display: "→"},
		Search:   Key{Code: '/', Display: "/"},
	}
}

func (s *Select) renderDetails(item interface{}) [][]byte {
	if s.Templates.details == nil {
		return nil
	}

	var buf bytes.Buffer
	w := ansiterm.NewTabWriter(&buf, 0, 0, 8, ' ', 0)

	err := s.Templates.details.Execute(w, item)
	if err != nil {
		fmt.Fprintf(w, "%v", item)
	}

	w.Flush()

	output := buf.Bytes()

	return bytes.Split(output, []byte("\n"))
}

func (s *Select) renderHelp(b bool) []byte {
	keys := struct {
		NextKey     string
		PrevKey     string
		PageDownKey string
		PageUpKey   string
		Search      bool
		SearchKey   string
	}{
		NextKey:     s.Keys.Next.Display,
		PrevKey:     s.Keys.Prev.Display,
		PageDownKey: s.Keys.PageDown.Display,
		PageUpKey:   s.Keys.PageUp.Display,
		SearchKey:   s.Keys.Search.Display,
		Search:      b,
	}

	return render(s.Templates.help, keys)
}

func render(tpl *template.Template, data interface{}) []byte {
	var buf bytes.Buffer
	err := tpl.Execute(&buf, data)
	if err != nil {
		return []byte(fmt.Sprintf("%v", data))
	}
	return buf.Bytes()
}
