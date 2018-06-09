package pathreader

import (
	"testing"
	"strings"
	"reflect"
		)

func Test_readItems(t *testing.T) {
	data := []struct {
		input string
		lines []string
	}{
		{"", []string{}},
		{"line 1", []string{"line 1"}},
		{"line 1\nline 2", []string{"line 1", "line 2"}},
		{"line 1\nline 2\n", []string{"line 1", "line 2"}},
		{"line 1\n\nline 2\n", []string{"line 1", "", "line 2"}},
	}

	for _, item := range data {
		expected := item.lines
		actual := toString(readItemsFromLines(strings.NewReader(item.input)))
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("got %#v, expected %#v", actual, expected)
		}

		input2 := strings.Replace(item.input, "\n", "\000", -1)
		actual2 := toString(readItemsFromNullDelimited(strings.NewReader(input2)))
		if !reflect.DeepEqual(expected, actual2) {
			t.Errorf("got %#v, expected %#v", actual2, expected)
		}
	}
}

func toString(c <-chan string) interface{} {
	lines := make([]string, 0)
	for line := range c {
		lines = append(lines, line)
	}
	return lines
}
