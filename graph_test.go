package graph_test

import (
	"os"
	"testing"

	"watuwat.com/graph"
)

var createStorage func() graph.Storage

func indexOf(element string, data []string) int {
	for i, v := range data {
		if element == v {
			return i
		}
	}
	return -1 //not found.
}

func TestSize(t *testing.T) {
	storage := createStorage()
	root := graph.NewNode("root", storage)

	users := root.Path("users")

	root.Path("users.1")
	root.Path("users.2")

	if root.Size() != 1 {
		t.Fatalf("expect root size to be 1 but got %d", root.Size())
	}

	if users.Size() != 2 {
		t.Fatalf("expect users size to be 2 but got %d", users.Size())
	}

	err := storage.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestVal(t *testing.T) {
	storage := createStorage()
	root := graph.NewNode("root", storage)

	users := root.Path("users")
	user1 := users.Path("1")
	user2 := users.Path("2")

	if user1.Val() != "1" {
		t.Fatalf("expect user1's val to be 1 but got %s", user1.Val())
	}

	if user2.Val() != "2" {
		t.Fatalf("expect user2's val to be 2 but got %s", user2.Val())
	}

	err := storage.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestMap(t *testing.T) {
	storage := createStorage()
	root := graph.NewNode("root", storage)

	users := root.Path("users")

	root.Path("users.1")
	root.Path("users.2")

	count := 0
	users.Map(func(n *graph.Node) bool {
		count++
		return false
	})

	if count != 2 {
		t.Fatalf("expect users map to be called 2 times but got %d", count)
	}

	count = 0
	users.Map(func(n *graph.Node) bool {
		count++
		return true
	})

	if count != 1 {
		t.Fatalf("expect users map to be called 1 times since we return true but got %d", count)
	}

	var ids []string
	users.Map(func(n *graph.Node) bool {
		ids = append(ids, n.Val())
		return false
	})

	if len(ids) != 2 {
		t.Fatalf("expect users map to be called 2 times with val but got %d", count)
	}

	if indexOf("1", ids) == -1 {
		t.Fatalf("expect one of edge to be 1 but not found")
	}

	if indexOf("2", ids) == -1 {
		t.Fatalf("expect one of edge to be 1 but not found")
	}

	err := storage.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadOnly(t *testing.T) {
	storage := createStorage()
	root := graph.NewNode("root", storage)

	root.Path("users.1")
	root.Path("users.2")

	readonlyRoot := root.Readonly()
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

	err := storage.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSet(t *testing.T) {
	storage := createStorage()
	root := graph.NewNode("root", storage)

	users := root.Path("users")
	admins := root.Path("admins")
	admin1 := root.Path("admin1")
	admin2 := root.Path("admin2")

	// similar to root.Path("users.admins")
	users.Set(admins)

	// similar to root.Path("users.admins.admin1")
	users.Path("admins").Set(admin1)
	// similar to root.Path("users.admins.admin2")
	users.Path("admins").Set(admin2)

	if admins.Size() != 2 {
		t.Fatalf("expect admins to have 2 users but got %d", admins.Size())
	}

	err := storage.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	createStorage = func() graph.Storage {
		return graph.NewMemoryStorage()
	}

	code := m.Run()

	if code != 0 {
		os.Exit(code)
	}
}
