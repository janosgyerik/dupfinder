package dupfinder

import (
	"os"
	"bytes"
	"io"
)

const chunkSize = 64000

func chunker(r io.Reader, ch chan <- []byte) {
	for {
		buf := make([]byte, chunkSize)
		_, err := r.Read(buf)

		if err != nil {
			close(ch)
			return
		}

		ch <- buf
	}
}

func CompareFiles(path1, path2 string) (int, error) {
	fd1, err := os.Open(path1)
	defer fd1.Close()
	if err != nil {
		return 0, err
	}

	fd2, err := os.Open(path2)
	defer fd2.Close()
	if err != nil {
		return 0, err
	}

	return CompareReaders(fd1, fd2)
}

func CompareReaders(fd1, fd2 io.Reader) (int, error) {
	ch1 := make(chan []byte)
	ch2 := make(chan []byte)

	go chunker(fd1, ch1)
	go chunker(fd2, ch2)

	for {
		buf1 := <-ch1
		buf2 := <-ch2

		cmp := bytes.Compare(buf1, buf2)
		if cmp != 0 {
			return cmp, nil
		}

		if len(buf1) == 0 {
			return 0, nil
		}
	}
}
