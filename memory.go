package graph

var _ Storage = &MemoryStorage{}
var _ Bucket = &MemoryBucket{}

type MemoryStorage struct {
	buckets map[string]*MemoryBucket
}

func (m *MemoryStorage) Bucket(name string) Bucket {
	bucket, ok := m.buckets[name]
	if !ok {
		bucket = &MemoryBucket{
			bucket: make(map[string]*Node),
		}
		m.buckets[name] = bucket
	}

	return bucket
}

func (m *MemoryStorage) Close() error {
	m.buckets = nil
	return nil
}

type MemoryBucket struct {
	bucket map[string]*Node
}

func (mb *MemoryBucket) Get(name string) (*Node, bool) {
	node, ok := mb.bucket[name]
	return node, ok
}

func (mb *MemoryBucket) Set(name string, node *Node) {
	if _, ok := mb.bucket[name]; ok {
		return
	}

	mb.bucket[name] = node
}

func (mb *MemoryBucket) Map(fn func(*Node) bool) {
	for name := range mb.bucket {
		if done := fn(mb.bucket[name]); done {
			break
		}
	}
}

func (mb *MemoryBucket) Size() int {
	return len(mb.bucket)
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		buckets: make(map[string]*MemoryBucket),
	}
}
