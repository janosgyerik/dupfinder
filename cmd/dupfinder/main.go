package main

import (
	"flag"
	"os"
	"github.com/janosgyerik/dupfinder/pathreader"
	"github.com/janosgyerik/dupfinder/finder"
	"fmt"
	"github.com/janosgyerik/dupfinder"
	"strconv"
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
	minSizePtr := flag.String("minSize", "100m", "minimum file size")
	stdinPtr := flag.Bool("stdin", false, "read paths from stdin")
	zeroPtr := flag.Bool("0", false, "read paths from stdin, null-delimited")
	silentPtr := flag.Bool("silent", false, "silent mode, do not print stats on stderr")
	includePtr := flag.String("include", ".", "include file paths that match regex")
	excludePtr := flag.String("exclude", defaultExclude, "exclude file paths that match regex")

	flag.Parse()

	minSize := toByteCount(*minSizePtr)

	var paths <-chan string
	if *zeroPtr {
		paths = pathreader.FromNullDelimited(os.Stdin)
	} else if *stdinPtr {
		paths = pathreader.FromLines(os.Stdin)
	} else if len(flag.Args()) > 0 {
		filters := []finder.Filter{
			finder.Filters.MinSize(minSize),
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
		minSize: minSize,
		verbose: !*silentPtr,
	}
}

func toByteCount(s string) int64 {
	numPart := s[0 : len(s)-1]
	unitPart := s[len(s)-1]

	var multiplier int64 = 1

	switch unitPart {
	case 'c':
		multiplier = 1
	case 'k':
		fallthrough
	case 'K':
		multiplier = 1<<10
	case 'm':
		fallthrough
	case 'M':
		multiplier = 1<<20
	case 'g':
		fallthrough
	case 'G':
		multiplier = 1<<30
	case 't':
		fallthrough
	case 'T':
		multiplier = 1<<40
	default:
		numPart = s
	}

	v, err := strconv.Atoi(numPart)
	check(err)
	return int64(v) * multiplier
}

func check(e error) {
	if e != nil {
		panic(e)
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

func (log *eventLogger) NewDuplicate(paths []string) {
	printLine()
	for _, path := range paths {
		printLine(path)
	}
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
