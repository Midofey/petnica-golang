package demo8

import (
	"fmt"

	"github.com/korizma/socnet-lab-go/demo2"
)

func Demo8() {
	// G, err := demo2.LoadZachary()
	// G, err := demo2.LoadFlorentine()
	// G, err := demo2.LoadWomen()
	G, err := demo2.LoadMiserables()
	if err != nil {
		fmt.Println("Error loading graph:", err)
		return
	}

	shell_index := ShellIndex(*G, true)

	for id, index := range shell_index {
		fmt.Printf("Node %d: Shell Index %d\n", id, index)
	}
}
