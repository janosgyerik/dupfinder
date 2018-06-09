package finder

import (
	"path/filepath"
	"os"
	"regexp"
)

type Finder interface {
	Find(basedir string) <-chan string
	FindAll(basedir string) []string
}

type Filter interface {
	Accept(path string, info os.FileInfo) bool
}

type minSizeFilter struct {
	size int64
}

func (filter minSizeFilter) Accept(path string, info os.FileInfo) bool {
	return info.Size() >= filter.size
}

type regexFilter struct {
	regex *regexp.Regexp
}

func newRegexFilter(regex string) regexFilter {
	return regexFilter{regexp.MustCompile(regex)}
}

func (filter regexFilter) Accept(path string, info os.FileInfo) bool {
	return filter.regex.MatchString(path)
}

var Filters = struct {
	MinSize      func(size int64) Filter
	IncludeRegex func(pattern string) Filter
}{
	MinSize:      func(size int64) Filter { return minSizeFilter{size} },
	IncludeRegex: func(regex string) Filter { return newRegexFilter(regex) },
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
		if info.Mode().IsRegular() {
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

func (finder defaultFinder) FindAll(basedir string) []string {
	var paths []string
	for path := range finder.Find(basedir) {
		paths = append(paths, path)
	}
	return paths
}

func NewFinder(filters ... Filter) Finder {
	return defaultFinder{filters: filters}
}
