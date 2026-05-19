package solutions

import (
	"fmt"

	"github.com/korizma/socnet-lab-go/demo5"
	"github.com/korizma/socnet-lab-go/lab"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

func LoadMiserables() (*simple.UndirectedGraph, error) {
	return lab.LoadGraph("les_miserables.txt")
}

func Sol2() {

	n := int64(100)
	p := 0.001
	g := demo5.GenerateErdosRenyiGraph(n, p)
	randComponents := topo.ConnectedComponents(g)
	for len(randComponents) != 1 {
		p += 0.001
		g = demo5.GenerateErdosRenyiGraph(n, p)
		randComponents = topo.ConnectedComponents(g)
	}
	fmt.Printf("Kriticna vrednost je %.2f", p)

}
