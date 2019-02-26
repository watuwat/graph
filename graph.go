package graph

import (
	"strings"

	"github.com/rs/xid"
)

/*

	# 1

	if node.isRoot {

	}

	- if node is root, we need to add name as key and generate a uuid as value
	-
	- prev.Set(name, uuid)

	root := graph.New(storage)

	root.Path("a.b")

	root: {
		a: <id-a>
		<id-b>: <id-b>
	}

	<id-a>: {
		b: <id-b>
	}

	<id-b>: {

	}

	# 2

	root := graph.New(storage)

	a := root.Path("a")
	b := root.Path("b")

	a.Set(b)

	root : {
		a: <id-a>
		b: <id-b>
	}

	<id-a>: {
		b: <id-b>
	}

	<id-b>: {}

*/

func uuid() string {
	return xid.New().String()
}

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
	id       string
	name     string
	storage  Storage
	edges    Bucket
	readonly bool
	isRoot   bool
}

func (m *Node) Path(p string) *Node {
	var ok bool

	names := strings.Split(p, ".")
	root := m.storage.Bucket("root")

	current := m
	prev := current
	for _, name := range names {
		current, ok = current.edges.Get(name)
		if !ok {
			if m.readonly {
				return nil
			}

			id := uuid()
			node := newNode(id, name, false, m.storage)

			// we need to create a node
			if prev.isRoot {
				root.Set(name, node)
			} else {
				root.Set(id, node)
				prev.edges.Set(name, node)
			}

			current = node
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

func newNode(id string, name string, isRoot bool, storage Storage) *Node {
	return &Node{
		id:       id,
		name:     name,
		storage:  storage,
		edges:    storage.Bucket(id),
		readonly: false,
		isRoot:   isRoot,
	}
}

func New(storage Storage) *Node {
	return newNode("root", "root", true, storage)
}
