package utils

import (
	"io/ioutil"
	"os"
	"testing"
	"path/filepath"
)

func TestFileSize(t *testing.T) {
	tests := []struct {
		name string
		want int64
	}{
		{"empty", 0},
		{"size 1", 1},
		{"size 5", 5},
	}
	for _, tt := range tests {
		f := newTempFile(tt.want)
		defer os.Remove(f)

		t.Run(tt.name, func(t *testing.T) {
			if got := FileSize(f); got != tt.want {
				t.Errorf("FileSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func newTempFile(size int64) string {
	tempfile, err := ioutil.TempFile("", "test")
	check(err)

	err = ioutil.WriteFile(tempfile.Name(), make([]byte, size), 0644)
	check(err)

	return tempfile.Name()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func TestIsFile(t *testing.T) {
	tempdir, err := ioutil.TempDir("", "test")
	check(err)

	defer os.RemoveAll(tempdir)

	file := filepath.Join(tempdir, "file")
	ioutil.WriteFile(file, []byte("foo"), 0644)

	linkToFile := filepath.Join(tempdir, "link-to-file")
	os.Symlink("file", linkToFile)

	dir := filepath.Join(tempdir, "dir")
	os.Mkdir(dir, 0755)

	linkToDir := filepath.Join(tempdir, "link-to-dir")
	os.Symlink("dir", linkToDir)

	tests := []struct {
		name string
		path string
		want bool
	}{
		{"file", file, true},
		{"symlink to file", linkToFile, false},
		{"dir", dir,false},
		{"symlink to dir", linkToDir, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFile(tt.path); got != tt.want {
				t.Errorf("IsFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
