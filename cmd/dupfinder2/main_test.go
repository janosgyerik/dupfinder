package main

import (
	"io/ioutil"
	"testing"
	"os"
	"os/exec"
	"path"
	)

var tempdir string

type fileData struct {
	name    string
	content string
}

func Test(t *testing.T) {
	fdata := []fileData{
		{"f1.txt", "foo"},
		{"a/f2.txt", "foo"},
		{"a/b/f3.txt", "foo"},
		{"b/f1.txt", "bar"},
		{"b/c/f2.txt", "bar"},
		{"c/f1.txt", "baz"},
	}
	expected := "a/b\n" +
		"a/b/c\n" +
		"a/b/c/d\n" +
		"\n" +
		"a/b2\n" +
		"a/b/c2\n" +
		"\n"

	createTempFiles(fdata)
	defer deleteTempFiles(fdata)

	if r := runForTempFiles(); r != expected {
		t.Fatalf("got:\n%#v\nexpected:\n%#v", r, expected)
	}
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createTempFiles(data []fileData) {
	tempdir, err := ioutil.TempDir("", "test")
	check(err)

	for _, v := range data {
		p := path.Join(tempdir, v.name)
		basedir := path.Dir(p)
		os.MkdirAll(basedir, 0755)
		err := ioutil.WriteFile(p, []byte(v.content), 0644)
		check(err)
	}
}

func deleteTempFiles(data []fileData) {
	err := os.RemoveAll(tempdir)
	check(err)
}

func runForTempFiles() string {
	out, err := exec.Command("go", "run", "main.go", ".").Output()
	check(err)
	return string(out)
}
