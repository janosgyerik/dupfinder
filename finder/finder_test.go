package finder

import (
	"testing"
	"os"
	"path/filepath"
)

const testDataDir = "testdata"

var basedir = testDataDir

func finderWithFilters(filters... Filter) Finder {
	return DefaultFinder{filters: filters}
}

func findPaths(finder Finder) []string {
	paths := []string{}
	for path := range finder.Find(basedir) {
		paths = append(paths, path)
	}
	return paths
}

func Test_should_find_all_files_and_only_files(t*testing.T) {
	finder := finderWithFilters()
	paths := findPaths(finder)

	expected := 9
	if len(paths) != expected {
		t.Fatalf("got %d files, expected %d", len(paths), expected)
	}

	assertPathsAreFiles(t, paths...)
}

func Test_should_find_size30_for_MinSize30(t*testing.T) {
	finder := finderWithFilters(MinSizeFilter{Size: 30})
	paths := findPaths(finder)

	if len(paths) != 3 {
		t.Fatalf("got %d files, expected 3", len(paths))
	}

	assertPathsAreFiles(t, paths...)
}

func assertPathsAreFiles(t*testing.T, paths... string) {
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			t.Error(err)
		} else if info.IsDir() {
			t.Errorf("got %s which is a directory, expected only files", info)
		}
	}
}

func Test_assertPathsAreFiles_should_pass_for_files(t*testing.T) {
	t2 := &testing.T{}
	assertPathsAreFiles(t2, filepath.Join(basedir, "size30.txt"))
	if t2.Failed() {
		t.Fatal("assertPathsAreFiles failed for a file, but it shouldn't have")
	}
}

func Test_assertPathsAreFiles_should_fail_for_dirs(t*testing.T) {
	t2 := &testing.T{}
	assertPathsAreFiles(t2, ".")
	if !t2.Failed() {
		t.Fatal("assertPathsAreFiles did not fail for a directory, but it should have")
	}
}

func Test_assertPathsAreFiles_should_fail_for_non_existent_paths(t*testing.T) {
	t2 := &testing.T{}
	assertPathsAreFiles(t2, "nonexistent")
	if !t2.Failed() {
		t.Fatal("assertPathsAreFiles did not fail for a nonexistent file, but it should have")
	}
}
