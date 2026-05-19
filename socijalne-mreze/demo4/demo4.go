package demo4

import (
	"fmt"

	"github.com/korizma/socnet-lab-go/demo2"
	"gonum.org/v1/gonum/graph"
)

func ClusteringCoefficient(g graph.Undirected, nodeID int64) float64 {
	neighbors_map := make(map[int64]bool)

	neighbors := graph.NodesOf(g.From(nodeID))
	for _, neighbor := range neighbors {
		neighbors_map[neighbor.ID()] = true
	}

	edges_between_neighbors := 0

	for _, neighbor := range neighbors {
		neighbor_neighbors := graph.NodesOf(g.From(neighbor.ID()))
		for _, neighbor_neighbor := range neighbor_neighbors {
			if neighbors_map[neighbor_neighbor.ID()] {
				edges_between_neighbors++
			}
		}
	}

	n := len(neighbors)
	if n < 2 {
		return 0.0
	}

	return float64(edges_between_neighbors) / float64(n*(n-1))
}

func Demo4() {
	// g, err := demo2.LoadZachary()
	// g, err := demo2.LoadFlorentine()
	// g, err := demo2.LoadWomen()
	g, err := demo2.LoadMiserables()
	if err != nil {
		fmt.Println("Error loading graph:", err)
		return
	}

	nodes := graph.NodesOf(g.Nodes())
	for _, node := range nodes {
		cc := ClusteringCoefficient(g, node.ID())
		fmt.Printf("Node %d: Clustering Coefficient = %.4f\n", node.ID(), cc)
	}
}
