package main

import (
	"flag"
	"os"
	"github.com/janosgyerik/dupfinder/pathreader"
	"github.com/janosgyerik/dupfinder/finder"
	"fmt"
	"github.com/janosgyerik/dupfinder"
)

var verbose bool

const defaultExclude = `^\.(DS_Store|git)$`

func exit() {
	flag.Usage()
	os.Exit(1)
}

type Params struct {
	paths   <-chan string
	minSize int64
	stdin   bool
	stdin0  bool
	verbose bool
	include []string
	exclude []string
}

func parseArgs() Params {
	minSizePtr := flag.Int64("minSize", 1, "minimum file size")
	stdinPtr := flag.Bool("stdin", false, "read paths from stdin")
	zeroPtr := flag.Bool("0", false, "read paths from stdin, null-delimited")
	silentPtr := flag.Bool("silent", false, "silent mode, do not print stats on stderr")
	includePtr := flag.String("include", ".", "include file paths that match regex")
	excludePtr := flag.String("exclude", defaultExclude, "exclude file paths that match regex")

	flag.Parse()

	var paths <-chan string
	if *zeroPtr {
		paths = pathreader.FromNullDelimited(os.Stdin)
	} else if *stdinPtr {
		paths = pathreader.FromLines(os.Stdin)
	} else if len(flag.Args()) > 0 {
		filters := []finder.Filter{
			finder.Filters.MinSize(*minSizePtr),
			finder.Filters.IncludeRegex(*includePtr),
			finder.Filters.ExcludeRegex(*excludePtr),
		}
		filefinder := finder.NewFinder(filters...)
		paths = findInAll(filefinder, flag.Args())
	} else {
		exit()
	}

	return Params{
		paths:   paths,
		minSize: *minSizePtr,
		verbose: !*silentPtr,
	}
}

func findInAll(f finder.Finder, args []string) <-chan string {
	agg := make(chan string)
	go func() {
		for _, path := range args {
			for msg := range f.Find(path) {
				agg <- msg
			}
		}
		close(agg)
	}()

	return agg
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

func (log *eventLogger) NewDuplicate(items []*dupfinder.FileItem, item *dupfinder.FileItem) {
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

	printLine("Collecting paths to check ...")

	filter := dupfinder.NewDefaultFilter()
	var paths []string
	i := 1
	for path := range params.paths {
		normalized := dupfinder.NormalizePath(path)
		if !filter(normalized) {
			continue
		}
		paths = append(paths, normalized)
		status("Found: %d", i)
		i += 1
	}
	printLine()

	tracker := dupfinder.NewTracker()
	logger := eventLogger{}
	tracker.SetLogger(&logger)

	i = 1
	for _, path := range paths {
		tracker.Add(path)
		status("Processing: %d / %d", i, len(paths))
		i += 1
	}
	printLine()

	for _, group := range tracker.Dups() {
		fmt.Println("# file sizes:", dupfinder.FileSize(group[0]))
		for _, path := range group {
			fmt.Println(path)
		}
		fmt.Println()
	}

	printLine("Total bytes read:", logger.bytesRead)
	printLine("Total files processed:", len(paths))
}
