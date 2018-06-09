package main

import (
	"flag"
	"os"
	"fmt"
	"github.com/janosgyerik/dupfinder/finder"
	"github.com/janosgyerik/dupfinder/pathreader"
	"github.com/janosgyerik/dupfinder/dupfinder3"
)

var verbose bool

func exit() {
	flag.Usage()
	os.Exit(1)
}

type Params struct {
	paths   []string
	minSize int64
	stdin   bool
	stdin0  bool
	verbose bool
}

func parseArgs() Params {
	minSizePtr := flag.Int64("minSize", 1, "minimum file size")
	stdinPtr := flag.Bool("stdin", false, "read paths from stdin")
	zeroPtr := flag.Bool("0", false, "read paths from stdin, null-delimited")
	verbosePtr := flag.Bool("verbose", false, "verbose mode, print stats on stderr")

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
		verbose: *verbosePtr,
	}
}

func printLine(args ...interface{}) {
	if !verbose {
		return
	}
	fmt.Fprintln(os.Stderr, args...)
}

func status(first string, args ...interface{}) {
	if !verbose {
		return
	}
	fmt.Fprintf(os.Stderr, "\r"+first, args...)
}

type eventLogger struct {
	bytesRead int64
}

func (log *eventLogger) NewDuplicate(items []*dupfinder3.FileItem, item *dupfinder3.FileItem) {
	printLine()
	for _, oldItem := range items {
		printLine(oldItem.Path)
	}
	printLine("->", item.Path, item.Size)
	printLine()
}

func (log *eventLogger) BytesRead(count int) {
	log.bytesRead += int64(count)
}

func main() {
	params := parseArgs()

	verbose = params.verbose

	filefinder := finder.NewFinder(finder.Filters.MinSize(params.minSize))

	printLine("Collecting paths to check ...")

	var paths []string
	for _, path := range params.paths {
		paths = append(paths, filefinder.FindAll(path)...)
	}

	printLine("Files:", len(paths))

	tracker := dupfinder3.NewTracker()
	logger := eventLogger{}
	tracker.SetLogger(&logger)

	i := 1
	for _, path := range paths {
		tracker.Add(path)
		status("Processing %d / %d", i, len(paths))
		i += 1
	}

	printLine()

	for _, group := range tracker.Dups() {
		fmt.Println("# file sizes:", dupfinder3.FileSize(group[0]))
		for _, path := range group {
			fmt.Println(path)
		}
		fmt.Println()
	}

	printLine("Total bytes read:", logger.bytesRead)
}
