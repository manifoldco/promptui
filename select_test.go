package promptui

import (
	"testing"
)

type example struct {
	start, end, selected, max, size int
}

func TestSelectTemplateRender(t *testing.T) {
	t.Run("when using default style", func(t *testing.T) {
		values := []string{"Zero"}
		s := Select{
			Label: "Select Number",
			Items: values,
		}
		err := s.prepareTemplates()
		if err != nil {
			t.Fatalf("Unexpected error preparing templates %v", err)
		}

		result := render(s.Templates.label, s.Label)
		exp := "\x1b[34m?\x1b[0m Select Number: "
		if result != exp {
			t.Errorf("Expected label to eq %q, got %q", exp, result)
		}

		result = render(s.Templates.active, s.items[0])
		exp = "\x1b[1m▸\x1b[0m \x1b[4mZero\x1b[0m"
		if result != exp {
			t.Errorf("Expected active item to eq %q, got %q", exp, result)
		}

		result = render(s.Templates.inactive, s.items[0])
		exp = "  Zero"
		if result != exp {
			t.Errorf("Expected inactive item to eq %q, got %q", exp, result)
		}

		result = render(s.Templates.selected, s.items[0])
		exp = "\x1b[32m\x1b[32m✔\x1b[0m \x1b[2mZero\x1b[0m"
		if result != exp {
			t.Errorf("Expected selected item to eq %q, got %q", exp, result)
		}
	})

	t.Run("when using custom style", func(t *testing.T) {
		type pepper struct {
			Name        string
			HeatUnit    int
			Peppers     int
			Description string
		}
		peppers := []pepper{
			{
				Name:        "Bell Pepper",
				HeatUnit:    0,
				Peppers:     1,
				Description: "Not very spicy!",
			},
		}

		templates := &SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "\U0001F525 {{ .Name | bold }} ({{ .HeatUnit | red | italic }})",
			Inactive: "   {{ .Name | bold }} ({{ .HeatUnit | red | italic }})",
			Selected: "\U0001F525 {{ .Name | red | bold }}",
			Details: `Name: {{.Name}}
Peppers: {{.Peppers}}
Description: {{.Description}}`,
		}

		s := Select{
			Label:     "Spicy Level",
			Items:     peppers,
			Templates: templates,
		}

		err := s.prepareTemplates()
		if err != nil {
			t.Fatalf("Unexpected error preparing templates %v", err)
		}

		result := render(s.Templates.label, s.Label)
		exp := "Spicy Level?"
		if result != exp {
			t.Errorf("Expected label to eq %q, got %q", exp, result)
		}

		result = render(s.Templates.active, s.items[0])
		exp = "🔥 \x1b[1mBell Pepper\x1b[0m (\x1b[3m\x1b[31m0\x1b[0m)"
		if result != exp {
			t.Errorf("Expected active item to eq %q, got %q", exp, result)
		}

		result = render(s.Templates.inactive, s.items[0])
		exp = "   \x1b[1mBell Pepper\x1b[0m (\x1b[3m\x1b[31m0\x1b[0m)"
		if result != exp {
			t.Errorf("Expected inactive item to eq %q, got %q", exp, result)
		}

		result = render(s.Templates.selected, s.items[0])
		exp = "🔥 \x1b[1m\x1b[31mBell Pepper\x1b[0m"
		if result != exp {
			t.Errorf("Expected selected item to eq %q, got %q", exp, result)
		}

		result = render(s.Templates.details, s.items[0])
		exp = "Name: Bell Pepper\nPeppers: 1\nDescription: Not very spicy!"
		if result != exp {
			t.Errorf("Expected selected item to eq %q, got %q", exp, result)
		}
	})

	t.Run("when a template is invalid", func(t *testing.T) {
		templates := &SelectTemplates{
			Label: "{{ . ",
		}

		s := Select{
			Label:     "Spicy Level",
			Templates: templates,
		}

		err := s.prepareTemplates()
		if err == nil {
			t.Fatalf("Expected error got none")
		}
	})

	t.Run("when items is nil", func(t *testing.T) {
		s := Select{}

		err := s.prepareTemplates()
		if err == nil {
			t.Fatalf("Expected error got none")
		}
	})

	t.Run("when items is not a slice", func(t *testing.T) {
		s := Select{
			Items: "1,2,3",
		}

		err := s.prepareTemplates()
		if err == nil {
			t.Fatalf("Expected error got none")
		}
	})

	t.Run("when a template render fails", func(t *testing.T) {
		templates := &SelectTemplates{
			Label: "{{ .InvalidName }}",
		}

		s := Select{
			Label:     struct{ Name string }{Name: "Pepper"},
			Items:     []string{},
			Templates: templates,
		}

		err := s.prepareTemplates()
		if err != nil {
			t.Fatalf("Unexpected error preparing templates %v", err)
		}

		result := render(s.Templates.label, s.Label)
		exp := "{Pepper}"
		if result != exp {
			t.Errorf("Expected label to eq %q, got %q", exp, result)
		}
	})
}

