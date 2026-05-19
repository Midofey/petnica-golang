package demo10

import (
	"fmt"

	"github.com/korizma/socnet-lab-go/demo2"
)

type LinkPredictionScore struct {
	FromID int64
	ToID   int64
	Score  float64
}

func Demo10() {
	// G, err := demo2.LoadZachary()
	// G, err := demo2.LoadFlorentine()
	// G, err := demo2.LoadWomen()
	G, err := demo2.LoadMiserables()
	if err != nil {
		fmt.Println("Error loading graph:", err)
		return
	}
	fmt.Println("Jaccard Coefficient:")
	printTop10(JaccardCoefficient(G))
	fmt.Println()

	fmt.Println("Adamic-Adar Index:")
	printTop10(AdamicAdarIndex(G))
}
