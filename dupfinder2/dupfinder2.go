package dupfinder2

import "fmt"

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