func TestPageDown(t *testing.T) {
	tcs := []struct {
		scenario string
		input    example
		output   example
	}{
		{
			scenario: "when list is too short",
			input:    example{start: 0, end: 2, selected: 0, max: 2, size: 5},
			output:   example{start: 0, end: 2, selected: 0},
		},
		{
			scenario: "when list is too long enough",
			input:    example{start: 0, end: 4, selected: 0, max: 9, size: 5},
			output:   example{start: 4, end: 8, selected: 4},
		},
		{
			scenario: "when list is in the middle",
			input:    example{start: 2, end: 6, selected: 2, max: 9, size: 5},
			output:   example{start: 5, end: 9, selected: 5},
		},
		{
			scenario: "when list is almost at the end",
			input:    example{start: 4, end: 8, selected: 7, max: 9, size: 5},
			output:   example{start: 5, end: 9, selected: 7},
		},
		{
			scenario: "when list has a large size",
			input:    example{start: 4, end: 8, selected: 7, max: 9, size: 10},
			output:   example{start: 0, end: 9, selected: 7},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			s := &Select{Size: tc.input.size}
			start, end, selected := s.pagedown(tc.input.start, tc.input.end, tc.input.selected, tc.input.max)
			result := example{start: start, end: end, selected: selected}

			if tc.output != result {
				t.Errorf("expected %v to equal %v", tc.output, result)
			}
		})
	}
}

func TestPageUp(t *testing.T) {
	tcs := []struct {
		scenario string
		input    example
		output   example
	}{
		{
			scenario: "when list is too short",
			input:    example{start: 0, end: 2, selected: 0, max: 2, size: 5},
			output:   example{start: 0, end: 2, selected: 0},
		},
		{
			scenario: "when list is in the beginning",
			input:    example{start: 2, end: 6, selected: 2, max: 9, size: 5},
			output:   example{start: 0, end: 4, selected: 0},
		},
		{
			scenario: "when list is in the middle",
			input:    example{start: 3, end: 7, selected: 4, max: 9, size: 5},
			output:   example{start: 0, end: 4, selected: 0},
		},
		{
			scenario: "when list is at the end",
			input:    example{start: 5, end: 9, selected: 7, max: 9, size: 5},
			output:   example{start: 1, end: 5, selected: 1},
		},
		{
			scenario: "when list has a large size",
			input:    example{start: 5, end: 9, selected: 7, max: 9, size: 10},
			output:   example{start: 0, end: 9, selected: 0},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			s := &Select{Size: tc.input.size}
			start, end, selected := s.pageup(tc.input.start, tc.input.end, tc.input.selected, tc.input.max)
			result := example{start: start, end: end, selected: selected}

			if tc.output != result {
				t.Errorf("expected %v to equal %v", tc.output, result)
			}
		})
	}
}
