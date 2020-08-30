package graph

type adjacencyList map[string][]string

// Node definition
type Node struct {
	Key        string
	Val        string
	neighbours []*Node
}

// Graph definition
type Graph struct {
	root    *Node
	adjList adjacencyList
	nodes   []string
}

// NewGraph from adjacency list and node list
func NewGraph(a adjacencyList, nodes []string) Graph {
	g := Graph{}
	m := make(map[string]*Node)

	g.adjList = a
	g.nodes = nodes

	// create nodes
	for _, n := range nodes {
		node := Node{Key: n}
		m[n] = &node
	}
	g.root = m[nodes[0]]

	// create edges
	for n, adj := range a {
		for _, neighbour := range adj {
			m[n].neighbours = append(m[n].neighbours, m[neighbour])
		}
	}

	return g
}
