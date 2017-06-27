package pathreader

import (
	"testing"
	"strings"
	"reflect"
)

func yes(string) bool {
	return true
}

func Test_readItems(t*testing.T) {
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
		actual := readItemsFromLines(strings.NewReader(item.input), yes)
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("got %#v, expected %#v", actual, expected)
		}

		input2 := strings.Replace(item.input, "\n", "\000", -1)
		actual2 := readItemsFromNullDelimited(strings.NewReader(input2), yes)
		if !reflect.DeepEqual(expected, actual2) {
			t.Errorf("got %#v, expected %#v", actual2, expected)
		}
	}
}

func Test_ReadPaths(t*testing.T) {
	data := []struct {
		input string
		lines []string
	}{
		{"", []string{}},
		{".\nnonexistent", []string{"."}},
		{".\n.", []string{"."}},
		{".\n.\n/", []string{".", "/"}},
		{".\n\n.\n.\n/\n/", []string{".", "/"}},
	}

	for _, item := range data {
		expected := item.lines
		actual := ReadPathsFromLines(strings.NewReader(item.input))
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("got %#v, expected %#v", actual, expected)
		}

		input2 := strings.Replace(item.input, "\n", "\000", -1)
		actual2 := ReadPathsFromNullDelimited(strings.NewReader(input2))
		if !reflect.DeepEqual(expected, actual2) {
			t.Errorf("got %#v, expected %#v", actual2, expected)
		}
	}
}

func Test_FilterPaths(t*testing.T) {
	data := []struct {
		input string
		paths []string
	}{
		{"", []string{}},
		{".\nnonexistent", []string{"."}},
		{".\n.", []string{"."}},
		{".\n.\n/", []string{".", "/"}},
		{".\n\n.\n.\n/\n/", []string{".", "/"}},
	}

	for _, item := range data {
		expected := item.paths
		actual := FilterPaths(strings.Split(item.input, "\n"))
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("got %#v, expected %#v", actual, expected)
		}
	}
}
