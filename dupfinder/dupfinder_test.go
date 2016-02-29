package dupfinder

import (
	"testing"
)

// TODO replace concrete files with fakes

func TestCompare_same_file(t*testing.T) {
	path := "/tmp/main.go"

	expected := 0

	cmp, err := Compare(path, path)
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

	cmp, err := Compare(smaller, bigger)
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

	cmp, err := Compare(bigger, smaller)
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

	cmp, err := Compare(lower, higher)
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

	cmp, err := Compare(higher, lower)
	if err != nil {
		t.Errorf("Compare(higher, lower) raised error: %v", err)
	}
	if cmp != expected {
		t.Errorf("Compare(higher, lower) == %v, want %v", cmp, expected)
	}
}
