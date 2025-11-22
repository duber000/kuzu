package art

import "iter"

// ART is the Adaptive Radix Tree
type ART struct {
	root *Node
	size int
}

// Node is the base interface for all node types
type Node interface {
	Insert(key []byte, value interface{}, depth int) (Node, bool)
	Search(key []byte, depth int) (interface{}, bool)
	Delete(key []byte, depth int) (Node, bool)
	Iterator() iter.Seq2[[]byte, interface{}]
}

// Node4 has 1-4 children
type Node4 struct {
	prefix      []byte
	keys        [4]byte
	children    [4]Node
	numChildren int
}

// Node16 has 5-16 children
type Node16 struct {
	prefix      []byte
	keys        [16]byte
	children    [16]Node
	numChildren int
}

// Node48 has 17-48 children
type Node48 struct {
	prefix      []byte
	childIndex  [256]byte  // maps key -> index in children
	children    [48]Node
	numChildren int
}

// Node256 has 49-256 children
type Node256 struct {
	prefix      []byte
	children    [256]Node
	numChildren int
}

// Leaf stores the actual value
type Leaf struct {
	key   []byte
	value interface{}
}

// New creates a new ART
func New() *ART {
	return &ART{}
}

// Insert inserts a key-value pair
func (art *ART) Insert(key []byte, value interface{}) {
	// TODO: Implement insertion
	// 1. Start at root
	// 2. Follow path, creating nodes as needed
	// 3. Handle prefix compression
	// 4. Grow nodes when full
}

// Search searches for a key
func (art *ART) Search(key []byte) (interface{}, bool) {
	// TODO: Implement search
	// 1. Start at root
	// 2. Follow path using appropriate node lookups
	// 3. Handle prefix matching
	return nil, false
}

// Delete deletes a key
func (art *ART) Delete(key []byte) bool {
	// TODO: Implement deletion
	// 1. Find and remove key
	// 2. Shrink nodes when sparse
	// 3. Merge nodes when possible
	return false
}

// Range returns an iterator over keys in [start, end)
func (art *ART) Range(start, end []byte) iter.Seq2[[]byte, interface{}] {
	// TODO: Implement range iterator
	return nil
}

// Size returns number of keys
func (art *ART) Size() int {
	return art.size
}

// Node4 methods
func (n *Node4) Insert(key []byte, value interface{}, depth int) (Node, bool) {
	// TODO: Implement Node4 insertion
	// Handle prefix matching
	// Insert into sorted position
	// Grow to Node16 if full
	return nil, false
}

func (n *Node4) Search(key []byte, depth int) (interface{}, bool) {
	// TODO: Implement Node4 search
	return nil, false
}

// Helper functions
func longestCommonPrefix(a, b []byte) int {
	// TODO: Implement longest common prefix
	return 0
}

func checkPrefix(node []byte, key []byte, depth int) int {
	// TODO: Check how much of prefix matches
	return 0
}
