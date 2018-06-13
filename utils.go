package dupfinder

import (
	"os"
	"path/filepath"
)

func FileSize(path string) int64 {
	fileInfo, e := os.Stat(path)
	if e != nil {
		panic(e)
	}
	return fileInfo.Size()
}

func NormalizePath(path string) string {
	return filepath.Clean(path)
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

func isFile(s string) bool {
	stat, err := os.Lstat(s)
	if err != nil {
		return false
	}
	return stat.Mode().IsRegular()
}

func NewDefaultFilter() PathFilter {
	isUnique := newUniqueFilter()

	return func(s string) bool {
		return isFile(s) && isUnique(s)
	}
}
