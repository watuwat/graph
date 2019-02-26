package graph

import (
	"strings"
)

type Bucket interface {
	Get(name string) (*Node, bool)
	Set(name string, node *Node)
	Map(fn func(*Node) bool)
	Size() int
}

type Storage interface {
	Bucket(name string) Bucket
	Close() error
}

type Node struct {
	name     string
	storage  Storage
	edges    Bucket
	readonly bool
}

func (m *Node) Path(p string) *Node {
	names := strings.Split(p, ".")

	var ok bool

	current := m
	prev := current
	for _, name := range names {
		current, ok = current.edges.Get(name)
		if !ok {
			if m.readonly {
				return nil
			}

			tmp := NewNode(name, m.storage)
			prev.edges.Set(name, tmp)
			current = tmp
		}
		prev = current
	}

	return current
}

func (m *Node) Set(n *Node) {
	if m.readonly {
		return
	}

	_, ok := m.edges.Get(n.name)
	if ok {
		return
	}

	m.edges.Set(n.name, n)
}

func (m *Node) Map(fn func(*Node) bool) {
	m.edges.Map(func(node *Node) bool {
		if m.readonly {
			node = node.Readonly()
		}
		return fn(node)
	})
}

func (m *Node) Size() int {
	return m.edges.Size()
}

func (m *Node) Val() string {
	return m.name
}

func (m *Node) Readonly() *Node {
	if m.readonly {
		return m
	}

	return &Node{
		name:     m.name,
		storage:  m.storage,
		edges:    m.edges,
		readonly: true,
	}
}

func NewNode(name string, storage Storage) *Node {
	return &Node{
		name:     name,
		storage:  storage,
		edges:    storage.Bucket(name),
		readonly: false,
	}
}
