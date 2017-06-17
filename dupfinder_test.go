package dupfinder

import (
	"testing"
	"strings"
	"io/ioutil"
	"os"
	"strconv"
	"github.com/janosgyerik/dupfinder/finder"
	"path/filepath"
	"errors"
)

const testdataBasedir = "testdata"

func Test_cmpResult(t*testing.T) {
	err := errors.New("dummy")

	cmpResultTests := []struct {
		label            string
		actual, expected cmpResult
	}{
		{
			"done 1",
			done(1),
			cmpResult{cmp: 1, success: true, done: true},
		},
		{
			"done -1",
			done(-1),
			cmpResult{cmp: -1, success: true, done: true},
		},
		{
			"done 0",
			done(0),
			cmpResult{cmp: 0, success: true, done: true},
		},
		{
			"undecided",
			undecided(),
			cmpResult{},
		},
		{
			"errFirst",
			errFirst(err),
			cmpResult{errFirst: err, done: true},
		},
		{
			"errSecond",
			errSecond(err),
			cmpResult{errSecond: err, done: true},
		},
	}
	for _, tt := range cmpResultTests {
		if tt.expected != tt.actual {
			t.Errorf("got %#v for %s, expected %#v", tt.actual, tt.label, tt.expected)
		}
	}
}

func Test_compareReaders_same_file(t*testing.T) {
	content := "dummy content"
	reader1 := strings.NewReader(content)
	reader2 := strings.NewReader(content)

	expected := 0

	r := compareReaders(reader1, reader2)
	if !r.success {
		t.Errorf("Compare(f, f) failed: %v", r)
	}
	if r.cmp != expected {
		t.Errorf("Compare(f, f) == %v, want %v", r.cmp, expected)
	}
}

func Test_compareReaders_size_ascending(t*testing.T) {
	smaller := strings.NewReader("dummy content")
	bigger := strings.NewReader("longer dummy content")

	expected := -1

	r := compareReaders(smaller, bigger)
	if !r.success {
		t.Errorf("Compare(smaller, bigger) failed: %v", r)
	}
	if r.cmp != expected {
		t.Errorf("Compare(smaller, bigger) == %v, want %v", r.cmp, expected)
	}
}

func Test_compareReaders_size_descending(t*testing.T) {
	smaller := strings.NewReader("dummy content")
	bigger := strings.NewReader("longer dummy content")

	expected := 1

	r := compareReaders(bigger, smaller)
	if !r.success {
		t.Errorf("Compare(bigger, smaller) failed: %v", r)
	}
	if r.cmp != expected {
		t.Errorf("Compare(bigger, smaller) == %v, want %v", r.cmp, expected)
	}
}

func Test_compareReaders_same_size_content_ascending(t*testing.T) {
	lower := strings.NewReader("dummy content a")
	higher := strings.NewReader("dummy content b")

	expected := -1

	r := compareReaders(lower, higher)
	if !r.success {
		t.Errorf("Compare(lower, higher) failed: %v", r)
	}
	if r.cmp != expected {
		t.Errorf("Compare(lower, higher) == %v, want %v", r.cmp, expected)
	}
}

func Test_compareReaders_same_size_content_descending(t*testing.T) {
	lower := strings.NewReader("dummy content a")
	higher := strings.NewReader("dummy content b")

	expected := 1

	r := compareReaders(higher, lower)
	if !r.success {
		t.Errorf("Compare(higher, lower) failed: %v", r)
	}
	if r.cmp != expected {
		t.Errorf("Compare(higher, lower) == %v, want %v", r.cmp, expected)
	}
}

func newTempFile(t*testing.T, name string) *os.File {
	file, err := ioutil.TempFile(os.TempDir(), name)

	if err != nil {
		t.Error("Failed to create temporary file")
	}

	return file
}

func newTempFileWithContent(t*testing.T, content string) *os.File {
	file := newTempFile(t, "dummy")
	file.WriteString(content)
	return file
}

