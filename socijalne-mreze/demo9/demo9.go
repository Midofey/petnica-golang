package demo9

import (
	"fmt"

	"github.com/korizma/socnet-lab-go/demo2"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/community"
)

func Louvain(g graph.Undirected) {
	reduced := community.Modularize(g, 1, nil)
	comms := reduced.Communities()
	examinePartition(g, comms)
}

func GMO(g graph.Undirected) {
	comms := GreedyModularityCommunities(g)
	examinePartition(g, comms)
}

func LP(g graph.Undirected) {
	comms := LabelPropagation(g)
	examinePartition(g, comms)
}

func examinePartition(g graph.Undirected, partition [][]graph.Node) {
	sortPartition(partition)

	fmt.Println("Broj zajednica", len(partition))
	for _, communityNodes := range partition {
		fmt.Println(nodeIDs(communityNodes))
	}

	q := community.Q(g, partition, 1)
	fmt.Println("Modularnost", q)
}

func Demo9() {
	// G, err := demo2.LoadZachary()
	// G, err := demo2.LoadFlorentine()
	// G, err := demo2.LoadWomen()
	G, err := demo2.LoadMiserables()
	if err != nil {
		fmt.Println("Error loading graph:", err)
		return
	}
	LP(G)
	fmt.Println()

	Louvain(G)
	fmt.Println()

	GMO(G)
}
