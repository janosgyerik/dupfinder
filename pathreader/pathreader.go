package pathreader

import (
	"bytes"
	"io"
	"bufio"
	"os"
	"path/filepath"
)

func scanNullDelimited(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, 0); i >= 0 {
		// We have a full null-terminated line.
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

type filter func(string) bool

func newUniqueFilter() filter {
	seen := make(map[string]bool)

	return func(s string) bool {
		normalized := normalize(s)
		if _, ok := seen[normalized]; !ok {
			seen[normalized] = true
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

func newDefaultFilter() filter {
	isUnique := newUniqueFilter()

	return func(s string) bool {
		return isFile(s) && isUnique(s)
	}
}

func readItems(reader io.Reader, splitter bufio.SplitFunc, filter filter) <-chan string {
	items := make(chan string)

	go func() {
		scanner := bufio.NewScanner(reader)
		scanner.Split(splitter)
		for scanner.Scan() {
			if item := scanner.Text(); filter(item) {
				items <- normalize(item)
			}
		}
		close(items)
	}()

	return items
}

func normalize(path string) string {
	return filepath.Clean(path)
}

func readItemsFromLines(reader io.Reader, filter filter) <-chan string {
	return readItems(reader, bufio.ScanLines, filter)
}

func readItemsFromNullDelimited(reader io.Reader, filter filter) <-chan string {
	return readItems(reader, scanNullDelimited, filter)
}

func readFilePaths(reader io.Reader, splitter bufio.SplitFunc) <-chan string {
	return readItems(reader, splitter, newDefaultFilter())
}

func FromLines(reader io.Reader) <-chan string {
	return readFilePaths(reader, bufio.ScanLines)
}

func FromNullDelimited(reader io.Reader) <-chan string {
	return readFilePaths(reader, scanNullDelimited)
}
