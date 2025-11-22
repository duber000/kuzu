package csrgraph

import "iter"

// NodeID represents a node identifier
type NodeID uint32

// CSRGraph represents a directed graph in Compressed Sparse Row format
type CSRGraph struct {
	nodeCount uint32
	edgeCount uint32
	offsets   []uint32  // nodeCount + 1 elements
	edges     []NodeID  // edgeCount elements
}

// GraphBuilder helps construct a CSR graph
type GraphBuilder struct {
	adjList map[NodeID][]NodeID
}

// NewBuilder creates a new graph builder
func NewBuilder() *GraphBuilder {
	return &GraphBuilder{
		adjList: make(map[NodeID][]NodeID),
	}
}

// AddNode adds a node to the graph
func (b *GraphBuilder) AddNode(node NodeID) {
	// TODO: Implement node addition
	// Initialize empty adjacency list if not exists
}

// AddEdge adds a directed edge from src to dst
func (b *GraphBuilder) AddEdge(src, dst NodeID) {
	// TODO: Implement edge addition
	// Add dst to src's adjacency list
}

// Build constructs the CSR graph from the adjacency list
func (b *GraphBuilder) Build() *CSRGraph {
	// TODO: Implement CSR construction
	// 1. Find max node ID
	// 2. Allocate offsets array (size = maxNode + 2)
	// 3. Fill offsets by counting edges
	// 4. Allocate edges array
	// 5. Copy edges from adjacency lists
	return nil
}

// NodeCount returns the number of nodes in the graph
func (g *CSRGraph) NodeCount() uint32 {
	return g.nodeCount
}

// EdgeCount returns the number of edges in the graph
func (g *CSRGraph) EdgeCount() uint32 {
	return g.edgeCount
}

// Degree returns the out-degree of a node
func (g *CSRGraph) Degree(node NodeID) uint32 {
	// TODO: Implement degree computation
	// degree = offsets[node+1] - offsets[node]
	return 0
}

// Neighbors returns an iterator over the neighbors of a node
// Uses Go 1.23 iter.Seq for efficient iteration
func (g *CSRGraph) Neighbors(node NodeID) iter.Seq[NodeID] {
	return func(yield func(NodeID) bool) {
		// TODO: Implement neighbor iteration
		// 1. Check node bounds
		// 2. Get start and end offsets
		// 3. Iterate edges[start:end]
		// 4. Call yield for each neighbor
		// 5. Return early if yield returns false
	}
}

// Edges returns an iterator over all edges in the graph
// Returns (src, dst) pairs using Go 1.23 iter.Seq2
func (g *CSRGraph) Edges() iter.Seq2[NodeID, NodeID] {
	return func(yield func(NodeID, NodeID) bool) {
		// TODO: Implement edge iteration
		// For each node, iterate its neighbors
	}
}

// Has2Hop checks if there is a path from src to dst within 2 hops
func (g *CSRGraph) Has2Hop(src, dst NodeID) bool {
	// TODO: Implement 2-hop connectivity check
	// 1. Check direct edge (1-hop)
	// 2. Check paths through intermediates (2-hop)
	return false
}

// TwoHopNeighbors returns an iterator over all nodes reachable in 2 hops
func (g *CSRGraph) TwoHopNeighbors(node NodeID) iter.Seq[NodeID] {
	return func(yield func(NodeID) bool) {
		// TODO: Implement 2-hop neighbor iteration
		// Use nested iteration over neighbors
	}
}

// Iterator composition helpers

// Filter returns an iterator that only yields elements matching the predicate
func Filter[T any](seq iter.Seq[T], pred func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		// TODO: Implement filter
		// Iterate seq, only yield elements where pred is true
	}
}

// Map transforms elements using the given function
func Map[T, U any](seq iter.Seq[T], fn func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		// TODO: Implement map
		// Iterate seq, yield transformed elements
	}
}

// Take returns an iterator that yields at most n elements
func Take[T any](seq iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		// TODO: Implement take
		// Count elements and stop after n
	}
}

// Collect gathers all elements from an iterator into a slice
func Collect[T any](seq iter.Seq[T]) []T {
	var result []T
	for v := range seq {
		result = append(result, v)
	}
	return result
}
