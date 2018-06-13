package pathreader

import (
	"bytes"
	"io"
	"bufio"
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

func readItems(reader io.Reader, splitter bufio.SplitFunc) <-chan string {
	items := make(chan string)

	go func() {
		scanner := bufio.NewScanner(reader)
		scanner.Split(splitter)
		for scanner.Scan() {
			item := scanner.Text()
			items <- item
		}
		close(items)
	}()

	return items
}

func readFilePaths(reader io.Reader, splitter bufio.SplitFunc) <-chan string {
	return readItems(reader, splitter)
}

func FromLines(reader io.Reader) <-chan string {
	return readFilePaths(reader, bufio.ScanLines)
}

func FromNullDelimited(reader io.Reader) <-chan string {
	return readFilePaths(reader, scanNullDelimited)
}
