package demo5

import (
	"fmt"
	"slices"

	"github.com/korizma/socnet-lab-go/demo2"
	"gonum.org/v1/gonum/graph/topo"
)

func Demo5() {
	// g, err := demo2.LoadZachary()
	// g, err := demo2.LoadFlorentine()
	// g, err := demo2.LoadWomen()
	g, err := demo2.LoadMiserables()
	if err != nil {
		fmt.Println("Error loading graph:", err)
		return
	}

	components := topo.ConnectedComponents(g)

	slices.SortFunc(components, cmp)

	fmt.Println("Connected components:", len(components))
	for i, comp := range components {
		fmt.Printf("Component %d: %d nodes\n", i, len(comp))
	}

	fmt.Println()

	n := int64(1000)
	p := 0.003
	randG := GenerateErdosRenyiGraph(n, p)

	fmt.Printf("Erdos-Renyi graph: n=%d, p=%f\n", n, p)
	fmt.Println("Nodes:", randG.Nodes().Len())
	fmt.Println("Edges:", randG.Edges().Len())

	randComponents := topo.ConnectedComponents(randG)
	slices.SortFunc(randComponents, cmp)

	fmt.Println("Connected components in random graph:", len(randComponents))
	for i, comp := range randComponents {
		fmt.Printf("Component %d: %d nodes\n", i, len(comp))
	}
}
