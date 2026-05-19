package demo7

import (
	"fmt"
	"math"
	"slices"

	"gonum.org/v1/gonum/graph"
)

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

func printTop10Centrality(centrality map[int64]float64, g graph.Graph) {
	type entry struct {
		nodeID int64
		value  float64
	}

	entries := make([]entry, 0, len(centrality))
	for id, val := range centrality {
		entries = append(entries, entry{id, val})
	}

	slices.SortFunc(entries, func(a, b entry) int {
		if a.value < b.value {
			return 1
		} else if a.value > b.value {
			return -1
		}
		return 0
	})

	limit := 10
	if len(entries) < 10 {
		limit = len(entries)
	}
	for i := 0; i < limit; i++ {
		fmt.Printf("%d. Node %d: %.6f\n", i+1, entries[i].nodeID, entries[i].value)
	}
}
