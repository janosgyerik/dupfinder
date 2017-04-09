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

type Tracker interface {
	Add(FileHandler, FileHandler)

	Groups() []Group
}

type myTracker struct {
	groups map[string]Group
}

func (tracker *myTracker) Add(f1, f2 FileHandler) {
	group, found := tracker.groups[f1.Id()]
	if found {
		group.Paths = append(group.Paths, f2.Id())
	} else {
		group = Group{Paths: []string{f1.Id(), f2.Id()}}
	}
	tracker.groups[f1.Id()] = group
}

func (tracker *myTracker) Groups() []Group {
	groups := []Group{}
	for _, group := range tracker.groups {
		groups = append(groups, group)
	}
	return groups
}

func NewTracker() Tracker {
	return &myTracker{
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

}

func (filter *SizeFilter) Accept(f FileHandler) bool {
	return false
}

type simpleIndex struct {
	files   []FileHandler

	filters []Filter

	groups  []Group

	tracker Tracker
}

func (index *simpleIndex) Add(f FileHandler) {
	files := index.files
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
	return []Group{}
}
