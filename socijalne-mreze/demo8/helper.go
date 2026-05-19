package demo8

import (
	"fmt"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

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

func ShellIndex(g simple.UndirectedGraph, print_flag bool) map[int64]int {
	new_graph := CopyGraph(g)
	shell_index := make(map[int64]int)

	k := 0

	for {
		k += 1

		for {
			to_be_removed := make([]graph.Node, 0)

			nodes := graph.NodesOf(new_graph.Nodes())
			for _, node := range nodes {

				if new_graph.From(node.ID()).Len() < k {
					shell_index[node.ID()] = k - 1
					to_be_removed = append(to_be_removed, node)
				}
			}

			if len(to_be_removed) == 0 {
				break
			}

			for _, node := range to_be_removed {
				new_graph.RemoveNode(node.ID())
			}
		}

		if new_graph.Nodes().Len() == 0 {
			break
		}

		if print_flag {
			fmt.Println("Core", k, "#Nodes", new_graph.Nodes().Len(), "#Links", new_graph.Edges().Len())
		}
	}

	return shell_index
}
