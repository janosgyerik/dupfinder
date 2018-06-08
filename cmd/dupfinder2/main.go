package main

import (
	"flag"
	"os"
	"fmt"
	"github.com/janosgyerik/dupfinder/finder"
	"github.com/janosgyerik/dupfinder/pathreader"
	"github.com/janosgyerik/dupfinder/dupfinder2"
)

func exit() {
	flag.Usage()
	os.Exit(1)
}

type Params struct {
	paths   []string
	minSize int64
	stdin   bool
	stdin0  bool
}

func parseArgs() Params {
	minSizePtr := flag.Int64("minSize", 1, "minimum file size")
	stdinPtr := flag.Bool("stdin", false, "read paths from stdin")
	zeroPtr := flag.Bool("0", false, "read paths from stdin, null-delimited")

	flag.Parse()

	var paths []string
	if *zeroPtr {
		paths = pathreader.ReadFilePathsFromNullDelimited(os.Stdin)
	} else if *stdinPtr {
		paths = pathreader.ReadFilePathsFromLines(os.Stdin)
	} else {
		paths = pathreader.FilterPaths(flag.Args())
	}

	if len(paths) == 0 {
		exit()
	}

	return Params{
		paths:   paths,
		minSize: *minSizePtr,
	}
}

func main() {
	params := parseArgs()

	filefinder := finder.NewFinder(finder.Filters.MinSize(params.minSize))

	paths := []string{}
	for _, path := range params.paths {
		paths = append(paths, filefinder.FindAll(path)...)
	}

	tracker := dupfinder2.NewTracker(dupfinder2.NewFileFilter())
	for _, path := range paths {
		tracker.Add(dupfinder2.NewFileItem(path))
	}

	// TODO consistent deterministic ordering
	for _, dup := range tracker.Dups() {
		for _, item := range dup.Items() {
			fmt.Println(item.(*dupfinder2.FileItem).Path)
		}
		fmt.Println()
	}
}
