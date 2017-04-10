package dupfinder

import (
	"os"
	"io/ioutil"
)

type Group struct {
	Paths []string
}

type FileHandler interface {
	Id() string

	Size() int64

	Digest() string

	Content() []byte
	// TODO implement this and drop Content()
	//NewReader() io.Reader
}

type fileHandler struct {
	id      string
	size    int64
	digest  string
	content []byte
}

func NewFileHandler(path string, file os.FileInfo) FileHandler {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return fileHandler{
		id: path,
		size: file.Size(),
		digest: string(file.Size()),
		content: content,
	}
}

func (f fileHandler) Id() string {
	return f.id
}

func (f fileHandler) Size() int64 {
	return f.size
}

func (f fileHandler) Digest() string {
	return f.digest
}

func (f fileHandler) Content() []byte {
	return f.content
}

type Tracker interface {
	Add(FileHandler, FileHandler)

	Groups() []Group
}

type simpleTracker struct {
	groups map[string]Group
}

func (tracker *simpleTracker) Add(f1, f2 FileHandler) {
	group, found := tracker.groups[f1.Id()]
	if found {
		group.Paths = append(group.Paths, f2.Id())
	} else {
		group = Group{Paths: []string{f1.Id(), f2.Id()}}
	}
	tracker.groups[f1.Id()] = group
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

func (filter contentFilter) Match(f FileHandler, other FileHandler) bool {
	// TODO very very very bad comparison
	return string(f.Content()) == string(other.Content())
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
