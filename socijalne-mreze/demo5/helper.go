package demo5

import (
	"math/rand"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

func GenerateErdosRenyiGraph(n int64, p float64) *simple.UndirectedGraph {

	G := simple.NewUndirectedGraph()

	for i := int64(0); i < n; i++ {
		G.AddNode(simple.Node(i))
	}

	rng := rand.New(rand.NewSource(42))

	nodes := graph.NodesOf(G.Nodes())

	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			if rng.Float64() < p {
				G.SetEdge(G.NewEdge(nodes[i], nodes[j]))
			}
		}
	}

	return G
}

func cmp(a []graph.Node, b []graph.Node) int {
	return len(b) - len(a)
}
