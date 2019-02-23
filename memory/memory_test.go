package memory_test

import (
	"testing"

	"watuwat.com/graph"
	"watuwat.com/graph/memory"
)

func TestMemoryGraph(t *testing.T) {
	root := memory.New("root")
	users := root.Path("users")
	user1 := root.Path("users.1")
	user2 := root.Path("users.2")

	if root.Size() != 1 {
		t.Fatalf("expect root size to be 1 but got %d", root.Size())
	}

	if users.Size() != 2 {
		t.Fatalf("expect users size to be 2 but got %d", users.Size())
	}

	if user1.Val() != "1" {
		t.Fatalf("expect user1's val to be 1 but got %s", user1.Val())
	}

	if user2.Val() != "2" {
		t.Fatalf("expect user2's val to be 2 but got %s", user2.Val())
	}

	count := 0
	users.Map(func(n graph.Node) bool {
		count++
		return false
	})

	if count != 2 {
		t.Fatalf("expect users map to be called 2 times but got %d", count)
	}

	count = 0
	users.Map(func(n graph.Node) bool {
		count++
		return true
	})

	if count != 1 {
		t.Fatalf("expect users map to be called 1 times since we return true but got %d", count)
	}

	var ids []string
	users.Map(func(n graph.Node) bool {
		ids = append(ids, n.Val())
		return false
	})

	if len(ids) != 2 {
		t.Fatalf("expect users map to be called 2 times with val but got %d", count)
	}

	if ids[0] != "1" {
		t.Fatalf("expect first edge to be 1 but got %s", ids[0])
	}

	if ids[1] != "2" {
		t.Fatalf("expect second edge to be 2 but got %s", ids[1])
	}

	readonlyRoot := graph.Readonly(root)
	if readonlyRoot == nil {
		t.Fatalf("expect root to be MemoryGraph which implements OnlyReader interface but got nil")
	}

	user3 := readonlyRoot.Path("users.3")
	if user3 != nil {
		t.Fatalf("expect readonly root not create a new path but apparently created one")
	}

	if readonlyRoot.Size() != 1 {
		t.Fatalf("expect readonly root size to be 1 but got %d", readonlyRoot.Size())
	}
}
