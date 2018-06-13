package finder

import (
	"testing"
	"os"
	"io/ioutil"
	"reflect"
	"path"
)

var tempdir string

type fileData struct {
	relpath string
	size    int
}

func Test_Find_MinSize(t *testing.T) {
	fdata := []fileData{
		{"f1.txt", 1},
		{"f2.txt", 2},
		{"f3.txt", 3},
	}

	createTempFiles(fdata)
	defer deleteTempFiles()

	data := []struct {
		minSize  int64
		expected []string
	}{
		{minSize: 1, expected: []string{"f1.txt", "f2.txt", "f3.txt"}},
		{minSize: 2, expected: []string{"f2.txt", "f3.txt"}},
		{minSize: 3, expected: []string{"f3.txt"}},
		{minSize: 4, expected: nil},
	}

	for _, item := range data {
		finder := NewFinder(Filters.MinSize(item.minSize))
		actual := normalize(findPaths(finder))
		if !reflect.DeepEqual(item.expected, actual) {
			t.Errorf("got %#v; expected %#v", actual, item.expected)
		}
	}
}

func Test_Find_IncludeRegex(t *testing.T) {
	fdata := []fileData{
		{relpath: "a/f1.txt"},
		{relpath: "b/f1.pdf"},
		{relpath: "c/f2.pdf"},
	}

	createTempFiles(fdata)
	defer deleteTempFiles()

	data := []struct {
		pattern  string
		expected []string
	}{
		{pattern: `\.txt$`, expected: []string{"a/f1.txt"}},
		{pattern: `\.pdf$`, expected: []string{"b/f1.pdf", "c/f2.pdf"}},
		{pattern: `f2.pdf`, expected: []string{"c/f2.pdf"}},
		{pattern: `f1.txt`, expected: []string{"a/f1.txt"}},
		{pattern: `^f1`, expected: []string{"a/f1.txt", "b/f1.pdf"}},
	}

	for _, item := range data {
		finder := NewFinder(Filters.IncludeRegex(item.pattern))
		actual := normalize(findPaths(finder))
		if !reflect.DeepEqual(item.expected, actual) {
			t.Errorf("got %#v; expected %#v", actual, item.expected)
		}
	}
}

func Test_Find_ExcludeRegex(t *testing.T) {
	fdata := []fileData{
		{relpath: "a/f1.txt"},
		{relpath: "b/f1.pdf"},
		{relpath: "c/f2.pdf"},
	}

	createTempFiles(fdata)
	defer deleteTempFiles()

	data := []struct {
		pattern  string
		expected []string
	}{
		{pattern: `\.txt$`, expected: []string{"b/f1.pdf", "c/f2.pdf"}},
		{pattern: `\.pdf$`, expected: []string{"a/f1.txt"}},
		{pattern: `f2.pdf`, expected: []string{"a/f1.txt", "b/f1.pdf"}},
		{pattern: `f1.txt`, expected: []string{"b/f1.pdf", "c/f2.pdf"}},
		{pattern: `^f1`, expected: []string{"c/f2.pdf"}},
	}

	for _, item := range data {
		finder := NewFinder(Filters.ExcludeRegex(item.pattern))
		actual := normalize(findPaths(finder))
		if !reflect.DeepEqual(item.expected, actual) {
			t.Errorf("got %#v; expected %#v", actual, item.expected)
		}
	}
}

func normalize(paths []string) []string {
	var stripped []string
	for _, p := range paths {
		stripped = append(stripped, p[len(tempdir)+1:])
	}
	return stripped
}

func createTempFiles(data []fileData) {
	var err error
	tempdir, err = ioutil.TempDir("", "test")
	check(err)

	for _, v := range data {
		p := path.Join(tempdir, v.relpath)
		basedir := path.Dir(p)
		os.MkdirAll(basedir, 0755)
		err := ioutil.WriteFile(p, make([]byte, v.size), 0644)
		check(err)
	}
}

func deleteTempFiles() {
	err := os.RemoveAll(tempdir)
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func findPaths(finder Finder) []string {
	var paths []string
	for p := range finder.Find(tempdir) {
		paths = append(paths, p)
	}
	return paths
}
