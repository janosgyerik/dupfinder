package utils

import (
	"os"
	)

func FileSize(path string) int64 {
	fileInfo, e := os.Stat(path)
	PanicIfFailed(e)
	return fileInfo.Size()
}

func IsFile(s string) bool {
	stat, err := os.Lstat(s)
	if err != nil {
		return false
	}
	return stat.Mode().IsRegular()
}

func PanicIfFailed(e error) {
	if e != nil {
		panic(e)
	}
}

type PathFilter func(string) bool

func newUniqueFilter() PathFilter {
	seen := make(map[string]bool)

	return func(s string) bool {
		if _, ok := seen[s]; !ok {
			seen[s] = true
			return true
		}
		return false
	}
}

func NewDefaultFilter() PathFilter {
	isUnique := newUniqueFilter()

	return func(s string) bool {
		return isUnique(s)
	}
}
