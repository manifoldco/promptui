package list

import (
	"fmt"
	"reflect"
	"testing"
)

func TestListMovement(t *testing.T) {
	letters := []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j'}

	l, err := New(letters, 4)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	tcs := []struct {
		expect   []rune
		move     string
		selected rune
	}{
		{move: "next", selected: 'b', expect: []rune{'a', 'b', 'c', 'd'}},
		{move: "prev", selected: 'a', expect: []rune{'a', 'b', 'c', 'd'}},
		{move: "prev", selected: 'a', expect: []rune{'a', 'b', 'c', 'd'}},
		{move: "next", selected: 'b', expect: []rune{'a', 'b', 'c', 'd'}},
		{move: "next", selected: 'c', expect: []rune{'a', 'b', 'c', 'd'}},
		{move: "next", selected: 'd', expect: []rune{'a', 'b', 'c', 'd'}},
		{move: "next", selected: 'e', expect: []rune{'b', 'c', 'd', 'e'}},
		{move: "prev", selected: 'd', expect: []rune{'b', 'c', 'd', 'e'}},
		{move: "up", selected: 'a', expect: []rune{'a', 'b', 'c', 'd'}},
		{move: "up", selected: 'a', expect: []rune{'a', 'b', 'c', 'd'}},
		{move: "down", selected: 'e', expect: []rune{'e', 'f', 'g', 'h'}},
		{move: "down", selected: 'g', expect: []rune{'g', 'h', 'i', 'j'}},
		{move: "down", selected: 'j', expect: []rune{'g', 'h', 'i', 'j'}},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("list %s", tc.move), func(t *testing.T) {
			switch tc.move {
			case "next":
				l.Next()
			case "prev":
				l.Prev()
			case "up":
				l.PageUp()
			case "down":
				l.PageDown()
			default:
				t.Fatalf("unknown move %q", tc.move)
			}

			got := castList(t, l.Display())

			if !reflect.DeepEqual(tc.expect, got) {
				t.Errorf("expected %q, got %q", tc.expect, got)
			}

			selected := l.Selected()

			if tc.selected != selected {
				t.Errorf("expected selected to be %q, got %q", tc.selected, selected)
			}
		})
	}
}

func castList(t *testing.T, list []interface{}) []rune {
	t.Helper()

	result := make([]rune, len(list))
	for i, l := range list {
		result[i] = l.(rune)
	}
	return result
}
