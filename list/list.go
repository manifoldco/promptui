package list

import (
	"fmt"
	"reflect"
)

type List struct {
	items  []interface{}
	cursor int
	size   int
	start  int
}

func New(items interface{}, size int) (*List, error) {
	if items == nil || reflect.TypeOf(items).Kind() != reflect.Slice {
		return nil, fmt.Errorf("items %v is not a slice", items)
	}

	l := &List{size: size}

	slice := reflect.ValueOf(items)
	for i := 0; i < slice.Len(); i++ {
		l.items = append(l.items, slice.Index(i).Interface())
	}

	return l, nil
}

func (l *List) Prev() {
	if l.cursor > 0 {
		l.cursor--
	}

	if l.start > l.cursor {
		l.start = l.cursor
	}
}

func (l *List) Next() {
	max := len(l.items) - 1

	if l.cursor < max {
		l.cursor++
	}

	if l.start+l.size <= l.cursor {
		l.start = l.cursor - l.size + 1
	}
}

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

func (l *List) Selected() interface{} {
	return l.items[l.cursor]
}

func (l *List) Display() []interface{} {
	var result []interface{}
	end := l.start + l.size

	for i := l.start; i < end; i++ {
		result = append(result, l.items[i])
	}

	return result
}
