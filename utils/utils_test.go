package utils

import (
	"testing"
	"io/ioutil"
	"os"
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
