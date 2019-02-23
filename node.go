package graph

type Node interface {
	Path(p string) Node
	Map(fn func(Node) bool)
	Size() int
	Val() string
}

type OnlyReader interface {
	Readonly() Node
}

func Readonly(node Node) Node {
	if n, ok := node.(OnlyReader); ok {
		return n.Readonly()
	}
	return nil
}
