package finder

import (
	"path/filepath"
	"os"
	"regexp"
)

type Finder interface {
	Find(basedir string) <-chan string
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
	regex    *regexp.Regexp
	negative bool
}

func newRegexFilter(regex string, negative bool) regexFilter {
	return regexFilter{regexp.MustCompile(regex), negative}
}

func (filter regexFilter) Accept(path string, info os.FileInfo) bool {
	return filter.negative != filter.regex.MatchString(path)
}

var Filters = struct {
	MinSize      func(size int64) Filter
	IncludeRegex func(pattern string) Filter
	ExcludeRegex func(pattern string) Filter
}{
	MinSize:      func(size int64) Filter { return minSizeFilter{size} },
	IncludeRegex: func(regex string) Filter { return newRegexFilter(regex, false) },
	ExcludeRegex: func(regex string) Filter { return newRegexFilter(regex, true) },
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

func NewFinder(filters ... Filter) Finder {
	return defaultFinder{filters: filters}
}
