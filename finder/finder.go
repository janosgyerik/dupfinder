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

type minSizeFilter struct {
	Size int64
}

func (filter minSizeFilter) Accept(path string, info os.FileInfo) bool {
	return info.Size() >= filter.Size
}

var Filters = struct {
	MinSize func(size int64) Filter
} {
	MinSize: func(size int64) Filter { return minSizeFilter{size} },
}

type defaultFinder struct {
	filters []Filter
}

func (finder defaultFinder) Find(basedir string) <-chan string {
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

func NewFinder(filters... Filter) Finder {
	return defaultFinder{filters: filters}
}
