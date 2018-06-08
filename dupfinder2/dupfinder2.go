package dupfinder2

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

type defaultTracker struct {
	groups []Group
}

func (t *defaultTracker) Add(item Item) {
	found := false
	for _, g := range t.groups {
		if g.Accepts(item) {
			g.Add(item)
			found = true
			break
		}
	}

	if !found {
		t.groups = append(t.groups, newGroup(item))
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

func NewTracker() Tracker {
	return &defaultTracker{}
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

func newGroup(item Item) Group {
	g := &defaultGroup{}
	g.Add(item)
	return g
}