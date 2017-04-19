package dupfinder

import (
	"os"
	"bytes"
	"io"
)

const chunkSize = 64000

type Group struct {
	Paths []string
}

type FileHandler interface {
	Path() string

	Size() int64

	Digest() string

	NewReadCloser() (io.ReadCloser, error)
}

type fileHandler struct {
	path   string
	size   int64
	digest string
}

func NewFileHandler(path string, file os.FileInfo) FileHandler {
	return fileHandler{
		path: path,
		size: file.Size(),
		digest: string(file.Size()),
	}
}

func (f fileHandler) Path() string {
	return f.path
}

func (f fileHandler) Size() int64 {
	return f.size
}

func (f fileHandler) Digest() string {
	return f.digest
}

func (f fileHandler) NewReadCloser() (io.ReadCloser, error) {
	return os.Open(f.Path())
}

type Tracker interface {
	Add(FileHandler, FileHandler)

	Groups() []Group
}

type simpleTracker struct {
	groups map[string]Group
}

func (tracker *simpleTracker) Add(f1, f2 FileHandler) {
	group, found := tracker.groups[f1.Path()]
	if found {
		group.Paths = append(group.Paths, f2.Path())
	} else {
		group = Group{Paths: []string{f1.Path(), f2.Path()}}
	}
	tracker.groups[f1.Path()] = group
}

func (tracker *simpleTracker) Groups() []Group {
	groups := []Group{}
	for _, group := range tracker.groups {
		groups = append(groups, group)
	}
	return groups
}

func NewTracker() Tracker {
	return &simpleTracker{
		groups: make(map[string]Group),
	}
}

type Index interface {
	Add(FileHandler)

	Groups() []Group
}

type Filter interface {
	Match(base, other FileHandler) bool
}

type sizeFilter struct{}

func (filter sizeFilter) Match(f FileHandler, other FileHandler) bool {
	return f.Size() == other.Size()
}

type digestFilter struct{}

func (filter digestFilter) Match(f FileHandler, other FileHandler) bool {
	return f.Digest() == other.Digest()
}

type contentFilter struct{}

// TODO return errors, let caller handle
func (filter contentFilter) Match(f FileHandler, other FileHandler) bool {
	fd1, err := f.NewReadCloser()
	defer fd1.Close()
	if err != nil {
		return false
	}

	fd2, err := other.NewReadCloser()
	defer fd2.Close()
	if err != nil {
		return false
	}

	cmp, err := CompareReaders(fd1, fd2)
	if err != nil {
		return false
	}
	return cmp == 0
}

type simpleIndex struct {
	files   []FileHandler
	filters []Filter
	tracker Tracker
}

func (index *simpleIndex) Add(f FileHandler) {
	files := index.files
	for _, filter := range index.filters {
		files = applyFilter(filter, files, f)
	}

	switch len(files) {
	case 0:
		index.files = append(index.files, f)
	case 1:
		index.tracker.Add(files[0], f)
	default:
		// TODO return error instead, let caller handle (+ unit test it)
		panic("more than one duplicates found in the unique index")
	}
}

func applyFilter(filter Filter, files []FileHandler, base FileHandler) []FileHandler {
	filtered := []FileHandler{}
	for _, file := range files {
		if filter.Match(base, file) {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

func (index *simpleIndex) Groups() []Group {
	return index.tracker.Groups()
}

func NewIndex() Index {
	index := simpleIndex{
		files: []FileHandler{},
		filters: []Filter{
			sizeFilter{},
			digestFilter{},
			contentFilter{},
		},
		tracker: NewTracker(),
	}
	return &index
}

func chunker(r io.Reader, ch chan <- []byte) {
	buf := make([]byte, chunkSize)
	for {
		n, err := r.Read(buf)

		if err != nil {
			if err == io.EOF {
				ch <- buf[:n]
			}
			close(ch)
			return
		}

		ch <- buf
	}
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
