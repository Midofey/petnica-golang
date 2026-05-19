package demo9

import (
	"slices"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/community"
)

func updateLabel(g graph.Graph, node_id int64, labels map[int64]int64, neighbours map[int64][]graph.Node) (int64, bool) {
	seen := make(map[int64]int)

	ns := neighbours[node_id]

	for _, neighbour := range ns {
		seen[labels[neighbour.ID()]]++
	}

	max_id := int64(-1)
	for id, val := range seen {
		if max_id == -1 {
			max_id = id
			continue
		}

		if seen[max_id] < val {
			max_id = id
		}
	}

	change := true
	if labels[node_id] == max_id {
		change = false
	}

	if max_id != -1 {
		return max_id, change
	}

	return labels[node_id], change
}

func LabelPropagation(g graph.Graph) [][]graph.Node {
	labels := make(map[int64]int64)
	new_labels := make(map[int64]int64)
	nodes := graph.NodesOf(g.Nodes())
	for _, node := range nodes {
		labels[node.ID()] = node.ID()
	}

	change := true
	max_iterations := 100
	iteration := 0

	// optimization things
	neighbours := make(map[int64][]graph.Node)

	for _, node := range nodes {
		neighbours[node.ID()] = graph.NodesOf(g.From(node.ID()))
	}

	for change {
		iteration++
		if iteration > max_iterations {
			break
		}
		change = false
		for _, node := range nodes {
			new_label, changed := updateLabel(g, node.ID(), labels, neighbours)

			change = change || changed

			new_labels[node.ID()] = new_label

		}
		labels = new_labels
		new_labels = make(map[int64]int64)
	}

	communities := make(map[int64][]graph.Node)

	for node_id, label := range labels {

		_, ok := communities[label]
		if !ok {
			communities[label] = []graph.Node{}
		}
		communities[label] = append(communities[label], g.Node(node_id))
	}

	comm_slice := make([][]graph.Node, 0, len(communities))

	for _, community := range communities {
		comm_slice = append(comm_slice, community)
	}

	return comm_slice
}

func getMockCommunities(communities [][]graph.Node, i int, j int) [][]graph.Node {
	new_c := make([][]graph.Node, 0)

	for indx, comm := range communities {
		if indx != i && indx != j {
			new_c = append(new_c, comm)
		}
	}

	ij_comm := make([]graph.Node, 0)

	for _, node := range communities[i] {
		ij_comm = append(ij_comm, node)
	}

	for _, node := range communities[j] {
		ij_comm = append(ij_comm, node)
	}

	new_c = append(new_c, ij_comm)

	return new_c
}

func GreedyModularityCommunities(g graph.Graph) [][]graph.Node {
	communities := make([][]graph.Node, 0)

	nodes := graph.NodesOf(g.Nodes())

	for _, node := range nodes {
		communities = append(communities, []graph.Node{node})
	}

	bestModularity := community.Q(g, communities, 1)

	for {
		bestImprovement := 0.0
		bestI := -1
		bestJ := -1

		for i := 0; i < len(communities); i++ {
			for j := i + 1; j < len(communities); j++ {
				communities_new := getMockCommunities(communities, i, j)

				modularity := community.Q(g, communities_new, 1)
				delta := modularity - bestModularity

				if delta > bestImprovement {
					bestImprovement = delta
					bestI = i
					bestJ = j
				}
			}
		}

		if bestImprovement == 0 {
			break
		}

		temp := communities[bestI]
		communities[bestI] = communities[len(communities)-1]
		communities[len(communities)-1] = temp

		for _, value := range temp {
			communities[bestJ] = append(communities[bestJ], value)
		}

		communities = communities[:len(communities)-1]
		bestModularity = community.Q(g, communities, 1)
	}

	return communities
}

func sortPartition(partition [][]graph.Node) {
	for _, communityNodes := range partition {
		slices.SortFunc(communityNodes, func(a graph.Node, b graph.Node) int {
			return int(a.ID() - b.ID())
		})
	}
	slices.SortFunc(partition, func(a []graph.Node, b []graph.Node) int {
		if len(a) == 0 && len(b) == 0 {
			return 0
		}
		if len(a) == 0 {
			return -1
		}
		if len(b) == 0 {
			return 1
		}
		return int(a[0].ID() - b[0].ID())
	})
}

func nodeIDs(nodes []graph.Node) []int64 {
	ids := make([]int64, 0, len(nodes))
	for _, node := range nodes {
		ids = append(ids, node.ID())
	}
	return ids
}
