package dupfinder

import (
	"testing"
)

// TODO replace concrete files with fakes
// TODO should raise error if first file could not be opened
// TODO should raise error if second file could not be opened
// TODO should raise error if error happens while reading first file
// TODO should raise error if error happens while reading second file

func TestCompare_same_file(t*testing.T) {
	path := "/tmp/main.go"

	expected := 0

	cmp, err := CompareFiles(path, path)
	if err != nil {
		t.Errorf("Compare(f, f) raised error: %v", err)
	}
	if cmp != expected {
		t.Errorf("Compare(f, f) == %v, want %v", cmp, expected)
	}
}

func TestCompare_size_ascending(t*testing.T) {
	smaller := "/tmp/a"
	bigger := "/tmp/main.go"

	expected := -1

	cmp, err := CompareFiles(smaller, bigger)
	if err != nil {
		t.Errorf("Compare(smaller, bigger) raised error: %v", err)
	}
	if cmp != expected {
		t.Errorf("Compare(smaller, bigger) == %v, want %v", cmp, expected)
	}
}

func TestCompare_size_descending(t*testing.T) {
	smaller := "/tmp/a"
	bigger := "/tmp/main.go"

	expected := 1

	cmp, err := CompareFiles(bigger, smaller)
	if err != nil {
		t.Errorf("Compare(bigger, smaller) raised error: %v", err)
	}
	if cmp != expected {
		t.Errorf("Compare(bigger, smaller) == %v, want %v", cmp, expected)
	}
}

func TestCompare_same_size_content_ascending(t*testing.T) {
	lower := "/tmp/a"
	higher := "/tmp/b"

	expected := -1

	cmp, err := CompareFiles(lower, higher)
	if err != nil {
		t.Errorf("Compare(lower, higher) raised error: %v", err)
	}
	if cmp != expected {
		t.Errorf("Compare(lower, higher) == %v, want %v", cmp, expected)
	}
}

func TestCompare_same_size_content_descending(t*testing.T) {
	lower := "/tmp/a"
	higher := "/tmp/b"

	expected := 1

	cmp, err := CompareFiles(higher, lower)
	if err != nil {
		t.Errorf("Compare(higher, lower) raised error: %v", err)
	}
	if cmp != expected {
		t.Errorf("Compare(higher, lower) == %v, want %v", cmp, expected)
	}
}
