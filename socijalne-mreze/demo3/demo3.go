package demo3

import (
	"fmt"

	"github.com/korizma/socnet-lab-go/demo2"
	"gonum.org/v1/gonum/graph"
)

func DegreeDistribution(g graph.Undirected) (map[int]float32, int) {
	degree_dist := make(map[int]float32)

	nodes := graph.NodesOf(g.Nodes())
	for _, node := range nodes {
		degree := g.From(node.ID()).Len()
		degree_dist[degree]++
	}

	max_degree := 0
	for degree := range degree_dist {
		degree_dist[degree] /= float32(len(nodes))
		max_degree = max(max_degree, degree)
	}

	return degree_dist, max_degree
}

func Demo3() {
	// g, err := demo2.LoadZachary()
	// g, err := demo2.LoadFlorentine()
	// g, err := demo2.LoadWomen()
	g, err := demo2.LoadMiserables()
	if err != nil {
		fmt.Println("Error loading graph:", err)
		return
	}

	dist, maxDegree := DegreeDistribution(g)

	fmt.Printf("Distribucija stepeni cvorova:\n")

	for i := 0; i <= maxDegree; i++ {
		_, ok := dist[i]
		if ok {
			fmt.Printf("Degree %d: %.2f%%\n", i, dist[i]*100)
		}
	}
}
