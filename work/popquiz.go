package main

import "fmt"

func printResults(score map[string][]int) {
	for key, val := range score {
		var ukupno int = 0
		var najbolji int = val[0]
		var prosek float32 = 0.0
		fmt.Println()
		fmt.Printf("Tim %s:\n", key)
		for i, v := range val {
			fmt.Printf("Rezultat %d. runde je: %d\n", i+1, v)
			ukupno += v
			if v > najbolji {
				najbolji = v
			}
		}
		prosek = float32(ukupno) / float32(len(val))
		fmt.Printf("Ukupan rezultat: %d, Najbolji rezultat: %d, Prosek: %.2f\n", ukupno, najbolji, prosek)
	}
}

func main() {
	teamscore := map[string][]int{
		"Alpha": []int{8, 5, 9, 6},
		"Beta":  []int{4, 7, 7, 10},
		"Gamma": []int{6, 6, 4, 8},
	}

	teamscore["Alpha"] = append(teamscore["Alpha"], 5)
	teamscore["Beta"] = append(teamscore["Beta"], 3)
	teamscore["Gamma"] = append(teamscore["Gamma"], 1)
	printResults(teamscore)
}