func Test_compareFiles_equal_if_both_same_empty_dummy(t*testing.T) {
	file := newTempFile(t, "dummy")
	defer os.Remove(file.Name())

	expected := 0

	r := compareFiles(file.Name(), file.Name())
	if !r.success {
		t.Errorf("Compare(dummy, dummy) failed: %v", r)
	}
	if r.cmp != expected {
		t.Errorf("Compare(dummy, dummy) == %v, want %v", r.cmp, expected)
	}
}

func Test_compareFiles_equal_if_both_empty_dummy(t*testing.T) {
	dummy1 := newTempFile(t, "dummy1")
	defer os.Remove(dummy1.Name())

	dummy2 := newTempFile(t, "dummy2")
	defer os.Remove(dummy2.Name())

	expected := 0

	r := compareFiles(dummy1.Name(), dummy2.Name())
	if !r.success {
		t.Errorf("Compare(dummy1, dummy2) failed: %v", r)
	}
	if r.cmp != expected {
		t.Errorf("Compare(dummy1, dummy2) == %v, want %v", r.cmp, expected)
	}
}

func Test_compareFiles_empty_comes_before_nonempty(t*testing.T) {
	empty := newTempFile(t, "empty")
	defer os.Remove(empty.Name())

	nonempty := newTempFile(t, "nonempty")
	defer os.Remove(nonempty.Name())

	nonempty.WriteString("something")

	expected := -1

	r := compareFiles(empty.Name(), nonempty.Name())
	if !r.success {
		t.Errorf("Compare(empty, nonempty) failed: %v", r)
	}
	if r.cmp != expected {
		t.Errorf("Compare(empty, nonempty) == %v, want %v", r.cmp, expected)
	}
}

func Test_compareFiles_fails_if_both_nonexistent(t*testing.T) {
	r := compareFiles("/nonexistent1", "/nonexistent2")
	if r.errFirst == nil {
		t.Error("Compare(nonexistent1, nonexistent2) should have raised error")
	}
}

func Test_compareFiles_fails_if_first_nonexistent(t*testing.T) {
	file := newTempFile(t, "dummy")
	defer os.Remove(file.Name())

	r := compareFiles("/nonexistent", file.Name())
	if r.success || r.errFirst == nil || r.errSecond != nil || !r.done {
		t.Error("Compare(nonexistent, dummy) should have raised error")
	}
}

func Test_compareFiles_fails_if_second_nonexistent(t*testing.T) {
	file := newTempFile(t, "dummy")
	defer os.Remove(file.Name())

	r := compareFiles(file.Name(), "/nonexistent")
	if r.success || r.errFirst != nil || r.errSecond == nil || !r.done {
		t.Error("Compare(dummy, nonexistent) should have raised error")
	}
}

func findDuplicates(t*testing.T, distinctCount int, groupSizes... int) []DupGroup {
	paths := make([]string, 0)

	for i := 0; i < distinctCount; i++ {
		file := newTempFileWithContent(t, "dummy distinct " + strconv.Itoa(i))
		defer os.Remove(file.Name())
		paths = append(paths, file.Name())
	}

	for group, groupSize := range groupSizes {
		content := "dummy group " + strconv.Itoa(group)
		for i := 0; i < groupSize; i++ {
			file := newTempFileWithContent(t, content)
			defer os.Remove(file.Name())
			paths = append(paths, file.Name())
		}
	}

	return FindDuplicates(paths).Groups
}

func Test_FindDuplicates_two_duplicates(t*testing.T) {
	duplicates := findDuplicates(t, 0, 2)
	if len(duplicates) != 1 {
		t.Errorf("Found %d duplicate groups, expected %d", len(duplicates), 1)
	}
	if duplicates[0].count() != 2 {
		t.Errorf("Found %d duplicate files, expected %d", duplicates[0].count(), 2)
	}
}

func Test_FindDuplicates_three_duplicates(t*testing.T) {
	duplicates := findDuplicates(t, 0, 3)
	if len(duplicates) != 1 {
		t.Errorf("Found %d duplicate groups, expected %d", len(duplicates), 1)
	}
	if duplicates[0].count() != 3 {
		t.Errorf("Found %d duplicate files, expected %d", duplicates[0].count(), 3)
	}
}

