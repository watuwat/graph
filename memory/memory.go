package memory

import (
	"strings"

	"watuwat.com/graph"
)

type MemoryNode struct {
	name     string
	edges    map[string]*MemoryNode
	readonly bool
}

var _ graph.Node = &MemoryNode{}
var _ graph.OnlyReader = &MemoryNode{}

func (m *MemoryNode) Path(p string) graph.Node {
	names := strings.Split(p, ".")

	var ok bool

	curr := m
	prev := curr
	for _, name := range names {
		curr, ok = curr.edges[name]
		if !ok {
			if m.readonly {
				return nil
			}

			tmp := New(name)
			prev.edges[name] = tmp
			curr = tmp
		}
		prev = curr
	}

	return curr
}

func (m *MemoryNode) Map(fn func(graph.Node) bool) {
	for key := range m.edges {
		if ok := fn(m.edges[key]); ok {
			break
		}
	}
}

func (m *MemoryNode) Size() int {
	return len(m.edges)
}

func (m *MemoryNode) Val() string {
	return m.name
}

func (m *MemoryNode) Readonly() graph.Node {
	return &MemoryNode{
		name:     m.name,
		edges:    m.edges,
		readonly: true,
	}
}

func New(name string) *MemoryNode {
	return &MemoryNode{
		name:     name,
		edges:    make(map[string]*MemoryNode),
		readonly: false,
	}
}
