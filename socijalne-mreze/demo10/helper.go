package demo10

import (
	"fmt"
	"math"
	"slices"

	"gonum.org/v1/gonum/graph"
)

func calculateJaccardCoefficient(g graph.Graph, node1 graph.Node, node2 graph.Node) LinkPredictionScore {
	node1_neigbour_map := make(map[int64]bool)

	common_neighbours := 0

	neighbours1 := graph.NodesOf(g.From(node1.ID()))
	neighbours2 := graph.NodesOf(g.From(node2.ID()))

	for _, neighbour := range neighbours1 {
		node1_neigbour_map[neighbour.ID()] = true
	}

	for _, neighbour := range neighbours2 {
		_, ok := node1_neigbour_map[neighbour.ID()]
		if ok {
			common_neighbours += 1
		}
	}

	union_neighbours := len(neighbours1) + len(neighbours2) - common_neighbours

	score := 0.0
	if union_neighbours > 0 {
		score = float64(common_neighbours) / float64(union_neighbours)
	}

	return LinkPredictionScore{FromID: node1.ID(), ToID: node2.ID(), Score: score}

}

func JaccardCoefficient(g graph.Graph) []LinkPredictionScore {
	scores := []LinkPredictionScore{}

	nodes := graph.NodesOf(g.Nodes())

	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			if g.HasEdgeBetween(nodes[i].ID(), nodes[j].ID()) {
				continue
			}
			score := calculateJaccardCoefficient(g, nodes[i], nodes[j])
			scores = append(scores, score)
		}
	}

	return scores
}

func calculateAdamicAdarIndex(g graph.Graph, node1 graph.Node, node2 graph.Node) LinkPredictionScore {
	node1_neigbour_map := make(map[int64]bool)

	neighbours1 := graph.NodesOf(g.From(node1.ID()))
	neighbours2 := graph.NodesOf(g.From(node2.ID()))

	score := 0.0

	for _, neighbour := range neighbours1 {
		node1_neigbour_map[neighbour.ID()] = true
	}

	for _, neighbour := range neighbours2 {
		_, ok := node1_neigbour_map[neighbour.ID()]
		if ok {
			degree := g.From(neighbour.ID()).Len()

			score += 1.0 / math.Log(float64(degree))
		}
	}

	return LinkPredictionScore{FromID: node1.ID(), ToID: node2.ID(), Score: score}

}

func AdamicAdarIndex(g graph.Graph) []LinkPredictionScore {
	scores := []LinkPredictionScore{}

	nodes := graph.NodesOf(g.Nodes())

	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			if g.HasEdgeBetween(nodes[i].ID(), nodes[j].ID()) {
				continue
			}

			score := calculateAdamicAdarIndex(g, nodes[i], nodes[j])
			scores = append(scores, score)
		}
	}

	return scores
}

func printTop10(scores []LinkPredictionScore) {
	slices.SortFunc(scores, func(a LinkPredictionScore, b LinkPredictionScore) int {
		if a.Score > b.Score {
			return -1
		}
		if a.Score < b.Score {
			return 1
		}
		if a.FromID != b.FromID {
			return int(a.FromID - b.FromID)
		}
		return int(a.ToID - b.ToID)
	})

	limit := min(10, len(scores))
	fmt.Printf("Node1\tNode2\tScore\n")
	for i := 0; i < limit; i++ {
		fmt.Printf("%d\t%d\t%.6f\n", scores[i].FromID, scores[i].ToID, scores[i].Score)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
