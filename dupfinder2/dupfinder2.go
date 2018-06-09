package dupfinder2

import (
	"fmt"
	"io/ioutil"
	"os"
)

type Item interface {
	Equals(Item) bool
}

type Group interface {
	Items() []Item
	accepts(Item) bool
	add(Item)
}

type Tracker interface {
	Add(Item)
	Dups() []Group
}

type Filter interface {
	CandidateGroups(Item) []Group
	Register(Item, Group)
}

type tracker struct {
	groups []Group
	filter Filter
}

func (t *tracker) Add(item Item) {
	candidates := t.filter.CandidateGroups(item)

	for _, g := range candidates {
		if g.accepts(item) {
			g.add(item)
			return
		}
	}

	group := newGroup(item)
	t.groups = append(t.groups, group)
	t.filter.Register(item, group)
}

// TODO consistent deterministic ordering
func (t *tracker) Dups() []Group {
	dups := make([]Group, 0)
	for _, g := range t.groups {
		if len(g.Items()) > 1 {
			dups = append(dups, g)
		}
	}
	return dups
}

func NewTracker(filter Filter) Tracker {
	return &tracker{filter: filter}
}

type group struct {
	items []Item
}

func (g *group) Items() []Item {
	return g.items
}

func (g *group) add(item Item) {
	g.items = append(g.items, item)
}

func (g *group) accepts(item Item) bool {
	return g.items[0].Equals(item)
}

func (g *group) String() string {
	return fmt.Sprintf("%v", g.items)
}

func newGroup(item Item) Group {
	g := &group{}
	g.add(item)
	return g
}

type FileItem struct {
	Path string
}

func (f *FileItem) Equals(other Item) bool {
	f2 := other.(*FileItem)

	s1, err1 := ioutil.ReadFile(f.Path)
	if err1 != nil {
		panic("could not read file: " + f.Path)
	}

	s2, err2 := ioutil.ReadFile(f2.Path)
	if err2 != nil {
		panic("could not read file: " + f.Path)
	}
	return string(s1) == string(s2)
}

func NewFileItem(path string) Item {
	return &FileItem{path}
}

type Key int

type KeyExtractor interface {
	Key(Item) Key
}

type sizeExtractor struct {
}

func (s *sizeExtractor) Key(item Item) Key {
	fi, e := os.Stat(item.(*FileItem).Path)
	if e != nil {
		return 0
	}
	return Key(fi.Size())
}

type defaultFilter struct {
	byKey        map[Key][]Group
	keyExtractor KeyExtractor
}

func (f *defaultFilter) CandidateGroups(item Item) []Group {
	if g, found := f.byKey[f.keyExtractor.Key(item)]; found {
		return g
	}
	return nil
}

func (f *defaultFilter) Register(item Item, g Group) {
	f.byKey[f.keyExtractor.Key(item)] = append(f.byKey[f.keyExtractor.Key(item)], g)
}

func NewFileFilter() Filter {
	return &defaultFilter{byKey: make(map[Key][]Group), keyExtractor: &sizeExtractor{}}
}
