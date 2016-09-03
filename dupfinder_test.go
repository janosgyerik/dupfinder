package dupfinder

import (
	"testing"
	"strings"
)

// TODO should raise error if first file could not be opened
// TODO should raise error if second file could not be opened
// TODO should raise error if error happens while reading first file
// TODO should raise error if error happens while reading second file

func TestCompareReaders_same_file(t*testing.T) {
	content := "dummy content"
	reader1 := strings.NewReader(content)
	reader2 := strings.NewReader(content)

	expected := 0

	cmp, err := CompareReaders(reader1, reader2)
	if err != nil {
		t.Errorf("Compare(f, f) raised error: %v", err)
	}
	if cmp != expected {
		t.Errorf("Compare(f, f) == %v, want %v", cmp, expected)
	}
}

func TestCompareReaders_size_ascending(t*testing.T) {
	smaller := strings.NewReader("dummy content")
	bigger := strings.NewReader("longer dummy content")

	expected := -1

	cmp, err := CompareReaders(smaller, bigger)
	if err != nil {
		t.Errorf("Compare(smaller, bigger) raised error: %v", err)
	}
	if cmp != expected {
		t.Errorf("Compare(smaller, bigger) == %v, want %v", cmp, expected)
	}
}

func TestCompareReaders_size_descending(t*testing.T) {
	smaller := strings.NewReader("dummy content")
	bigger := strings.NewReader("longer dummy content")

	expected := 1

	cmp, err := CompareReaders(bigger, smaller)
	if err != nil {
		t.Errorf("Compare(bigger, smaller) raised error: %v", err)
	}
	if cmp != expected {
		t.Errorf("Compare(bigger, smaller) == %v, want %v", cmp, expected)
	}
}

func TestCompareReaders_same_size_content_ascending(t*testing.T) {
	lower := strings.NewReader("dummy content a")
	higher := strings.NewReader("dummy content b")

	expected := -1

	cmp, err := CompareReaders(lower, higher)
	if err != nil {
		t.Errorf("Compare(lower, higher) raised error: %v", err)
	}
	if cmp != expected {
		t.Errorf("Compare(lower, higher) == %v, want %v", cmp, expected)
	}
}

func TestCompareReaders_same_size_content_descending(t*testing.T) {
	lower := strings.NewReader("dummy content a")
	higher := strings.NewReader("dummy content b")

	expected := 1

	cmp, err := CompareReaders(higher, lower)
	if err != nil {
		t.Errorf("Compare(higher, lower) raised error: %v", err)
	}
	if cmp != expected {
		t.Errorf("Compare(higher, lower) == %v, want %v", cmp, expected)
	}
}

func TestCompareFiles_fails_if_both_nonexistent(t*testing.T) {
	_, err := CompareFiles("/nonexistent1", "/nonexistent2")
	if err == nil {
		t.Error("Compare(nonexistent1, nonexistent2) should have raised error")
	}
}
