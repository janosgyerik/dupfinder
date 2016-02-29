package dupfinder

import (
	"os"
	"io"
	"bytes"
)

const chunkSize = 64000

func compareInts(a, b int) int {
	if a == b {
		return 0
	}
	if a < b {
		return -1
	}
	return 1
}

func Compare(path1, path2 string) (int, error) {
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

	for {
		buf1 := make([]byte, chunkSize)
		count1, err := fd1.Read(buf1)
		if err != nil && err != io.EOF{
			return 0, err
		}

		buf2 := make([]byte, chunkSize)
		count2, err := fd2.Read(buf2)
		if err != nil && err != io.EOF {
			return 0, err
		}

		if count1 != count2 {
			return compareInts(count1, count2), nil
		}

		cmp := bytes.Compare(buf1, buf2)
		if cmp != 0 {
			return cmp, nil
		}

		if err == io.EOF {
			return 0, nil
		}
	}
}
