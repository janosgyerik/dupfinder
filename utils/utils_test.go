package utils

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
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
	PanicIfFailed(err)

	err = ioutil.WriteFile(tempfile.Name(), make([]byte, size), 0644)
	PanicIfFailed(err)

	return tempfile.Name()
}

func TestIsFile(t *testing.T) {
	tempdir, err := ioutil.TempDir("", "test")
	PanicIfFailed(err)

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
		{"dir", dir, false},
		{"symlink to dir", linkToDir, false},
		{"nonexistent", "nonexistent", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFile(tt.path); got != tt.want {
				t.Errorf("IsFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPanicIfFailed(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		recovered bool
	}{
		{"success", nil, false},
		{"failure", io.EOF, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tt.recovered {
					if r == nil {
						t.Errorf("PanicIfFailed() did not panic but should have")
					}
				} else {
					if r != nil {
						t.Errorf("PanicIfFailed() paniced but should not have")
					}
				}
			}()
			PanicIfFailed(tt.err)
		})
	}
}

func Test_UniqueFilter(t *testing.T) {
	type AddWant struct {
		add  string
		want bool
	}
	tests := []struct {
		name string
		seq  []AddWant
	}{
		{"all different", []AddWant{
			{"foo", true},
			{"bar", true},
			{"baz", true}}},
		{"all same", []AddWant{
			{"foo", true},
			{"foo", false},
			{"foo", false}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uf := NewUniqueFilter()
			for _, item := range tt.seq {
				if got := uf.Add(item.add); got != item.want {
					t.Errorf("UniqueFilter.Add(%v) = %v, want %v", item.add, got, item.want)
				}
			}
		})
	}
}
