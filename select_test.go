package promptui

import "testing"

func TestForward(t *testing.T) {
	type example struct {
		start, end, selected, max int
	}

	tcs := []struct {
		scenario string
		input    example
		output   example
	}{
		{
			scenario: "when list is too short",
			input:    example{start: 0, end: 2, selected: 0, max: 3},
			output:   example{start: 0, end: 2, selected: 0},
		},
		{
			scenario: "when list is too long enough",
			input:    example{start: 0, end: 4, selected: 0, max: 10},
			output:   example{start: 4, end: 8, selected: 4},
		},
		{
			scenario: "when list is in the middle",
			input:    example{start: 2, end: 6, selected: 2, max: 10},
			output:   example{start: 5, end: 9, selected: 5},
		},
		{
			scenario: "when list is almost at the end",
			input:    example{start: 4, end: 8, selected: 7, max: 10},
			output:   example{start: 5, end: 9, selected: 7},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			start, end, selected := forward(tc.input.start, tc.input.end, tc.input.selected, tc.input.max)
			result := example{start: start, end: end, selected: selected}

			if tc.output != result {
				t.Errorf("expected %v to equal %v", tc.output, result)
			}
		})
	}
}

func TestBackward(t *testing.T) {
	type example struct {
		start, end, selected, max int
	}

	tcs := []struct {
		scenario string
		input    example
		output   example
	}{
		{
			scenario: "when list is too short",
			input:    example{start: 0, end: 2, selected: 0, max: 3},
			output:   example{start: 0, end: 2, selected: 0},
		},
		{
			scenario: "when list is in the beggining",
			input:    example{start: 2, end: 6, selected: 2, max: 10},
			output:   example{start: 0, end: 4, selected: 0},
		},
		{
			scenario: "when list is in the middle",
			input:    example{start: 3, end: 7, selected: 4, max: 10},
			output:   example{start: 0, end: 4, selected: 0},
		},
		{
			scenario: "when list is at the end",
			input:    example{start: 5, end: 9, selected: 7, max: 10},
			output:   example{start: 1, end: 5, selected: 1},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			start, end, selected := backward(tc.input.start, tc.input.end, tc.input.selected, tc.input.max)
			result := example{start: start, end: end, selected: selected}

			if tc.output != result {
				t.Errorf("expected %v to equal %v", tc.output, result)
			}
		})
	}
}
