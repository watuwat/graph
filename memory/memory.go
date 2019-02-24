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

func (m *MemoryNode) Set(n graph.Node) {
	if m.readonly {
		return
	}

	given, ok := n.(*MemoryNode)
	if !ok {
		panic("given node in Set is not MemoryNode")
	}

	_, ok = m.edges[given.name]
	if ok {
		return
	}

	m.edges[given.name] = given
}

func (m *MemoryNode) Map(fn func(graph.Node) bool) {
	for key := range m.edges {
		node := m.edges[key]
		if m.readonly {
			node = node.Readonly().(*MemoryNode)
		}

		if ok := fn(node); ok {
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
	if m.readonly {
		return m
	}

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
