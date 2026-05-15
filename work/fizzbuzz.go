package main

import "fmt"

func main() {
	var fizz int = 0
	var buzz int = 0
	var fizzbuzz int = 0
	for i := 1; i <= 100; i++ {
		fmt.Printf("%d ", i)
		if i%3 == 0 {
			fmt.Printf("Fizz")
			fizz++
		}
		if i%5 == 0 {
			fmt.Printf("Buzz")
			buzz++
		}
		if i%3 == 0 && i%5 == 0 {
			fizzbuzz++
			fizz--
			buzz--
		}
		fmt.Printf("\n")
	}
	fmt.Printf("Fizz: %d, Buzz: %d, FizzBuzz: %d", fizz, buzz, fizzbuzz)
}
