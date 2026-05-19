package demo7

import (
	"fmt"

	"github.com/korizma/socnet-lab-go/demo2"
	"gonum.org/v1/gonum/graph/network"
	"gonum.org/v1/gonum/graph/path"
)

func Demo7() {
	// g, err := demo2.LoadZachary()
	// g, err := demo2.LoadFlorentine()
	// g, err := demo2.LoadWomen()
	g, err := demo2.LoadMiserables()
	if err != nil {
		fmt.Println("Error loading graph:", err)
		return
	}

	// Betweenness centrality
	betweenness := network.Betweenness(g)

	fmt.Println("Betweeness centrality:")
	printTop10Centrality(betweenness, g)

	fmt.Println()

	// Closeness centrality
	closeness := network.Closeness(g, path.DijkstraAllPaths(g))

	fmt.Println("Closeness centrality:")
	printTop10Centrality(closeness, g)

	fmt.Println()

	// Eigenvector centrality (stub)
	eigenvector := EigenvectorCentrality(g)
	fmt.Println("Eigenvector centrality:")
	printTop10Centrality(eigenvector, g)
}
