package graph

type Node interface {
	Set(n Node)
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
