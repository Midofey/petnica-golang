package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scores := map[string][]int{
		"Sweden":   []int{10, 7, 12},
		"Ukraine":  []int{12, 8, 10},
		"Portugal": []int{6, 12, 8},
		"Norway":   []int{8, 5, 7},
	}
	var opcija int = 0
	reader := bufio.NewReader(os.Stdin)
	for opcija != 4 {

		fmt.Println("1. Add a country")
		fmt.Println("2. Add a score for a country")
		fmt.Println("3. Print standings")
		fmt.Println("4. Quit")
		fmt.Scan(&opcija)
		reader.ReadString('\n')

		switch opcija {
		case 1:
			fmt.Print("Enter country: ")
			country, _ := reader.ReadString('\n')
			country = strings.TrimSpace(country)
			scores[country] = []int{}
		case 2:
			var score int = 0
			fmt.Print("Enter country: ")
			country, _ := reader.ReadString('\n')
			country = strings.TrimSpace(country)
			_, ok := scores[country]
			if !ok {
				fmt.Println("Invalid country, please add it first")
			}
			fmt.Print("Enter score: ")
			fmt.Scan(&score)
			for score > 12 || score == 9 || score == 11 {
				fmt.Print("Invalid score, try again: ")
				fmt.Scan(&score)
			}

			scores[country] = append(scores[country], score)

		case 3:

			for key, val := range scores {
				if len(val) == 0 {
					continue
				}

				var total int = 0
				var average float32 = 0.0
				var highest int = val[0]
				fmt.Printf("%s: ", key)
				for _, v := range val {
					total += v
					if v > highest {
						highest = v
					}
					fmt.Printf("%d, ", v)
				}
				average = float32(total) / float32(len(val))
				fmt.Printf("\nTotal:%d\nAverage:%.2f\nHighest:%d\n", total, average, highest)
				fmt.Println()
			}
		}
	}
}
