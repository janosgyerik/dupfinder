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

type UniqueFilter interface {
	Add(string) bool
}

type uniqueFilter struct {
	seen map[string]bool
}

func (uf *uniqueFilter) Add(s string) bool {
	if _, ok := uf.seen[s]; !ok {
		uf.seen[s] = true
		return true
	}
	return false
}

func NewUniqueFilter() UniqueFilter {
	uf := &uniqueFilter{}
	uf.seen = make(map[string]bool)
	return uf
}
