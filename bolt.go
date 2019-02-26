package graph

import (
	"fmt"

	"go.etcd.io/bbolt"
)

// each bucket should have a random unique name
// each of those name will store in `buckets` bucket.
// for every unique node, we will have a bucket

var _ Storage = &BoltStorage{}
var _ Bucket = &BoltBucket{}

type BoltStorage struct {
	db *bbolt.DB
}

func (b *BoltStorage) Bucket(name string) Bucket {
	return &BoltBucket{
		name: []byte(name),
		db:   b.db,
	}
}

func (b *BoltStorage) Close() error {
	return b.db.Close()
}

type BoltBucket struct {
	name []byte
	db   *bbolt.DB
}

func (bb *BoltBucket) Get(name string) (*Node, bool) {
	return nil, false
}

func (bb *BoltBucket) Set(name string, node *Node) {
	bb.db.Update(func(tx *bbolt.Tx) error {
		return nil
	})
}

func (bb *BoltBucket) Map(fn func(*Node) bool) {
	bb.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bb.name)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			///
			fmt.Print(v)
			///
		}

		return nil
	})
}

func (bb *BoltBucket) Size() int {
	size := 0

	bb.db.View(func(tx *bbolt.Tx) error {
		size = int(tx.Size())
		return nil
	})

	return size
}

func NewBoltStorage(path string) (*BoltStorage, error) {
	db, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}

	return &BoltStorage{
		db: db,
	}, nil
}
