package lab

import (
	"bufio"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/graph/simple"
)

func LoadGraph(filename string) (*simple.UndirectedGraph, error) {
	g := simple.NewUndirectedGraph()

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}

		fromID, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			continue
		}
		toID, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			continue
		}
		if fromID == toID {
			continue
		}

		from, _ := g.NodeWithID(fromID)
		to, _ := g.NodeWithID(toID)
		g.SetEdge(g.NewEdge(from, to))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return g, nil
}

// - l     -- broj zajednica
// - k     -- broj čvorova u svakoj od zajednica
// - p_in  -- verovatnoća vezivanja čvorova unutar iste zajednice
// - p_out -- verovatnoća vezivanja čvorova koji su u različitim zajednicama.
func PlantedPartitionGraph(num_com int, num_nodes int, p_in float64, p_out float64) *simple.UndirectedGraph {
	g := simple.NewUndirectedGraph()

	for i := 0; i < num_com*num_nodes; i++ {
		comm_i := i / num_nodes
		node_i, _ := g.NodeWithID(int64(i))
		for j := i + 1; j < num_com*num_nodes; j++ {
			node_j, _ := g.NodeWithID(int64(j))
			comm_j := j / num_nodes

			if comm_j == comm_i {
				if rand.Float64() < p_in {
					newedge := g.NewEdge(node_j, node_i)
					g.SetEdge(newedge)
				}
			} else {
				if rand.Float64() < p_out {
					newedge := g.NewEdge(node_j, node_i)
					g.SetEdge(newedge)
				}
			}
		}
	}

	return g
}
