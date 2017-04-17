package finder

import (
	"path/filepath"
	"os"
)

type Finder interface {
	Find(basedir string) <-chan string
}

type Filter interface {
	Accept(path string, info os.FileInfo) bool
}

type MinSizeFilter struct {
	Size int64
}

func (filter MinSizeFilter) Accept(path string, info os.FileInfo) bool {
	return info.Size() >= filter.Size
}

type DefaultFinder struct {
	filters []Filter
}

func (finder DefaultFinder) Find(basedir string) <-chan string {
	paths := make(chan string)
	walkfn := func(path string, info os.FileInfo, err error) error {
		for _, filter := range finder.filters {
			if !filter.Accept(path, info) {
				return nil
			}
		}
		if !info.IsDir() {
			paths <- path
		}
		return nil
	}
	go func() {
		filepath.Walk(basedir, walkfn)
		close(paths)
	}()
	return paths
}
