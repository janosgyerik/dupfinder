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
	Accepts(Item) bool
	Add(Item)
}

type Tracker interface {
	Add(Item)
	Dups() []Group
}

type Filter interface {
	CandidateGroups(Item) []Group
	Register(Item, Group)
}

type defaultTracker struct {
	groups []Group
	filter Filter
}

func (t *defaultTracker) Add(item Item) {
	candidates := t.filter.CandidateGroups(item)

	found := false
	for _, g := range candidates {
		if g.Accepts(item) {
			g.Add(item)
			found = true
			break
		}
	}

	if !found {
		group := newGroup(item)
		t.groups = append(t.groups, group)
		t.filter.Register(item, group)
	}
}

func (t *defaultTracker) Dups() []Group {
	dups := make([]Group, 0)
	for _, g := range t.groups {
		if len(g.Items()) > 1 {
			dups = append(dups, g)
		}
	}
	return dups
}

func NewTracker(filter Filter) Tracker {
	return &defaultTracker{filter: filter}
}

type defaultGroup struct {
	items []Item
}

func (g *defaultGroup) Items() []Item {
	return g.items
}

func (g *defaultGroup) Add(item Item) {
	g.items = append(g.items, item)
}

func (g *defaultGroup) Accepts(item Item) bool {
	return g.items[0].Equals(item)
}

func (g *defaultGroup) String() string {
	return fmt.Sprintf("%v", g.items)
}

func newGroup(item Item) Group {
	g := &defaultGroup{}
	g.Add(item)
	return g
}


type FileItem struct {
	path string
}

func (f *FileItem) Equals(other Item) bool {
	f2 := other.(*FileItem)

	s1, err1 := ioutil.ReadFile(f.path)
	if err1 != nil {
		panic("could not read file: " + f.path)
	}

	s2, err2 := ioutil.ReadFile(f2.path)
	if err2 != nil {
		panic("could not read file: " + f.path)
	}
	return string(s1) == string(s2)
}

type Key int

type KeyExtractor interface {
	Key(Item) Key
}

type sizeExtractor struct {
}

func (s *sizeExtractor) Key(item Item) Key {
	fi, e := os.Stat(item.(*FileItem).path)
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

func newFileFilter() Filter {
	return &defaultFilter{byKey: make(map[Key][]Group), keyExtractor: &sizeExtractor{}}
}