func Test_FindDuplicates_two_different(t*testing.T) {
	duplicates := findDuplicates(t, 2)
	if len(duplicates) != 0 {
		t.Errorf("Found %d duplicate groups, expected %d", len(duplicates), 0)
	}
}

func Test_FindDuplicates_three_different(t*testing.T) {
	duplicates := findDuplicates(t, 3)
	if len(duplicates) != 0 {
		t.Errorf("Found %d duplicate groups, expected %d", len(duplicates), 0)
	}
}

func Test_FindDuplicates_two_duplicate_groups(t*testing.T) {
	duplicates := findDuplicates(t, 0, 2, 3)
	if len(duplicates) != 2 {
		t.Errorf("Found %d duplicate groups, expected %d", len(duplicates), 2)
	}
	if duplicates[0].count() != 2 {
		t.Errorf("Found %d duplicate files in group 1, expected %d", duplicates[0].count(), 2)
	}
	if duplicates[1].count() != 3 {
		t.Errorf("Found %d duplicate files in group 2, expected %d", duplicates[1].count(), 3)
	}
}

func verifyFailure(t*testing.T, paths []string, failurePath string) {
	result := FindDuplicates(paths)

	if len(result.Groups) != 0 {
		t.Errorf("Got %d duplicate groups, expected none", len(result.Groups))
	}

	if len(result.Failures) != 1 {
		t.Errorf("Got %d failures, expected 1", len(result.Failures))
	}

	if actual := result.Failures[0].Path; actual != failurePath {
		t.Errorf("Got %s, expected %s as failed path", actual, failurePath)
	}
}

func Test_FindDuplicates_nonexistent_files(t*testing.T) {
	file := newTempFile(t, "dummy")
	defer os.Remove(file.Name())

	failurePath := "/nonexistent"

	verifyFailure(t, []string{failurePath, file.Name()}, failurePath)
	verifyFailure(t, []string{file.Name(), failurePath}, failurePath)
}

func Test_FindDuplicates_unreadable_files(t*testing.T) {
	file := newTempFile(t, "dummy")
	defer os.Remove(file.Name())

	unreadable := newTempFile(t, "unreadable")
	defer os.Remove(unreadable.Name())

	os.Chmod(unreadable.Name(), 0)

	failurePath := unreadable.Name()

	verifyFailure(t, []string{failurePath, file.Name()}, failurePath)
	verifyFailure(t, []string{file.Name(), failurePath}, failurePath)
}

func Test_dupTracker_add_and_merge(t*testing.T) {
	tracker := newDupTracker()
	tracker.add("path1-1", "path1-2")
	tracker.add("path1-2", "path1-3")
	tracker.add("path2-1", "path2-2")
	tracker.add("path2-3", "path2-2")
	tracker.add("path1-1", "path2-2")

	duplicates := tracker.getDupGroups()
	if len(duplicates) != 1 {
		t.Errorf("Found %d duplicate groups, expected %d", len(duplicates), 1)
	}
	if duplicates[0].count() != 6 {
		t.Errorf("Found %d duplicate files, expected %d", duplicates[0].count(), 6)
	}
}

func findPaths(basename string) []string {
	return finder.NewFinder().FindAll(filepath.Join(testdataBasedir, basename))
}

func Test_nodups(t*testing.T) {
	duplicates := FindDuplicates(findPaths("nodups")).Groups

	if len(duplicates) > 0 {
		t.Fatal("found duplicates in different files:", duplicates)
	}
}

func Test_samesize(t*testing.T) {
	duplicates := FindDuplicates(findPaths("samesize")).Groups

	if len(duplicates) > 0 {
		t.Fatal("found duplicates in different files with same size:", duplicates)
	}
}

func Test_alldups(t*testing.T) {
	duplicates := FindDuplicates(findPaths("alldups")).Groups

	if len(duplicates) != 1 {
		t.Fatalf("got %d duplicate groups, expected 1", len(duplicates))
	}
}
