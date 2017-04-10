package dupfinder

type Group struct {
	Paths []string
}

type FileHandler interface {
	Id() string

	Size() int

	Digest() string

	Content() string
	// TODO implement this and drop Content()
	//NewReader() io.Reader
}

type fileHandler struct {
	id      string
	size    int
	digest  string
	content string
}

func NewFileHandler(id string) FileHandler {
	return fileHandler{
		id: id,
	}
}

func (f fileHandler) Id() string {
	return f.id
}

func (f fileHandler) Size() int {
	return f.size
}

func (f fileHandler) Digest() string {
	return f.digest
}

func (f fileHandler) Content() string {
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
	Accept(FileHandler) bool
}

type SizeFilter struct {
	size int
}

func (filter SizeFilter) Accept(f FileHandler) bool {
	return filter.size == f.Size()
}

type DigestFilter struct {
	digest string
}

func (filter DigestFilter) Accept(f FileHandler) bool {
	return filter.digest == f.Digest()
}

type ContentFilter struct {
	content string
}

func (filter ContentFilter) Accept(f FileHandler) bool {
	return filter.content == f.Content()
}

type simpleIndex struct {
	files   []FileHandler
	filters []Filter
	tracker Tracker
}

func (index *simpleIndex) Add(f FileHandler) {
	files := index.files
	index.filters = []Filter{
		SizeFilter{f.Size()},
		DigestFilter{f.Digest()},
		ContentFilter{f.Content()},
	}
	for _, filter := range index.filters {
		files = applyFilter(filter, files)
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

func applyFilter(filter Filter, files []FileHandler) []FileHandler {
	filtered := []FileHandler{}
	for _, file := range files {
		if filter.Accept(file) {
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
		filters: []Filter{},
		tracker: NewTracker(),
	}
	return &index
}
