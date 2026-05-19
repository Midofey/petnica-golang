package demo2

import (
	"fmt"

	"github.com/korizma/socnet-lab-go/lab"
	"gonum.org/v1/gonum/graph/simple"
)

func LoadZachary() (*simple.UndirectedGraph, error) {
	return lab.LoadGraph("zachary.txt")
}

func LoadMiserables() (*simple.UndirectedGraph, error) {
	return lab.LoadGraph("les_miserables.txt")
}

func LoadWomen() (*simple.UndirectedGraph, error) {
	return lab.LoadGraph("southern_women.txt")
}

func LoadFlorentine() (*simple.UndirectedGraph, error) {
	return lab.LoadGraph("florentine.txt")
}

func Demo2() {
	g, err := LoadZachary()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("Velicine grafova:\n")

	fmt.Printf("Nodes: %d, Edges: %d\n", g.Nodes().Len(), g.Edges().Len())

	g, err = LoadFlorentine()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("Nodes: %d, Edges: %d\n", g.Nodes().Len(), g.Edges().Len())

	g, err = LoadMiserables()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("Nodes: %d, Edges: %d\n", g.Nodes().Len(), g.Edges().Len())

	g, err = LoadWomen()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("Nodes: %d, Edges: %d\n", g.Nodes().Len(), g.Edges().Len())
}
