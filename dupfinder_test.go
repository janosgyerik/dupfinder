package dupfinder

import (
	"testing"
	"strings"
	"io/ioutil"
	"os"
	"strconv"
)

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

func TestCompareFiles_equal_if_both_same_empty_dummy(t*testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "dummy")
	defer os.Remove(file.Name())

	expected := 0

	cmp, err := CompareFiles(file.Name(), file.Name())
	if err != nil {
		t.Errorf("Compare(dummy, dummy) raised error: %v", err)
	}
	if cmp != expected {
		t.Errorf("Compare(dummy, dummy) == %v, want %v", cmp, expected)
	}
}

func TestCompareFiles_equal_if_both_empty_dummy(t*testing.T) {
	dummy1, err := ioutil.TempFile(os.TempDir(), "dummy1")
	defer os.Remove(dummy1.Name())

	dummy2, err := ioutil.TempFile(os.TempDir(), "dummy2")
	defer os.Remove(dummy2.Name())

	expected := 0

	cmp, err := CompareFiles(dummy1.Name(), dummy2.Name())
	if err != nil {
		t.Errorf("Compare(dummy1, dummy2) raised error: %v", err)
	}
	if cmp != expected {
		t.Errorf("Compare(dummy1, dummy2) == %v, want %v", cmp, expected)
	}
}

func TestCompareFiles_empty_comes_before_nonempty(t*testing.T) {
	empty, err := ioutil.TempFile(os.TempDir(), "empty")
	defer os.Remove(empty.Name())

	nonempty, err := ioutil.TempFile(os.TempDir(), "nonempty")
	defer os.Remove(nonempty.Name())

	nonempty.WriteString("something")

	expected := -1

	cmp, err := CompareFiles(empty.Name(), nonempty.Name())
	if err != nil {
		t.Errorf("Compare(empty, nonempty) raised error: %v", err)
	}
	if cmp != expected {
		t.Errorf("Compare(empty, nonempty) == %v, want %v", cmp, expected)
	}
}

func TestCompareFiles_fails_if_both_nonexistent(t*testing.T) {
	_, err := CompareFiles("/nonexistent1", "/nonexistent2")
	if err == nil {
		t.Error("Compare(nonexistent1, nonexistent2) should have raised error")
	}
}

func TestCompareFiles_fails_if_first_nonexistent(t*testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "dummy")
	defer os.Remove(file.Name())

	_, err = CompareFiles("/nonexistent", file.Name())
	if err == nil {
		t.Error("Compare(nonexistent, dummy) should have raised error")
	}
}

func TestCompareFiles_fails_if_second_nonexistent(t*testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "dummy")
	defer os.Remove(file.Name())

	_, err = CompareFiles(file.Name(), "/nonexistent")
	if err == nil {
		t.Error("Compare(dummy, nonexistent) should have raised error")
	}
}

// TODO missing test: compare equal non-empty content

func createTempFileWithContent(content string) *os.File {
	file, _ := ioutil.TempFile(os.TempDir(), "prefix")
	file.WriteString(content)
	return file
}

func findDuplicates(distinctCount int, groupSizes... int) []Duplicates {
	paths := make([]string, 0)

	for i := 0; i < distinctCount; i++ {
		file := createTempFileWithContent("dummy distinct " + strconv.Itoa(i))
		defer os.Remove(file.Name())
		paths = append(paths, file.Name())
	}

	return FindDuplicates(paths...)
}

func TestFindDuplicates_two_duplicates(t*testing.T) {
	content := "dummy content"

	file1 := createTempFileWithContent(content)
	defer os.Remove(file1.Name())

	file2 := createTempFileWithContent(content)
	defer os.Remove(file2.Name())

	duplicates := FindDuplicates(file1.Name(), file2.Name())
	if len(duplicates) != 1 {
		t.Errorf("Found %d duplicate groups, expected %d", len(duplicates), 1)
	}
	if len(duplicates[0]) != 2 {
		t.Errorf("Found %d duplicate files, expected %d", len(duplicates[0]), 2)
	}
}

func TestFindDuplicates_three_duplicates(t*testing.T) {
	content := "dummy content"

	file1 := createTempFileWithContent(content)
	defer os.Remove(file1.Name())

	file2 := createTempFileWithContent(content)
	defer os.Remove(file2.Name())

	file3 := createTempFileWithContent(content)
	defer os.Remove(file3.Name())

	duplicates := FindDuplicates(file1.Name(), file2.Name(), file3.Name())
	if len(duplicates) != 1 {
		t.Errorf("Found %d duplicate groups, expected %d", len(duplicates), 1)
	}
	if len(duplicates[0]) != 3 {
		t.Errorf("Found %d duplicate files, expected %d", len(duplicates[0]), 3)
	}
}

func TestFindDuplicates_two_different(t*testing.T) {
	duplicates := findDuplicates(2)
	if len(duplicates) != 0 {
		t.Errorf("Found %d duplicate groups, expected %d", len(duplicates), 0)
	}
}

func TestFindDuplicates_three_different(t*testing.T) {
	duplicates := findDuplicates(3)
	if len(duplicates) != 0 {
		t.Errorf("Found %d duplicate groups, expected %d", len(duplicates), 0)
	}
}

func TestFindDuplicates_two_duplicate_groups(t*testing.T) {
	content1 := "dummy content 1"

	file1_1 := createTempFileWithContent(content1)
	defer os.Remove(file1_1.Name())

	file1_2 := createTempFileWithContent(content1)
	defer os.Remove(file1_2.Name())

	content2 := "dummy content 2"

	file2_1 := createTempFileWithContent(content2)
	defer os.Remove(file2_1.Name())

	file2_2 := createTempFileWithContent(content2)
	defer os.Remove(file2_2.Name())

	duplicates := FindDuplicates(file1_1.Name(), file1_2.Name(), file2_1.Name(), file2_2.Name())
	if len(duplicates) != 2 {
		t.Errorf("Found %d duplicate groups, expected %d", len(duplicates), 2)
	}
	if len(duplicates[0]) != 2 {
		t.Errorf("Found %d duplicate files in group 1, expected %d", len(duplicates[0]), 2)
	}
	if len(duplicates[1]) != 2 {
		t.Errorf("Found %d duplicate files in group 2, expected %d", len(duplicates[1]), 2)
	}
}
