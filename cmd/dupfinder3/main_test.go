package main

import (
	"testing"
	"io/ioutil"
	"path"
	"os"
	"reflect"
	"os/exec"
	"strings"
)

var tempdir string

type fileData struct {
	relpath string
	content string
}

func Test_find_no_groups_from_two_distinct(t *testing.T) {
	fdata := []fileData{
		{"f1.txt", "foo"},
		{"f2.txt", "bar"},
	}
	var expected [][]string

	createTempFiles(fdata)
	defer deleteTempFiles()

	if actual := normalize(run()); !reflect.DeepEqual(expected, actual) {
		t.Fatalf("got:\n%#v\nexpected:\n%#v", actual, expected)
	}
}

func Test_find_one_group(t *testing.T) {
	fdata := []fileData{
		{"f1.txt", "foo"},
		{"f2.txt", "foo"},
	}
	expected := [][]string{{"f1.txt", "f2.txt"}}

	createTempFiles(fdata)
	defer deleteTempFiles()

	if actual := normalize(run()); !reflect.DeepEqual(expected, actual) {
		t.Fatalf("got:\n%#v\nexpected:\n%#v", actual, expected)
	}
}

func Test_find_multiple_varied_sized_groups(t *testing.T) {
	fdata := []fileData{
		{"f1.txt", "foo"},
		{"a/f2.txt", "foo"},
		{"a/b/f3.txt", "foo"},
		{"b/f1.txt", "bar"},
		{"b/c/f2.txt", "bar"},
		{"c/f1.txt", "baz"},
	}
	expected := [][]string{
		{
			"a/b/f3.txt",
			"a/f2.txt",
			"f1.txt",
		},
		{
			"b/c/f2.txt",
			"b/f1.txt",
		},
	}

	createTempFiles(fdata)
	defer deleteTempFiles()

	if actual := normalize(run()); !reflect.DeepEqual(expected, actual) {
		t.Fatalf("got:\n%#v\nexpected:\n%#v", actual, expected)
	}
}

func normalize(out string) [][]string {
	var result [][]string
	var current []string

	for _, p := range strings.Split(out, "\n") {
		if len(p) == 0 {
			if len(current) > 0 {
				result = append(result, current)
			}
			current = make([]string, 0)
		} else {
			current = append(current, p[len(tempdir)+1:])
		}
	}
	return result
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createTempFiles(data []fileData) {
	var err error
	tempdir, err = ioutil.TempDir("", "test")
	check(err)

	for _, v := range data {
		p := path.Join(tempdir, v.relpath)
		basedir := path.Dir(p)
		os.MkdirAll(basedir, 0755)
		err := ioutil.WriteFile(p, []byte(v.content), 0644)
		check(err)
	}
}

func deleteTempFiles() {
	err := os.RemoveAll(tempdir)
	check(err)
}

func run() string {
	out, err := exec.Command("go", "run", "main.go", tempdir).Output()
	check(err)
	return string(out)
}
