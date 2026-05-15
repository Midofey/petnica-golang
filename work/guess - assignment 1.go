package main

import (
	"fmt"
	"math/rand/v2"
)

func main() {
	var broj int = rand.IntN(100) + 1
	var pokusaji int = 7
	var score int = 0
	var runda = 1
	for runda <= 5 {
		for pokusaji > 0 {
			var trenutno int = 0
			fmt.Printf("Unesite vas pokusaj: ")
			fmt.Scan(&trenutno)
			if trenutno == broj {
				fmt.Printf("Nasli ste tacan broj!\n")
				break
			} else if trenutno > broj {
				fmt.Printf("Broj je manji\n")
			} else {
				fmt.Printf("Broj je veci\n")
			}
			pokusaji--
		}
		switch pokusaji {
		case 6:
			score += 10
		case 5:
			score += 7
		case 4:
			score += 5
		case 3:
			score += 3
		default:
			score += 1
		}
		fmt.Printf("Bilo ti je potrebno %d pokusaja, nakon %d runde, imas %d poena\n", 7-pokusaji, runda, score)
		runda++
		pokusaji = 7
	}

	if score >= 35 {
		fmt.Printf("Incredible")
	} else if score >= 20 {
		fmt.Printf("Solid")
	} else if score >= 10 {
		fmt.Printf("Room to improve")
	} else {
		fmt.Printf("Better luck next time")
	}

}
