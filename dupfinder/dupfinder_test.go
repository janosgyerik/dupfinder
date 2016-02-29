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

func TestCompare_shorter_size_to_longer_size(t*testing.T) {
	shorter := "/tmp/a"
	longer := "/tmp/main.go"

	expected := -1

	cmp, err := Compare(shorter, longer)
	if err != nil {
		t.Errorf("Compare(shorter, longer) raised error: %v", err)
	}
	if cmp != expected {
		t.Errorf("Compare(shorter, longer) == %v, want %v", cmp, expected)
	}
}

func TestCompare_longer_size_to_shorter_size(t*testing.T) {
	shorter := "/tmp/a"
	longer := "/tmp/main.go"

	expected := 1

	cmp, err := Compare(longer, shorter)
	if err != nil {
		t.Errorf("Compare(longer, shorter) raised error: %v", err)
	}
	if cmp != expected {
		t.Errorf("Compare(longer, shorter) == %v, want %v", cmp, expected)
	}
}
