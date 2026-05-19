package main

import (
	"fmt"
	"os"

	"github.com/korizma/socnet-lab-go/demo1"
	"github.com/korizma/socnet-lab-go/demo10"
	"github.com/korizma/socnet-lab-go/demo2"
	"github.com/korizma/socnet-lab-go/demo3"
	"github.com/korizma/socnet-lab-go/demo4"
	"github.com/korizma/socnet-lab-go/demo5"
	"github.com/korizma/socnet-lab-go/demo6"
	"github.com/korizma/socnet-lab-go/demo7"
	"github.com/korizma/socnet-lab-go/demo8"
	"github.com/korizma/socnet-lab-go/demo9"
	"github.com/korizma/socnet-lab-go/solutions"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run main.go <demo number>")
		return
	}

	demo_num := os.Args[1]

	switch demo_num {
	case "1":
		demo1.Demo1()
	case "2":
		demo2.Demo2()
	case "3":
		demo3.Demo3()
	case "4":
		demo4.Demo4()
	case "5":
		demo5.Demo5()
	case "6":
		demo6.Demo6()
	case "7":
		demo7.Demo7()
	case "8":
		demo8.Demo8()
	case "9":
		demo9.Demo9()
	case "10":
		demo10.Demo10()
	case "sol1":
		solutions.Sol1()
	case "sol2":
		solutions.Sol2()
	case "sol3":
		solutions.Sol3()
	case "sol4":
		solutions.Sol4()
	case "sol5":
		solutions.Sol5()
	}
}
