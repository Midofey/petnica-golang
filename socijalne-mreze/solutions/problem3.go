package solutions

import (
	"fmt"

	"github.com/korizma/socnet-lab-go/lab"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

func LoadWomen() (*simple.UndirectedGraph, error) {
	return lab.LoadGraph("southern_women.txt")
}

func dist(shortest path.AllShortest, node1 graph.Node, node2 graph.Node) int {
	_, dist, _ := shortest.Between(node1.ID(), node2.ID())
	return int(dist)
}

func Sol3() {
	distribucija := make(map[int]float64)
	g, _ := LoadWomen()
	slice := graph.NodesOf(g.Nodes())
	shortest := path.DijkstraAllPaths(g)
	for _, x := range slice {
		for _, y := range slice {
			if dist(shortest, x, y) <= 0 {
				continue
			}
			distribucija[dist(shortest, x, y)]++
		}
	}
	for distanca, broj := range distribucija {
		verovatnoca := broj / (float64)(len(slice)*len(slice)-len(slice))

		fmt.Printf("Verovatnoca za distancu %d je: %.2f%s\n", distanca, verovatnoca*100, "%")
	}

}
