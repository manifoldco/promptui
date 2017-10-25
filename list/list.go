package list

import (
	"fmt"
	"reflect"
)

// List holds a collection of items that can be displayed with an N number of
// visible items. The list can be moved up, down by one item of time or an
// entire page (ie: visible size). It keeps track of the current selected item.
type List struct {
	items  []interface{}
	cursor int // cursor holds the index of the current selected item
	size   int // size is the number of visible options
	start  int
}

// New creates and initializes a list. Items must be a slice type and size must
// be greater than 0.
func New(items interface{}, size int) (*List, error) {
	if size < 1 {
		return nil, fmt.Errorf("list size %d must be greater than 0", size)
	}

	if items == nil || reflect.TypeOf(items).Kind() != reflect.Slice {
		return nil, fmt.Errorf("items %v is not a slice", items)
	}

	l := &List{size: size}

	slice := reflect.ValueOf(items)
	for i := 0; i < slice.Len(); i++ {
		item := slice.Index(i).Interface()
		l.items = append(l.items, item)
	}

	return l, nil
}

// Prev moves the visible list back one item. If the selected item is out of
// view, the new select item becomes the last visible item. If the list is
// already at the top, nothing happens.
func (l *List) Prev() {
	if l.cursor > 0 {
		l.cursor--
	}

	if l.start > l.cursor {
		l.start = l.cursor
	}
}

// Next moves the visible list forward one item. If the selected item is out of
// view, the new select item becomes the first visible item. If the list is
// already at the bottom, nothing happens.
func (l *List) Next() {
	max := len(l.items) - 1

	if l.cursor < max {
		l.cursor++
	}

	if l.start+l.size <= l.cursor {
		l.start = l.cursor - l.size + 1
	}
}

// PageUp moves the visible list backward by x items. Where x is the size of the
// visible items on the list. The selected item becomes the first visible item.
// If the list is already at the bottom, the selected item becomes the last
// visible item.
func (l *List) PageUp() {
	start := l.start - l.size
	if start < 0 {
		l.start = 0
	} else {
		l.start = start
	}

	cursor := l.start

	if cursor < l.cursor {
		l.cursor = cursor
	}
}

// PageDown moves the visible list forward by x items. Where x is the size of
// the visible items on the list. The selected item becomes the first visible
// item.
func (l *List) PageDown() {
	start := l.start + l.size
	max := len(l.items) - l.size

	if start > max {
		l.start = max
	} else {
		l.start = start
	}

	cursor := l.start

	if cursor == l.cursor {
		l.cursor = len(l.items) - 1

	} else if cursor > l.cursor {
		l.cursor = cursor
	}
}

// CanPageDown returns whether a list can still PageDown().
func (l *List) CanPageDown() bool {
	max := len(l.items)
	return l.start+l.size < max
}

// CanPageUp returns whether a list can still PageUp().
func (l *List) CanPageUp() bool {
	return l.start > 0
}

// Index returns the index of the item currently selected.
func (l *List) Index() int {
	return l.cursor
}

// Items returns a slice equal to the size of the list with the current visible
// items and the index of the active item in this list.
func (l *List) Items() ([]interface{}, int) {
	var result []interface{}
	max := len(l.items)
	end := l.start + l.size

	if end > max {
		end = max
	}

	active := 0

	for i, j := l.start, 0; i < end; i, j = i+1, j+1 {
		if l.cursor == i {
			active = j
		}

		result = append(result, l.items[i])
	}

	return result, active
}
