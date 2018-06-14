package dupfinder

import (
	"testing"
	"io/ioutil"
	"path"
	"os"
	"reflect"
	"github.com/janosgyerik/dupfinder/utils"
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

	if actual := normalize(run(fdata)); !reflect.DeepEqual(expected, actual) {
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

	if actual := normalize(run(fdata)); !reflect.DeepEqual(expected, actual) {
		t.Fatalf("got:\n%#v\nexpected:\n%#v", actual, expected)
	}
}

func Test_find_multiple_varied_sized_groups(t *testing.T) {
	fdata := []fileData{
		{"z1.txt", "zo"},
		{"z2.txt", "zo"},
		{"size5-1.txt", "apple"},
		{"size5-2.txt", "apple"},
		{"f1.txt", "foo"},
		{"a/f2.txt", "foo"},
		{"a/b/f3.txt", "foo"},
		{"b/f1.txt", "bar"},
		{"b/c/f2.txt", "bar"},
		{"c/f1.txt", "baz"},
	}
	expected := [][]string{
		{
			"z1.txt",
			"z2.txt",
		},
		{
			"a/b/f3.txt",
			"a/f2.txt",
			"f1.txt",
		},
		{
			"b/c/f2.txt",
			"b/f1.txt",
		},
		{
			"size5-1.txt",
			"size5-2.txt",
		},
	}

	createTempFiles(fdata)
	defer deleteTempFiles()

	if actual := normalize(run(fdata)); !reflect.DeepEqual(expected, actual) {
		t.Fatalf("got:\n%#v\nexpected:\n%#v", actual, expected)
	}
}

func normalize(groups [][]string) [][]string {
	var result [][]string
	for _, g := range groups {
		var stripped []string
		for _, p := range g {
			stripped = append(stripped, p[len(tempdir)+1:])
		}
		result = append(result, stripped)
	}
	return result
}

func createTempFiles(data []fileData) {
	var err error
	tempdir, err = ioutil.TempDir("", "test")
	utils.PanicIfFailed(err)

	for _, v := range data {
		p := path.Join(tempdir, v.relpath)
		basedir := path.Dir(p)
		os.MkdirAll(basedir, 0755)
		err := ioutil.WriteFile(p, []byte(v.content), 0644)
		utils.PanicIfFailed(err)
	}
}

func deleteTempFiles() {
	err := os.RemoveAll(tempdir)
	utils.PanicIfFailed(err)
}

func run(fdata []fileData) [][]string {
	t := NewTracker()
	t.SetEventListener(&nullEventListener{})
	for _, v := range fdata {
		t.Add(path.Join(tempdir, v.relpath))
	}
	return t.Dups()
}
