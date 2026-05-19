package solutions

import (
	"fmt"
	"math"

	"github.com/korizma/socnet-lab-go/lab"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/network"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

func LoadZachary() (*simple.UndirectedGraph, error) {
	return lab.LoadGraph("zachary.txt")
}

func canPropagate(g graph.Graph) bool {
	shortest := path.DijkstraAllPaths(g)
	path, _, _ := shortest.Between(4, 7)

	if path != nil {
		fmt.Println("Nije izgubio funckionalnost propagiranja")
		return true
	} else {
		fmt.Println("Izgubio je funkcionalnost propagiranja")
		return false
	}
}

func CopyGraph(g simple.UndirectedGraph) *simple.UndirectedGraph {
	copy := simple.NewUndirectedGraph()

	// Copy nodes
	nodes := graph.NodesOf(g.Nodes())
	for _, node := range nodes {
		copy.AddNode(simple.Node(node.ID()))
	}

	// Copy edges
	edges := graph.EdgesOf(g.Edges())
	for _, edge := range edges {
		copy.SetEdge(copy.NewEdge(simple.Node(edge.From().ID()), simple.Node(edge.To().ID())))
	}

	return copy
}

func EigenvectorCentrality(g graph.Graph) map[int64]float64 {
	ec_map_prev := make(map[int64]float64)
	ec_map_curr := make(map[int64]float64)

	nodes := graph.NodesOf(g.Nodes())

	for _, node := range nodes {
		ec_map_prev[node.ID()] = 1
	}

	iterations := 100

	for iterations > 0 {
		iterations -= 1

		for _, node := range nodes {
			ec_map_curr[node.ID()] = 0
			neighbours := graph.NodesOf(g.From(node.ID()))

			for _, neighbour := range neighbours {
				ec_map_curr[node.ID()] += ec_map_prev[neighbour.ID()]
			}
		}

		norm := 0.0
		for _, val := range ec_map_curr {
			norm += val * val
		}
		norm = math.Sqrt(norm)

		for id := range ec_map_curr {
			ec_map_curr[id] /= norm
		}

		temp := ec_map_curr
		ec_map_curr = ec_map_prev
		ec_map_prev = temp
	}
	return ec_map_prev
}

type entry struct {
	nodeID int64
	value  float64
}

func Sol1() {
	g, err := LoadZachary()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	shortest := path.DijkstraAllPaths(g)

	// Betweenness centrality
	betweenness := network.Betweenness(g)
	fmt.Println("Betweeness centrality:")

	entries := make([]entry, 0, len(betweenness))
	for id, val := range betweenness {
		entries = append(entries, entry{id, val})
	}
	node := entries[0]

	for _, n := range entries {
		if n.value > node.value {
			node = n
		}
	}
	copyg := CopyGraph(*g)
	g.RemoveNode(node.nodeID)
	canPropagate(g)

	// Closeness centrality
	fmt.Println("Closeness centrality:")
	closeness := network.Closeness(copyg, shortest)
	entries = make([]entry, 0, len(closeness))
	for id, val := range closeness {
		entries = append(entries, entry{id, val})
	}
	node = entries[0]
	for _, n := range entries {
		if n.value > node.value {
			node = n
		}
	}
	newcopyg := CopyGraph(*copyg)
	copyg.RemoveNode(node.nodeID)
	canPropagate(copyg)

	// fmt.Println()

	// Eigenvector centrality (stub)
	eigenvector := EigenvectorCentrality(g)
	fmt.Println("Eigenvector centrality:")
	entries = make([]entry, 0, len(eigenvector))
	for id, val := range eigenvector {
		entries = append(entries, entry{id, val})
	}
	node = entries[0]
	for _, n := range entries {
		if n.value > node.value {
			node = n
		}
	}
	canPropagate(newcopyg)

}
